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
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	maintenancemodels "pondok-tani-backend/models/maintenance"

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
	HasSubmittedProof bool      `json:"has_submitted_proof"`
}

type myMaintenanceActivityResponse struct {
	ID               uint       `json:"id"`
	ScheduleID       uint       `json:"schedule_id"`
	ScheduleCode     string     `json:"schedule_code"`
	ActivityType     string     `json:"activity_type"`
	ActivityDate     time.Time  `json:"activity_date"`
	PhotoURL         *string    `json:"photo_url,omitempty"`
	ValidationStatus string     `json:"validation_status"`
	ValidationNotes  *string    `json:"validation_notes,omitempty"`
	ValidatedAt      *time.Time `json:"validated_at,omitempty"`
}

type submitMaintenanceActivityRequest struct {
	ScheduleID   uint    `json:"schedule_id" binding:"required,gt=0"`
	ActivityDate string  `json:"activity_date" binding:"required"`
	Description  *string `json:"description"`
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

	if err := query.Order("next_due_date ASC").Find(&schedules).Error; err != nil {
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
		Order("next_due_date ASC").
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
		if hasLatest {
			latestPtr = &latest
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
		Where("schedule_id IN (?)", config.DB.Model(&maintenancemodels.MaintenanceSchedule{}).Select("id").Where("user_id = ?", userID)).
		Order("activity_date DESC").
		Find(&activities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch maintenance activities"})
		return
	}

	response := make([]myMaintenanceActivityResponse, 0, len(activities))
	for _, item := range activities {
		var validatedAt *time.Time
		if item.ValidationStatus == "approved" || item.ValidationStatus == "verified" || item.ValidationStatus == "rejected" {
			t := item.UpdatedAt
			validatedAt = &t
		}

		response = append(response, myMaintenanceActivityResponse{
			ID:               item.ID,
			ScheduleID:       item.ScheduleID,
			ScheduleCode:     fmt.Sprintf("MA-%03d", item.ScheduleID),
			ActivityType:     item.ActivityType,
			ActivityDate:     item.ActivityDate,
			PhotoURL:         item.PhotoURL,
			ValidationStatus: item.ValidationStatus,
			ValidationNotes:  item.ValidationNotes,
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
