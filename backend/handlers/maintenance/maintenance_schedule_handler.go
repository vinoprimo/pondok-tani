package maintenance

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"pondok-tani-backend/config"
	notificationhandler "pondok-tani-backend/handlers/notification"
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	maintenancemodels "pondok-tani-backend/models/maintenance"
	notificationmodels "pondok-tani-backend/models/notification"
	"pondok-tani-backend/utils"

	"github.com/gin-gonic/gin"
)

func syncOverdueMaintenanceSchedules() error {
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	// Normalize legacy pending status to the new status naming.
	if err := config.DB.Exec(`
		UPDATE maintenance_schedules
		SET status = 'mendatang', updated_at = NOW()
		WHERE status = 'pending'
	`).Error; err != nil {
		return err
	}

	return config.DB.Exec(`
		UPDATE maintenance_schedules AS ms
		SET status = 'overdue', updated_at = NOW()
		WHERE ms.status = 'mendatang'
		  AND DATE(ms.next_due_date) < ?
		  AND NOT EXISTS (
			SELECT 1
			FROM maintenance_activities ma
			WHERE ma.schedule_id = ms.id
		  )
	`, startOfDay).Error
}

type createMaintenanceScheduleRequest struct {
	UserID        string `json:"user_id" binding:"required,uuid"`
	PlantBatchID  uint   `json:"plant_batch_id" binding:"required,gt=0"`
	ActivityType  string `json:"activity_type" binding:"required,min=2,max=60"`
	FrequencyDays uint   `json:"frequency_days"`
	NextDueDate   string `json:"next_due_date" binding:"required"`
}

