package maintenance

import (
	"net/http"
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

	return config.DB.Exec(`
		UPDATE maintenance_schedules AS ms
		SET status = 'overdue', updated_at = NOW()
		WHERE ms.status = 'pending'
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
		Where("status = ? AND DATE(next_due_date) >= ?", "pending", todayKey).
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
		Where("status = ? OR (status = ? AND DATE(next_due_date) < ?)", "overdue", "pending", todayKey).
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
		Status:        "pending",
	}

	if err := config.DB.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create maintenance schedule"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Maintenance schedule created", "id": schedule.ID})
}
