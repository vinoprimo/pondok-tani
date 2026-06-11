package notification

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"pondok-tani-backend/config"
	notificationmodels "pondok-tani-backend/models/notification"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createEmailTemplateRequest struct {
	Name     string `json:"name" binding:"required"`
	Subject  string `json:"subject" binding:"required"`
	Body     string `json:"body" binding:"required"`
	IsActive *bool  `json:"is_active"`
}

type updateEmailTemplateRequest struct {
	Name     *string `json:"name"`
	Subject  *string `json:"subject"`
	Body     *string `json:"body"`
	IsActive *bool   `json:"is_active"`
}

type sendEmailRequest struct {
	UserID         *string `json:"user_id"`
	RecipientEmail string  `json:"recipient_email" binding:"required,email"`
	TemplateName   string  `json:"template_name" binding:"required"`
	Subject        string  `json:"subject" binding:"required"`
	Body           string  `json:"body" binding:"required"`
}

type checkEmailDuplicateRequest struct {
	RecipientEmail string `json:"recipient_email" binding:"required,email"`
	TemplateName   string `json:"template_name" binding:"required"`
}

func CreateEmailTemplate(c *gin.Context) {
	var req createEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	template := notificationmodels.EmailTemplate{
		Name:     strings.TrimSpace(req.Name),
		Subject:  strings.TrimSpace(req.Subject),
		Body:     strings.TrimSpace(req.Body),
		IsActive: true,
	}

	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}

	if err := config.DB.Create(&template).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create email template"})
		return
	}

	c.JSON(http.StatusCreated, template)
}

func ListEmailTemplates(c *gin.Context) {
	var templates []notificationmodels.EmailTemplate
	query := config.DB.Order("created_at DESC")

	if err := query.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email templates"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

func UpdateEmailTemplate(c *gin.Context) {
	templateID := c.Param("id")
	var req updateEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = strings.TrimSpace(*req.Name)
	}
	if req.Subject != nil {
		updates["subject"] = strings.TrimSpace(*req.Subject)
	}
	if req.Body != nil {
		updates["body"] = strings.TrimSpace(*req.Body)
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No changes provided"})
		return
	}

	if err := config.DB.Model(&notificationmodels.EmailTemplate{}).Where("id = ?", templateID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email template updated"})
}

func DeleteEmailTemplate(c *gin.Context) {
	templateID := c.Param("id")
	if err := config.DB.Delete(&notificationmodels.EmailTemplate{}, "id = ?", templateID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete email template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email template deleted"})
}

func ListEmailLogs(c *gin.Context) {
	var logs []notificationmodels.EmailLog
	query := config.DB.Order("created_at DESC")

	if recipient := strings.TrimSpace(c.Query("recipient_email")); recipient != "" {
		query = query.Where("recipient_email = ?", recipient)
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Preload("Template").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func ListEmailMonitoring(c *gin.Context) {
	var items []notificationmodels.EmailMonitoring
	query := config.DB.Order("monitor_date DESC")

	if recipient := strings.TrimSpace(c.Query("recipient_email")); recipient != "" {
		query = query.Where("recipient_email = ?", recipient)
	}
	if template := strings.TrimSpace(c.Query("template_name")); template != "" {
		query = query.Where("template_name = ?", template)
	}

	if err := query.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email monitoring"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func CheckEmailDuplicate(c *gin.Context) {
	var req checkEmailDuplicateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	canSend, err := canSendEmail(req.RecipientEmail, req.TemplateName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check duplicate status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"duplicate": !canSend})
}

func SendEmailLog(c *gin.Context) {
	var req sendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	canSend, err := canSendEmail(req.RecipientEmail, req.TemplateName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate email delivery"})
		return
	}
	if !canSend {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terkirim untuk template ini pada hari yang sama"})
		return
	}

	now := time.Now().UTC()
	sendStatus := "sent"
	var errorMessage *string

	emailBody := strings.TrimSpace(req.Body)
	emailSubject := strings.TrimSpace(req.Subject)
	if err := SendSMTPEmail(req.RecipientEmail, emailSubject, emailBody); err != nil {
		errMsg := err.Error()
		sendStatus = "failed"
		errorMessage = &errMsg
	}

	logEntry := notificationmodels.EmailLog{
		UserID:         req.UserID,
		RecipientEmail: strings.TrimSpace(req.RecipientEmail),
		TemplateID:     nil,
		Subject:        emailSubject,
		Body:           emailBody,
		Status:         sendStatus,
		ErrorMessage:   errorMessage,
		SentAt:         &now,
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	if err := tx.Create(&logEntry).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email log"})
		return
	}

	if sendStatus == "sent" {
		if err := trackEmailMonitoring(tx, req.RecipientEmail, req.TemplateName, now); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email monitoring"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize email log"})
		return
	}

	if sendStatus == "failed" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email", "details": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent and logged", "log_id": logEntry.ID})
}