type maintenanceScheduleResponse struct {
	ID            uint      `json:"id"`
	UserID        string    `json:"user_id"`
	UserName      string    `json:"user_name"`
	UserEmail     string    `json:"user_email"`
	PlantBatchID  uint      `json:"plant_batch_id"`
	BatchCode     string    `json:"batch_code"`
	ActivityType  string    `json:"activity_type"`
	FrequencyDays uint      `json:"frequency_days"`
	NextDueDate   time.Time `json:"next_due_date"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type maintenanceScheduleSummaryResponse struct {
	PendingCount  int64 `json:"pending_count"`
	DueTodayCount int64 `json:"due_today_count"`
	OverdueCount  int64 `json:"overdue_count"`
	TotalCount    int64 `json:"total_count"`
}

type myMaintenanceScheduleResponse struct {
	ID                uint      `json:"id"`
	ScheduleCode      string    `json:"schedule_code"`
	PlantBatchID      uint      `json:"plant_batch_id"`
	BatchCode         string    `json:"batch_code"`
	ActivityType      string    `json:"activity_type"`
	FrequencyDays     uint      `json:"frequency_days"`
	NextDueDate       time.Time `json:"next_due_date"`
	Status            string    `json:"status"`
	Remark            *string   `json:"remark,omitempty"`
	HasSubmittedProof bool      `json:"has_submitted_proof"`
}

type myMaintenanceActivityResponse struct {
	ID               uint       `json:"id"`
	ScheduleID       uint       `json:"schedule_id"`
	ScheduleCode     string     `json:"schedule_code"`
	BatchCode        string     `json:"batch_code"`
	ActivityType     string     `json:"activity_type"`
	Description      *string    `json:"description,omitempty"`
	ActivityDate     time.Time  `json:"activity_date"`
	PhotoURL         *string    `json:"photo_url,omitempty"`
	ValidationStatus string     `json:"validation_status"`
	ValidationNotes  *string    `json:"validation_notes,omitempty"`
	Remark           *string    `json:"remark,omitempty"`
	ValidatedAt      *time.Time `json:"validated_at,omitempty"`
}

type submitMaintenanceActivityRequest struct {
	ScheduleID   uint    `json:"schedule_id" binding:"required,gt=0"`
	ActivityDate string  `json:"activity_date" binding:"required"`
	Description  *string `json:"description"`
}

type reviewMaintenanceActivityRequest struct {
	Action string  `json:"action" binding:"required"`
	Notes  *string `json:"notes"`
}

type adminMaintenanceActivitySummaryResponse struct {
	WaitingCount  int64 `json:"waiting_count"`
	ApprovedToday int64 `json:"approved_today"`
	RejectedToday int64 `json:"rejected_today"`
	ReviewedMonth int64 `json:"reviewed_month"`
}

type adminMaintenanceActivityItem struct {
	ID             uint      `json:"id"`
	ActivityCode   string    `json:"activity_code"`
	UserID         string    `json:"user_id"`
	UserName       string    `json:"user_name"`
	BatchCode      string    `json:"batch_code"`
	ActivityType   string    `json:"activity_type"`
	ActivityDate   time.Time `json:"activity_date"`
	SubmissionDate time.Time `json:"submission_date"`
	Description    *string   `json:"description,omitempty"`
	PhotoURL       *string   `json:"photo_url,omitempty"`
	Status         string    `json:"status"`
}

type adminMaintenanceActivityRaw struct {
	ID               uint
	ScheduleID       uint
	ActivityType     string
	ActivityDate     time.Time
	Description      *string
	PhotoURL         *string
	ValidationStatus string
	ValidationNotes  *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UserID           string
	UserName         string
	BatchCode        string
}

func saveMaintenanceActivityImage(c *gin.Context, fieldKey string) (string, error) {
	fileHeader, err := c.FormFile(fieldKey)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return "", errors.New("unsupported file extension")
	}

	if err := os.MkdirAll("uploads/maintenance-activities", 0o755); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("maintenance_activity_%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join("uploads", "maintenance-activities", fileName)
	if err := c.SaveUploadedFile(fileHeader, fullPath); err != nil {
		return "", err
	}

	return "/uploads/maintenance-activities/" + fileName, nil
}

func parseSubmitActivityInput(c *gin.Context) (submitMaintenanceActivityRequest, bool, error) {
	var req submitMaintenanceActivityRequest
	isMultipart := strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "multipart/form-data")

	if isMultipart {
		scheduleIDValue := strings.TrimSpace(c.PostForm("schedule_id"))
		if scheduleIDValue == "" {
			return req, true, errors.New("schedule_id is required")
		}
		scheduleIDParsed, err := strconv.ParseUint(scheduleIDValue, 10, 64)
		if err != nil || scheduleIDParsed == 0 {
			return req, true, errors.New("invalid schedule_id")
		}

		req.ScheduleID = uint(scheduleIDParsed)
		req.ActivityDate = strings.TrimSpace(c.PostForm("activity_date"))
		description := strings.TrimSpace(c.PostForm("description"))
		if description != "" {
			req.Description = &description
		}

		if req.ActivityDate == "" {
			return req, true, errors.New("activity_date is required")
		}

		return req, true, nil
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return req, false, err
	}

	return req, false, nil
}

func deriveInvestorScheduleStatus(scheduleStatus string, latestActivity *maintenancemodels.MaintenanceActivity) string {
	if latestActivity != nil {
		switch latestActivity.ValidationStatus {
		case "approved", "verified":
			return "terverifikasi"
		case "rejected":
			return "ditolak"
		default:
			return "menunggu_verifikasi"
		}
	}

	if scheduleStatus == "overdue" {
		return "overdue"
	}

	return "mendatang"
}

func GetAdminMaintenanceActivitySummary(c *gin.Context) {
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	startOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	var waitingCount int64
	if err := config.DB.Model(&maintenancemodels.MaintenanceActivity{}).
		Where("validation_status = ?", "pending").
		Count(&waitingCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count waiting activities"})
		return
	}

	var approvedToday int64
	if err := config.DB.Model(&maintenancemodels.MaintenanceActivity{}).
		Where("validation_status IN ? AND updated_at >= ?", []string{"approved", "verified"}, startOfDay).
		Count(&approvedToday).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count approved activities"})
		return
	}

	var rejectedToday int64
	if err := config.DB.Model(&maintenancemodels.MaintenanceActivity{}).
		Where("validation_status = ? AND updated_at >= ?", "rejected", startOfDay).
		Count(&rejectedToday).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count rejected activities"})
		return
	}

	var reviewedMonth int64
	if err := config.DB.Model(&maintenancemodels.MaintenanceActivity{}).
		Where("validation_status IN ? AND updated_at >= ?", []string{"approved", "verified", "rejected"}, startOfMonth).
		Count(&reviewedMonth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count monthly reviewed activities"})
		return
	}

	c.JSON(http.StatusOK, adminMaintenanceActivitySummaryResponse{
		WaitingCount:  waitingCount,
		ApprovedToday: approvedToday,
		RejectedToday: rejectedToday,
		ReviewedMonth: reviewedMonth,
	})
}

func ListAdminMaintenanceActivities(c *gin.Context) {
	status := strings.ToLower(strings.TrimSpace(c.Query("status")))
	if status == "" {
		status = "pending"
	}

	pagination := utils.GeneratePaginationFromRequest(c)

	query := config.DB.Table("maintenance_activities AS ma").
		Select(`
			ma.id,
			ma.schedule_id,
			ma.activity_type,
			ma.activity_date,
			ma.description,
			ma.photo_url,
			ma.validation_status,
			ma.validation_notes,
			ma.created_at,
			ma.updated_at,
			ms.user_id,
			u.name AS user_name,
			pb.batch_code
		`).
		Joins("JOIN maintenance_schedules ms ON ms.id = ma.schedule_id").
		Joins("JOIN users u ON u.id = ms.user_id").
		Joins("JOIN plant_batches pb ON pb.id = ma.plant_batch_id")

	if status != "all" {
		query = query.Where("ma.validation_status = ?", status)
	}

	if pagination.Search != "" {
		searchPattern := "%" + pagination.Search + "%"
		query = query.Where("u.name ILIKE ? OR pb.batch_code ILIKE ? OR ma.activity_type ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// Override default sort for specific column name
	sortOrder := pagination.Sort
	if sortOrder == "created_at desc" {
		sortOrder = "ma.created_at DESC"
	}

	query = query.Order(sortOrder)

	var rows []adminMaintenanceActivityRaw
	if err := utils.Paginate(query, &pagination, &rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenance activities"})
		return
	}

	response := make([]adminMaintenanceActivityItem, 0, len(rows))
	for _, row := range rows {
		statusText := "Menunggu"
		switch row.ValidationStatus {
		case "approved", "verified":
			statusText = "Disetujui"
		case "rejected":
			statusText = "Ditolak"
		}

		response = append(response, adminMaintenanceActivityItem{
			ID:             row.ID,
			ActivityCode:   fmt.Sprintf("MA-%03d", row.ID),
			UserID:         row.UserID,
			UserName:       row.UserName,
			BatchCode:      row.BatchCode,
			ActivityType:   row.ActivityType,
			ActivityDate:   row.ActivityDate,
			SubmissionDate: row.CreatedAt,
			Description:    row.Description,
			PhotoURL:       row.PhotoURL,
			Status:         statusText,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
		"meta": map[string]interface{}{
			"limit":       pagination.GetLimit(),
			"page":        pagination.GetPage(),
			"sort":        pagination.GetSort(),
			"total_rows":  pagination.TotalRows,
			"total_pages": pagination.TotalPages,
		},
	})
}

func ReviewMaintenanceActivity(c *gin.Context) {
	activityID := c.Param("id")
	adminIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	adminID := adminIDValue.(string)

	var req reviewMaintenanceActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	action := strings.ToLower(strings.TrimSpace(req.Action))
	if action != "approve" && action != "reject" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action. Use approve or reject"})
		return
	}

	var activity maintenancemodels.MaintenanceActivity
	if err := config.DB.First(&activity, "id = ?", activityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance activity not found"})
		return
	}
	if activity.ValidationStatus != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Activity already reviewed"})
		return
	}

	var schedule maintenancemodels.MaintenanceSchedule
	if err := config.DB.First(&schedule, "id = ?", activity.ScheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance schedule not found"})
		return
	}

	newValidationStatus := "approved"
	newScheduleStatus := "terverifikasi"
	if action == "reject" {
		newValidationStatus = "rejected"
		newScheduleStatus = "ditolak"
	}

	notes := strings.TrimSpace(func() string {
		if req.Notes == nil {
			return ""
		}
		return *req.Notes
	}())

	if action == "reject" && notes == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alasan penolakan wajib diisi"})
		return
	}
	if action == "approve" && notes == "" {
		notes = "Disetujui oleh admin"
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	activityUpdates := map[string]interface{}{
		"validation_status": newValidationStatus,
		"validated_by":      adminID,
		"validation_notes":  notes,
		"updated_at":        time.Now(),
	}

	if err := tx.Model(&maintenancemodels.MaintenanceActivity{}).
		Where("id = ?", activity.ID).
		Updates(activityUpdates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance activity"})
		return
	}

	if err := tx.Model(&maintenancemodels.MaintenanceSchedule{}).
		Where("id = ?", activity.ScheduleID).
		Updates(map[string]interface{}{"status": newScheduleStatus, "updated_at": time.Now()}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance schedule"})
		return
	}

	if action == "approve" && schedule.FrequencyDays > 0 {
		nextDueDate := schedule.NextDueDate.AddDate(0, 0, int(schedule.FrequencyDays))

		var existingCount int64
		if err := tx.Model(&maintenancemodels.MaintenanceSchedule{}).
			Where("user_id = ? AND plant_batch_id = ? AND activity_type = ? AND next_due_date = ?", schedule.UserID, schedule.PlantBatchID, schedule.ActivityType, nextDueDate).
			Count(&existingCount).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate recurring schedule"})
			return
		}

		if existingCount == 0 {
			nextSchedule := maintenancemodels.MaintenanceSchedule{
				UserID:        schedule.UserID,
				PlantBatchID:  schedule.PlantBatchID,
				ActivityType:  schedule.ActivityType,
				FrequencyDays: schedule.FrequencyDays,
				NextDueDate:   nextDueDate,
				Status:        "mendatang",
			}

			if err := tx.Create(&nextSchedule).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recurring schedule"})
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance activity reviewed"})
}

func GetMaintenanceScheduleSummary(c *gin.Context) {
	if err := syncOverdueMaintenanceSchedules(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync overdue maintenance schedules"})
		return
	}

	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	todayKey := startOfDay.Format("2006-01-02")

	var totalCount int64
	if err := config.DB.Model(&maintenancemodels.MaintenanceSchedule{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count maintenance schedules"})
		return
	}

	var pendingCount int64
	if err := config.DB.Model(&maintenancemodels.MaintenanceSchedule{}).
		Where("status = ? AND DATE(next_due_date) >= ?", "mendatang", todayKey).
		Count(&pendingCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count pending schedules"})
		return
	}

	var dueTodayCount int64
	if err := config.DB.Model(&maintenancemodels.MaintenanceSchedule{}).
		Where("DATE(next_due_date) = ?", todayKey).
		Count(&dueTodayCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count today's schedules"})
		return
	}

	var overdueCount int64
	if err := config.DB.Model(&maintenancemodels.MaintenanceSchedule{}).
		Where("status = ? OR (status = ? AND DATE(next_due_date) < ?)", "overdue", "mendatang", todayKey).
		Count(&overdueCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count overdue schedules"})
		return
	}

	c.JSON(http.StatusOK, maintenanceScheduleSummaryResponse{
		PendingCount:  pendingCount,
		DueTodayCount: dueTodayCount,
		OverdueCount:  overdueCount,
		TotalCount:    totalCount,
	})
}

func ListMaintenanceSchedules(c *gin.Context) {
	if err := syncOverdueMaintenanceSchedules(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync overdue maintenance schedules"})
		return
	}

	var schedules []maintenancemodels.MaintenanceSchedule
	query := config.DB.Preload("User").Preload("PlantBatch")

	if userID := strings.TrimSpace(c.Query("user_id")); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("status = ?", strings.ToLower(status))
	}

	if err := query.Order("created_at DESC").Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenance schedules"})
		return
	}

	response := make([]maintenanceScheduleResponse, 0, len(schedules))
	for _, item := range schedules {
		response = append(response, maintenanceScheduleResponse{
			ID:            item.ID,
			UserID:        item.UserID,
			UserName:      item.User.Name,
			UserEmail:     item.User.Email,
			PlantBatchID:  item.PlantBatchID,
			BatchCode:     item.PlantBatch.BatchCode,
			ActivityType:  item.ActivityType,
			FrequencyDays: item.FrequencyDays,
			NextDueDate:   item.NextDueDate,
			Status:        item.Status,
			CreatedAt:     item.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func CreateMaintenanceSchedule(c *gin.Context) {
	var req createMaintenanceScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	nextDueDate, err := time.Parse("2006-01-02", req.NextDueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid next_due_date format. Use YYYY-MM-DD"})
		return
	}

	var user authmodels.User
	if err := config.DB.First(&user, "id = ?", req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var total int64
	if err := config.DB.Model(&coremodels.PlantBatch{}).
		Joins("JOIN investments ON investments.id = plant_batches.investment_id").
		Where("plant_batches.id = ? AND investments.user_id = ?", req.PlantBatchID, req.UserID).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate plant batch"})
		return
	}
	if total == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plant batch does not belong to selected user"})
		return
	}

	schedule := maintenancemodels.MaintenanceSchedule{
		UserID:        req.UserID,
		PlantBatchID:  req.PlantBatchID,
		ActivityType:  strings.TrimSpace(req.ActivityType),
		FrequencyDays: req.FrequencyDays,
		NextDueDate:   nextDueDate,
		Status:        "mendatang",
	}

	if err := config.DB.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create maintenance schedule"})
		return
	}

	// === TRIGGER NOTIFIKASI LANGSUNG ===
	go func(sch maintenancemodels.MaintenanceSchedule, u authmodels.User) {
		var template notificationmodels.EmailTemplate
		hasTemplate := true
		if err := config.DB.Where("name = ?", "maintenance_reminder").First(&template).Error; err != nil {
			hasTemplate = false
		}

		var emailSubject, emailBody string

		if hasTemplate && template.IsActive {
			emailSubject = template.Subject
			emailSubject = strings.ReplaceAll(emailSubject, "{{user_name}}", u.Name)
			emailSubject = strings.ReplaceAll(emailSubject, "{{activity_name}}", sch.ActivityType)
			emailSubject = strings.ReplaceAll(emailSubject, "{{due_date}}", sch.NextDueDate.Format("02 Jan 2006"))

			emailBody = template.Body
			emailBody = strings.ReplaceAll(emailBody, "{{user_name}}", u.Name)
			emailBody = strings.ReplaceAll(emailBody, "{{activity_name}}", sch.ActivityType)
			emailBody = strings.ReplaceAll(emailBody, "{{due_date}}", sch.NextDueDate.Format("02 Jan 2006"))
		} else {
			emailSubject = "Jadwal Pemeliharaan Baru - Pondok Tani"
			emailBody = fmt.Sprintf(`
				<div style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
					<h3 style="color: #2e7d32;">Halo %s,</h3>
					<p>Jadwal pemeliharaan tanaman baru telah ditambahkan untuk Anda.</p>
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
					<p>Silakan periksa dashboard Pondok Tani Anda.</p>
					<br>
					<p>Terima kasih,<br><strong>Tim Pondok Tani</strong></p>
				</div>
			`, u.Name, sch.ActivityType, sch.NextDueDate.Format("02 Jan 2006"))
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

		notif := notificationmodels.Notification{
			UserID:   u.ID,
			UserName: u.Name,
			Type:     "maintenance_reminder",
			Subject:  emailSubject,
			Message:  plainTextMessage,
		}
		config.DB.Create(&notif)

		if u.Email != "" {
			sendStatus := "sent"
			var errorMessage *string

			err := notificationhandler.SendSMTPEmail(u.Email, emailSubject, emailBody)
			if err != nil {
				errMsg := err.Error()
				errorMessage = &errMsg
				sendStatus = "failed"
				fmt.Printf("Gagal mengirim email notifikasi ke %s: %v\n", u.Email, err)
			} else {
				fmt.Printf("Email notifikasi berhasil dikirim ke %s\n", u.Email)
			}

			now := time.Now().UTC()
			logEntry := notificationmodels.EmailLog{
				UserID:         &u.ID,
				RecipientEmail: u.Email,
				Subject:        emailSubject,
				Body:           emailBody,
				Status:         sendStatus,
				ErrorMessage:   errorMessage,
				SentAt:         &now,
			}
			config.DB.Create(&logEntry)
		} else {
			fmt.Println("User tidak memiliki email, skip pengiriman email")
		}
	}(schedule, user)
	// ===================================

	c.JSON(http.StatusCreated, gin.H{"message": "Maintenance schedule created", "id": schedule.ID})
}

func ListMyMaintenanceSchedules(c *gin.Context) {
	if err := syncOverdueMaintenanceSchedules(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync overdue maintenance schedules"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var schedules []maintenancemodels.MaintenanceSchedule
	if err := config.DB.Preload("PlantBatch").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenance schedules"})
		return
	}

	scheduleIDs := make([]uint, 0, len(schedules))
	for _, item := range schedules {
		scheduleIDs = append(scheduleIDs, item.ID)
	}

	latestActivityMap := map[uint]maintenancemodels.MaintenanceActivity{}
	if len(scheduleIDs) > 0 {
		var activities []maintenancemodels.MaintenanceActivity
		if err := config.DB.
			Where("schedule_id IN ?", scheduleIDs).
			Order("schedule_id ASC, created_at DESC").
			Find(&activities).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenance activities"})
			return
		}

		for _, activity := range activities {
			if _, exists := latestActivityMap[activity.ScheduleID]; !exists {
				latestActivityMap[activity.ScheduleID] = activity
			}
		}
	}

	response := make([]myMaintenanceScheduleResponse, 0, len(schedules))
	for _, item := range schedules {
		latest, hasLatest := latestActivityMap[item.ID]
		var latestPtr *maintenancemodels.MaintenanceActivity
		var remark *string
		if hasLatest {
			latestPtr = &latest
			if latest.ValidationNotes != nil && strings.TrimSpace(*latest.ValidationNotes) != "" {
				remark = latest.ValidationNotes
			} else if latest.ValidationStatus == "approved" || latest.ValidationStatus == "verified" {
				defaultRemark := "Disetujui oleh admin"
				remark = &defaultRemark
			}
		}

		response = append(response, myMaintenanceScheduleResponse{
			ID:                item.ID,
			ScheduleCode:      fmt.Sprintf("MA-%03d", item.ID),
			PlantBatchID:      item.PlantBatchID,
			BatchCode:         item.PlantBatch.BatchCode,
			ActivityType:      item.ActivityType,
			FrequencyDays:     item.FrequencyDays,
			NextDueDate:       item.NextDueDate,
			Status:            deriveInvestorScheduleStatus(item.Status, latestPtr),
			Remark:            remark,
			HasSubmittedProof: hasLatest && latest.ValidationStatus == "pending",
		})
	}

	c.JSON(http.StatusOK, response)
}

func ListMyMaintenanceActivities(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var activities []maintenancemodels.MaintenanceActivity
	if err := config.DB.
		Preload("Schedule").
		Preload("PlantBatch").
		Where("schedule_id IN (?)", config.DB.Model(&maintenancemodels.MaintenanceSchedule{}).Select("id").Where("user_id = ?", userID)).
		Order("activity_date DESC").
		Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenance activities"})
		return
	}

	response := make([]myMaintenanceActivityResponse, 0, len(activities))
	for _, item := range activities {
		var validatedAt *time.Time
		var remark *string
		if item.ValidationStatus == "approved" || item.ValidationStatus == "verified" || item.ValidationStatus == "rejected" {
			t := item.UpdatedAt
			validatedAt = &t
		}
		if item.ValidationNotes != nil && strings.TrimSpace(*item.ValidationNotes) != "" {
			remark = item.ValidationNotes
		} else if item.ValidationStatus == "approved" || item.ValidationStatus == "verified" {
			defaultRemark := "Disetujui oleh admin"
			remark = &defaultRemark
		}

		response = append(response, myMaintenanceActivityResponse{
			ID:               item.ID,
			ScheduleID:       item.ScheduleID,
			ScheduleCode:     fmt.Sprintf("MA-%03d", item.ScheduleID),
			BatchCode:        item.PlantBatch.BatchCode,
			ActivityType:     item.ActivityType,
			Description:      item.Description,
			ActivityDate:     item.ActivityDate,
			PhotoURL:         item.PhotoURL,
			ValidationStatus: item.ValidationStatus,
			ValidationNotes:  item.ValidationNotes,
			Remark:           remark,
			ValidatedAt:      validatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func SubmitMyMaintenanceActivity(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDValue.(string)

	req, isMultipart, err := parseSubmitActivityInput(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	activityDate, err := time.Parse("2006-01-02", req.ActivityDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity_date format. Use YYYY-MM-DD"})
		return
	}

	var schedule maintenancemodels.MaintenanceSchedule
	if err := config.DB.First(&schedule, "id = ? AND user_id = ?", req.ScheduleID, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance schedule not found"})
		return
	}

	var pendingSubmissionCount int64
	if err := config.DB.Model(&maintenancemodels.MaintenanceActivity{}).
		Where("schedule_id = ? AND validation_status = ?", req.ScheduleID, "pending").
		Count(&pendingSubmissionCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate existing activity"})
		return
	}
	if pendingSubmissionCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Laporan untuk jadwal ini sedang menunggu verifikasi"})
		return
	}

	photoURL := ""
	if isMultipart {
		uploadedPhotoURL, uploadErr := saveMaintenanceActivityImage(c, "photo")
		if uploadErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Photo upload failed"})
			return
		}
		photoURL = uploadedPhotoURL
	}
	if strings.TrimSpace(photoURL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Photo is required"})
		return
	}

	photoURLValue := photoURL

	activity := maintenancemodels.MaintenanceActivity{
		ScheduleID:       schedule.ID,
		PlantBatchID:     schedule.PlantBatchID,
		ActivityType:     schedule.ActivityType,
		Description:      req.Description,
		PhotoURL:         &photoURLValue,
		ActivityDate:     activityDate,
		ValidationStatus: "pending",
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	if err := tx.Create(&activity).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit maintenance activity"})
		return
	}

	if err := tx.Model(&maintenancemodels.MaintenanceSchedule{}).
		Where("id = ?", schedule.ID).
		Updates(map[string]interface{}{"status": "menunggu_verifikasi", "updated_at": time.Now()}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule status"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save maintenance activity"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Laporan aktivitas berhasil dikirim"})
}
