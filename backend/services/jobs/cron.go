package jobs

import (
	"fmt"
	"log"
	"strings"
	"time"

	"pondok-tani-backend/config"
	notificationhandler "pondok-tani-backend/handlers/notification"
	maintenancemodels "pondok-tani-backend/models/maintenance"
	notificationmodels "pondok-tani-backend/models/notification"

	"github.com/robfig/cron/v3"
)

func StartCronJobs() {
	c := cron.New(cron.WithLocation(time.Local))

	// Jalankan setiap hari pada pukul 08:00
	_, err := c.AddFunc("0 8 * * *", processMaintenanceReminders)
	if err != nil {
		log.Printf("Gagal mendaftarkan cron job: %v", err)
	}

	c.Start()
	log.Println("Cron jobs started successfully")
}

func processMaintenanceReminders() {
	log.Println("Running processMaintenanceReminders...")

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	sevenDaysFromNow := startOfToday.AddDate(0, 0, 7)

	var schedules []maintenancemodels.MaintenanceSchedule
	// Mengambil jadwal yang belum selesai (mendatang) dan jatuh tempo dalam 7 hari ke depan (atau sudah lewat)
	if err := config.DB.Preload("User").Preload("PlantBatch").
		Where("status = ?", "mendatang").
		Where("next_due_date <= ?", sevenDaysFromNow).
		Find(&schedules).Error; err != nil {
		log.Printf("Error fetching maintenance schedules for reminders: %v", err)
		return
	}

	// Fetch Email Template
	var template notificationmodels.EmailTemplate
	hasTemplate := true
	if err := config.DB.Where("name = ?", "maintenance_reminder").First(&template).Error; err != nil {
		log.Printf("Template email 'maintenance_reminder' tidak ditemukan. Menggunakan fallback template. Error: %v", err)
		hasTemplate = false
	}

	for _, schedule := range schedules {
		var emailSubject string
		var emailBody string

		if hasTemplate && template.IsActive {
			// Resolve Subject
			emailSubject = template.Subject
			emailSubject = strings.ReplaceAll(emailSubject, "{{user_name}}", schedule.User.Name)
			emailSubject = strings.ReplaceAll(emailSubject, "{{activity_name}}", schedule.ActivityType)
			emailSubject = strings.ReplaceAll(emailSubject, "{{due_date}}", schedule.NextDueDate.Format("02 Jan 2006"))

			// Resolve Body
			emailBody = template.Body
			emailBody = strings.ReplaceAll(emailBody, "{{user_name}}", schedule.User.Name)
			emailBody = strings.ReplaceAll(emailBody, "{{activity_name}}", schedule.ActivityType)
			emailBody = strings.ReplaceAll(emailBody, "{{due_date}}", schedule.NextDueDate.Format("02 Jan 2006"))
		} else {
			emailSubject = "Pengingat Aktivitas Pemeliharaan Tanaman - Pondok Tani"
			emailBody = fmt.Sprintf(`
				<div style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
					<h3 style="color: #2e7d32;">Halo %s,</h3>
					<p>Ini adalah pengingat otomatis bahwa Anda memiliki aktivitas pemeliharaan tanaman yang akan datang.</p>
					<table style="border-collapse: collapse; width: 100%%; max-width: 600px; margin-top: 15px; margin-bottom: 15px;">
						<tr>
							<td style="padding: 8px; border: 1px solid #ddd; font-weight: bold; width: 35%%;">Aktivitas</td>
							<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
						</tr>
						<tr>
							<td style="padding: 8px; border: 1px solid #ddd; font-weight: bold;">Tanggal Jatuh Tempo</td>
							<td style="padding: 8px; border: 1px solid #ddd;">%s</td>
						</tr>
					</table>
					<p>Silakan login ke dashboard Pondok Tani Anda untuk melihat detail lebih lanjut.</p>
					<br>
					<p>Terima kasih,<br><strong>Tim Pondok Tani</strong></p>
				</div>
			`, schedule.User.Name, schedule.ActivityType, schedule.NextDueDate.Format("02 Jan 2006"))
		}

		// Buat versi plain text untuk bell notification
		plainTextMessage := strings.ReplaceAll(emailBody, "\n", " ")
		plainTextMessage = strings.ReplaceAll(plainTextMessage, "\r", "")
		// Sederhanakan strip HTML
		for strings.Contains(plainTextMessage, "<") && strings.Contains(plainTextMessage, ">") {
			start := strings.Index(plainTextMessage, "<")
			end := strings.Index(plainTextMessage, ">")
			if start < end {
				plainTextMessage = plainTextMessage[:start] + " " + plainTextMessage[end+1:]
			} else {
				break
			}
		}
		// Bersihkan spasi ganda
		plainTextMessage = strings.Join(strings.Fields(plainTextMessage), " ")

		// 1. Buat Bell Notification (In-App)
		notif := notificationmodels.Notification{
			UserID:   schedule.UserID,
			UserName: schedule.User.Name,
			Type:     "maintenance_reminder",
			Subject:  emailSubject,
			Message:  plainTextMessage, // Gunakan versi plain text
		}
		if err := config.DB.Create(&notif).Error; err != nil {
			log.Printf("Gagal membuat notifikasi bell untuk user %s: %v", schedule.UserID, err)
		}

		// 2. Kirim Email Notification
		if schedule.User.Email != "" {
			go func(email, subj, body string, userID string) {
				sendStatus := "sent"
				var errorMessage *string

				if err := notificationhandler.SendSMTPEmail(email, subj, body); err != nil {
					errMsg := err.Error()
					errorMessage = &errMsg
					sendStatus = "failed"
					log.Printf("Gagal mengirim email pengingat ke %s: %v", email, err)
				} else {
					log.Printf("Email pengingat berhasil dikirim ke %s", email)
				}

				now := time.Now().UTC()
				logEntry := notificationmodels.EmailLog{
					UserID:         &userID,
					RecipientEmail: email,
					Subject:        subj,
					Body:           body,
					Status:         sendStatus,
					ErrorMessage:   errorMessage,
					SentAt:         &now,
				}
				config.DB.Create(&logEntry)
			}(schedule.User.Email, emailSubject, emailBody, schedule.UserID)
		}
	}

	log.Println("Selesai menjalankan processMaintenanceReminders.")
}