func SendSMTPEmail(recipientEmail, subject, body string) error {
	smtpHost := strings.TrimSpace(config.SMTP.Host)
	smtpPort := strings.TrimSpace(config.SMTP.Port)
	smtpUser := strings.TrimSpace(config.SMTP.Email)
	smtpPass := strings.TrimSpace(config.SMTP.Password)
	smtpDisplayName := strings.TrimSpace(config.SMTP.DisplayName)
	smtpEnableSSL := config.SMTP.EnableSSL
	timeout := config.SMTP.Timeout

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP configuration incomplete")
	}

	fromHeader := smtpUser
	if smtpDisplayName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", smtpDisplayName, smtpUser)
	}

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", fromHeader),
		fmt.Sprintf("To: %s", recipientEmail),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	dialer := &net.Dialer{Timeout: timeout}

	if smtpEnableSSL {
		if smtpPort == "465" {
			tlsConfig := &tls.Config{
				ServerName: smtpHost,
			}

			conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
			if err != nil {
				return err
			}
			defer conn.Close()

			client, err := smtp.NewClient(conn, smtpHost)
			if err != nil {
				return err
			}
			defer client.Quit()

			if err := client.Auth(auth); err != nil {
				return err
			}

			if err := client.Mail(smtpUser); err != nil {
				return err
			}
			if err := client.Rcpt(recipientEmail); err != nil {
				return err
			}

			wc, err := client.Data()
			if err != nil {
				return err
			}

			_, err = wc.Write([]byte(msg))
			if err != nil {
				wc.Close()
				return err
			}
			return wc.Close()
		}

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		client, err := smtp.NewClient(conn, smtpHost)
		if err != nil {
			return err
		}
		defer client.Quit()

		tlsConfig := &tls.Config{
			ServerName: smtpHost,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return err
		}

		if err := client.Auth(auth); err != nil {
			return err
		}

		if err := client.Mail(smtpUser); err != nil {
			return err
		}
		if err := client.Rcpt(recipientEmail); err != nil {
			return err
		}

		wc, err := client.Data()
		if err != nil {
			return err
		}

		_, err = wc.Write([]byte(msg))
		if err != nil {
			wc.Close()
			return err
		}
		return wc.Close()
	}

	return smtp.SendMail(addr, auth, smtpUser, []string{recipientEmail}, []byte(msg))
}

func canSendEmail(recipientEmail, templateName string) (bool, error) {
	dateOnly := time.Now().UTC().Truncate(24 * time.Hour)
	var existing notificationmodels.EmailMonitoring
	if err := config.DB.Where("recipient_email = ? AND template_name = ? AND monitor_date = ?", strings.TrimSpace(recipientEmail), strings.TrimSpace(templateName), dateOnly).First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, err
	}

	return existing.MessageCount == 0, nil
}

func trackEmailMonitoring(tx *gorm.DB, recipientEmail, templateName string, now time.Time) error {
	dateOnly := now.UTC().Truncate(24 * time.Hour)
	var monitoring notificationmodels.EmailMonitoring
	if err := tx.Where("recipient_email = ? AND template_name = ? AND monitor_date = ?", strings.TrimSpace(recipientEmail), strings.TrimSpace(templateName), dateOnly).First(&monitoring).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			monitoring = notificationmodels.EmailMonitoring{
				RecipientEmail: strings.TrimSpace(recipientEmail),
				TemplateName:   strings.TrimSpace(templateName),
				MonitorDate:    dateOnly,
				MessageCount:   1,
				FirstSentAt:    &now,
				LastSentAt:     &now,
			}
			return tx.Create(&monitoring).Error
		}
		return err
	}

	monitoring.MessageCount++
	monitoring.LastSentAt = &now
	return tx.Save(&monitoring).Error
}
