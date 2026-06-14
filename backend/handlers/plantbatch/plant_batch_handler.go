package plantbatch

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	maintenancemodels "pondok-tani-backend/models/maintenance"

	"github.com/gin-gonic/gin"
)

type plantBatchSummary struct {
	ID           uint      `json:"id"`
	InvestmentID uint      `json:"investment_id"`
	PackageID    uint      `json:"package_id"`
	PackageName  string    `json:"package_name"`
	BatchCode    string    `json:"batch_code"`
	PlantingDate time.Time `json:"planting_date"`
	Location     string    `json:"location"`
	SeedCount    uint      `json:"seed_count"`
	LandArea     float64   `json:"land_area"`
	PhotoURL     string    `json:"photo_url"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type activatePlantBatchItem struct {
	PackageID    uint    `json:"package_id" binding:"required,gt=0"`
	BatchCode    string  `json:"batch_code" binding:"required,min=3,max=60"`
	PlantingDate string  `json:"planting_date" binding:"required"`
	Location     string  `json:"location" binding:"required,min=3,max=150"`
	SeedCount    uint    `json:"seed_count" binding:"required,gt=0"`
	LandArea     float64 `json:"land_area" binding:"required,gt=0"`
	PhotoURL     string  `json:"photo_url"`
}

type activatePackageRequest struct {
	PlantBatches []activatePlantBatchItem `json:"plant_batches" binding:"required"`
}

var phaseProgressMap = map[string]int{
	"penanaman":        0,
	"pertumbuhan_awal": 15,
	"vegetatif":        40,
	"pembungaan":       70,
	"pembuahan":        85,
	"panen":            100,
}

func saveBatchImage(c *gin.Context, fieldKey string) (string, error) {
	fileHeader, err := c.FormFile(fieldKey)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return "", fmt.Errorf("unsupported file extension")
	}

	if err := os.MkdirAll("uploads/batches", 0o755); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("batch_%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join("uploads", "batches", fileName)
	if err := c.SaveUploadedFile(fileHeader, fullPath); err != nil {
		return "", err
	}

	return "/uploads/batches/" + fileName, nil
}

func getPhaseProgress(phase string) (int, bool) {
	progress, ok := phaseProgressMap[strings.ToLower(strings.TrimSpace(phase))]
	return progress, ok
}

func calculateAgeProgress(plantingDate time.Time) int {
	ageMonths := int(math.Floor(time.Since(plantingDate).Hours() / 24 / 30.4375))
	if ageMonths < 0 {
		ageMonths = 0
	}

	switch {
	case ageMonths < 3:
		return int(math.Round(float64(ageMonths) / 3 * 15))
	case ageMonths < 12:
		return int(math.Round(15 + float64(ageMonths-3)/9*25))
	case ageMonths < 24:
		return int(math.Round(40 + float64(ageMonths-12)/12*30))
	case ageMonths < 30:
		return int(math.Round(70 + float64(ageMonths-24)/6*15))
	case ageMonths < 36:
		return int(math.Round(85 + float64(ageMonths-30)/6*10))
	default:
		return 100
	}
}

func calculatePlantProgress(plantingDate time.Time, phaseKey string) int {
	ageProgress := calculateAgeProgress(plantingDate)
	phaseBase, exists := phaseProgressMap[phaseKey]
	if !exists {
		phaseBase = 0
	}

	progress := int(math.Max(float64(ageProgress), float64(phaseBase)))
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	return progress
}

func ActivateUserPackage(c *gin.Context) {
	userID := c.Param("id")

	var req activatePackageRequest
	isMultipart := strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "multipart/form-data")
	if isMultipart {
		payload := c.PostForm("plant_batches")
		if strings.TrimSpace(payload) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "plant_batches is required"})
			return
		}
		if err := json.Unmarshal([]byte(payload), &req.PlantBatches); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant batch form"})
			return
		}
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant batch form"})
			return
		}
	}
	if len(req.PlantBatches) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plant_batches is required"})
		return
	}

	createdBy, exists := c.Get("user_id")
	createdByStr, ok := createdBy.(string)
	if !exists || !ok || strings.TrimSpace(createdByStr) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	var user authmodels.User
	if err := config.DB.Preload("Investments").Preload("SelectedPackage").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if len(user.Investments) == 0 && user.SelectedPackageID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has not selected a package yet"})
		return
	}

	now := time.Now()
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	onProcessInvestments := make([]coremodels.Investment, 0)
	for _, investment := range user.Investments {
		if investment.Status == "on_process" {
			onProcessInvestments = append(onProcessInvestments, investment)
		}
	}

	if len(user.Investments) == 0 && user.SelectedPackageID != nil {
		var packageItem coremodels.InvestmentPackage
		if err := tx.First(&packageItem, "id = ? AND status = ?", *user.SelectedPackageID, "active").Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Investment package not found"})
			return
		}

		legacyInvestment := coremodels.Investment{
			UserID:         userID,
			PackageID:      *user.SelectedPackageID,
			Amount:         packageItem.Price,
			Status:         "on_process",
			InvestmentDate: now,
		}

		if err := tx.Create(&legacyInvestment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate package"})
			return
		}

		onProcessInvestments = append(onProcessInvestments, legacyInvestment)
	}

	if len(onProcessInvestments) == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "User must be in on_process status before activation"})
		return
	}

	if len(req.PlantBatches) != len(onProcessInvestments) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Each on_process package must have one plant batch form"})
		return
	}

	investmentIDs := make([]uint, 0, len(onProcessInvestments))
	onProcessByPackageID := make(map[uint]coremodels.Investment, len(onProcessInvestments))
	for _, investment := range onProcessInvestments {
		investmentIDs = append(investmentIDs, investment.ID)
		onProcessByPackageID[investment.PackageID] = investment
	}

	usedPackageIDs := map[uint]struct{}{}

	for _, item := range req.PlantBatches {
		_, exists := onProcessByPackageID[item.PackageID]
		if !exists {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Package %d is not in on_process status", item.PackageID)})
			return
		}

		if _, duplicate := usedPackageIDs[item.PackageID]; duplicate {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Duplicate plant batch for package %d", item.PackageID)})
			return
		}
		usedPackageIDs[item.PackageID] = struct{}{}
	}

	for _, item := range req.PlantBatches {
		investment := onProcessByPackageID[item.PackageID]

		plantingDate, err := time.Parse("2006-01-02", item.PlantingDate)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting_date format. Use YYYY-MM-DD"})
			return
		}

		batchCode := strings.TrimSpace(item.BatchCode)
		location := strings.TrimSpace(item.Location)
		if batchCode == "" || location == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "batch_code and location are required"})
			return
		}

		photoURL := strings.TrimSpace(item.PhotoURL)
		if isMultipart {
			fieldKey := fmt.Sprintf("batch_image_%d", item.PackageID)
			uploadedPhotoURL, err := saveBatchImage(c, fieldKey)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Photo upload failed for package %d", item.PackageID)})
				return
			}
			photoURL = uploadedPhotoURL
		}

		if photoURL == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Photo is required for package %d", item.PackageID)})
			return
		}

		plantBatch := coremodels.PlantBatch{
			InvestmentID: investment.ID,
			BatchCode:    batchCode,
			PlantingDate: plantingDate,
			Location:     location,
			SeedCount:    item.SeedCount,
			LandArea:     item.LandArea,
			PhotoURL:     photoURL,
			Status:       "active",
		}

		if err := tx.Create(&plantBatch).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save plant batch. Ensure batch_code is unique"})
			return
		}

		progress := calculatePlantProgress(plantingDate, "penanaman")
		initialNote := fmt.Sprintf("Validasi penanaman awal untuk batch %s", batchCode)
		creatorID := createdByStr
		initialMonitoring := maintenancemodels.PlantMonitoring{
			PlantBatchID:  plantBatch.ID,
			Phase:         "penanaman",
			Progress:      progress,
			HealthStatus:  "sehat",
			Disease:       "",
			DiseaseNote:   "",
			AffectedCount: 0,
			TotalPlants:   int(item.SeedCount),
			Note:          initialNote,
			CreatedBy:     &creatorID,
			PhotoURL:      photoURL,
		}

		if err := tx.Create(&initialMonitoring).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create initial plant monitoring"})
			return
		}
	}

	if err := tx.Model(&coremodels.Investment{}).
		Where("id IN ?", investmentIDs).
		Updates(map[string]interface{}{
			"status":     "active",
			"updated_at": now,
		}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate package"})
		return
	}

	updates := map[string]interface{}{
		"package_status":       "active",
		"package_activated_at": now,
	}

	if err := tx.Model(&authmodels.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate package"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate package"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Package activated after plant batch validation"})
}

func ListUserPlantBatches(c *gin.Context) {
	userID := c.Param("id")

	var user authmodels.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	subQuery := config.DB.Model(&coremodels.Investment{}).
		Select("id").
		Where("user_id = ?", userID)

	var plantBatches []coremodels.PlantBatch
	if err := config.DB.
		Preload("Investment").
		Preload("Investment.Package").
		Where("investment_id IN (?)", subQuery).
		Order("planting_date DESC").
		Find(&plantBatches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plant batch data"})
		return
	}

	response := make([]plantBatchSummary, 0, len(plantBatches))
	for _, item := range plantBatches {
		response = append(response, plantBatchSummary{
			ID:           item.ID,
			InvestmentID: item.InvestmentID,
			PackageID:    item.Investment.PackageID,
			PackageName:  item.Investment.Package.PackageName,
			BatchCode:    item.BatchCode,
			PlantingDate: item.PlantingDate,
			Location:     item.Location,
			SeedCount:    item.SeedCount,
			LandArea:     item.LandArea,
			PhotoURL:     item.PhotoURL,
			Status:       item.Status,
			CreatedAt:    item.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
