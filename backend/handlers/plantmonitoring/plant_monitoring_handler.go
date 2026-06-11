package plantmonitoring

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
	coremodels "pondok-tani-backend/models/core"
	maintenancemodels "pondok-tani-backend/models/maintenance"

	"github.com/gin-gonic/gin"
)

var phaseProgressMap = map[string]int{
	"penanaman":        0,
	"pertumbuhan_awal": 20,
	"vegetatif":        40,
	"pra-berbunga":     60,
	"berbunga":         80,
	"panen":            100,
}

var phaseLabelToKey = map[string]string{
	"penanaman":        "penanaman",
	"pertumbuhan awal": "pertumbuhan_awal",
	"pertumbuhan_awal": "pertumbuhan_awal",
	"vegetatif":        "vegetatif",
	"pra-berbunga":     "pra-berbunga",
	"pra berbunga":     "pra-berbunga",
	"berbunga":         "berbunga",
	"panen":            "panen",
}

var healthStatusMap = map[string]string{
	"sehat":              "sehat",
	"sebagian terdampak": "sebagian_terdampak",
	"sebagian_terdampak": "sebagian_terdampak",
	"mati":               "mati",
}

var diseaseMap = map[string]string{
	"":             "",
	"null":         "",
	"tidak_ada":    "",
	"tidak ada":    "",
	"batang busuk": "batang_busuk",
	"batang_busuk": "batang_busuk",
	"lainnya":      "lainnya",
}

type createPlantMonitoringInput struct {
	PlantBatchID  uint
	Phase         string
	HealthStatus  string
	Disease       string
	DiseaseNote   string
	AffectedCount int
	TotalPlants   int
	Note          string
}

type plantMonitoringResponse struct {
	ID            uint      `json:"id"`
	PlantBatchID  uint      `json:"plant_batch_id"`
	BatchCode     string    `json:"batch_code"`
	PackageName   string    `json:"package_name"`
	Phase         string    `json:"phase"`
	Progress      int       `json:"progress"`
	HealthStatus  string    `json:"health_status"`
	Disease       string    `json:"disease"`
	DiseaseNote   string    `json:"disease_note"`
	AffectedCount int       `json:"affected_count"`
	TotalPlants   int       `json:"total_plants"`
	Note          string    `json:"note"`
	CreatedBy     string    `json:"created_by"`
	PhotoURL      string    `json:"photo_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func normalizePhase(input string) (string, int, error) {
	phaseKey, ok := phaseLabelToKey[strings.ToLower(strings.TrimSpace(input))]
	if !ok {
		return "", 0, errors.New("invalid phase")
	}

	progress, exists := phaseProgressMap[phaseKey]
	if !exists {
		return "", 0, errors.New("phase progress not configured")
	}

	return phaseKey, progress, nil
}

func normalizeHealthStatus(input string) (string, error) {
	status, ok := healthStatusMap[strings.ToLower(strings.TrimSpace(input))]
	if !ok {
		return "", errors.New("invalid health_status")
	}
	return status, nil
}

func normalizeDisease(input string) (string, error) {
	status, ok := diseaseMap[strings.ToLower(strings.TrimSpace(input))]
	if !ok {
		return "", errors.New("invalid disease")
	}
	return status, nil
}

func savePlantMonitoringImage(c *gin.Context, fieldKey string) (string, error) {
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

	if err := os.MkdirAll("uploads/plant-monitorings", 0o755); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("monitoring_%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join("uploads", "plant-monitorings", fileName)
	if err := c.SaveUploadedFile(fileHeader, fullPath); err != nil {
		return "", err
	}

	return "/uploads/plant-monitorings/" + fileName, nil
}

func parseCreateInput(c *gin.Context) (createPlantMonitoringInput, bool, error) {
	var input createPlantMonitoringInput
	isMultipart := strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "multipart/form-data")

	if isMultipart {
		batchIDValue := strings.TrimSpace(c.PostForm("plant_batch_id"))
		if batchIDValue == "" {
			return input, true, errors.New("plant_batch_id is required")
		}
		batchIDParsed, err := strconv.ParseUint(batchIDValue, 10, 64)
		if err != nil || batchIDParsed == 0 {
			return input, true, errors.New("invalid plant_batch_id")
		}
		input.PlantBatchID = uint(batchIDParsed)
		input.Phase = c.PostForm("phase")
		input.HealthStatus = c.PostForm("health_status")
		input.Disease = c.PostForm("disease")
		input.DiseaseNote = c.PostForm("disease_note")
		input.Note = c.PostForm("note")

		affectedCount := strings.TrimSpace(c.PostForm("affected_count"))
		if affectedCount != "" {
			parsedAffected, err := strconv.Atoi(affectedCount)
			if err != nil {
				return input, true, errors.New("invalid affected_count")
			}
			input.AffectedCount = parsedAffected
		}

		totalPlants := strings.TrimSpace(c.PostForm("total_plants"))
		if totalPlants != "" {
			parsedTotal, err := strconv.Atoi(totalPlants)
			if err != nil {
				return input, true, errors.New("invalid total_plants")
			}
			input.TotalPlants = parsedTotal
		}

		return input, true, nil
	}

	var jsonPayload struct {
		PlantBatchID  uint   `json:"plant_batch_id"`
		Phase         string `json:"phase"`
		HealthStatus  string `json:"health_status"`
		Disease       string `json:"disease"`
		DiseaseNote   string `json:"disease_note"`
		AffectedCount int    `json:"affected_count"`
		TotalPlants   int    `json:"total_plants"`
		Note          string `json:"note"`
	}
	if err := c.ShouldBindJSON(&jsonPayload); err != nil {
		return input, false, err
	}

	input = createPlantMonitoringInput{
		PlantBatchID:  jsonPayload.PlantBatchID,
		Phase:         jsonPayload.Phase,
		HealthStatus:  jsonPayload.HealthStatus,
		Disease:       jsonPayload.Disease,
		DiseaseNote:   jsonPayload.DiseaseNote,
		AffectedCount: jsonPayload.AffectedCount,
		TotalPlants:   jsonPayload.TotalPlants,
		Note:          jsonPayload.Note,
	}

	return input, false, nil
}

func CreatePlantMonitoring(c *gin.Context) {
	input, isMultipart, err := parseCreateInput(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form payload"})
		return
	}
	if input.PlantBatchID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plant_batch_id is required"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	userID, ok := userIDValue.(string)
	if !exists || !ok || strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		return
	}

	phaseKey, progress, err := normalizePhase(input.Phase)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phase"})
		return
	}

	healthStatus, err := normalizeHealthStatus(input.HealthStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid health_status"})
		return
	}

	disease, err := normalizeDisease(input.Disease)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid disease"})
		return
	}

	if input.AffectedCount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "affected_count cannot be negative"})
		return
	}

	var plantBatch coremodels.PlantBatch
	query := config.DB.Preload("Investment").Preload("Investment.Package").Where("plant_batches.id = ?", input.PlantBatchID)
	if roleValue, ok := c.Get("role"); ok {
		role, _ := roleValue.(string)
		if role == "investor" || role == "mitra" {
			query = query.Joins("JOIN investments ON investments.id = plant_batches.investment_id").Where("investments.user_id = ?", userID)
		}
	}
	if err := query.First(&plantBatch).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plant batch not found"})
		return
	}

	totalPlants := input.TotalPlants
	if totalPlants <= 0 {
		totalPlants = int(plantBatch.SeedCount)
	}
	if totalPlants <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "total_plants must be greater than 0"})
		return
	}
	if input.AffectedCount > totalPlants {
		c.JSON(http.StatusBadRequest, gin.H{"error": "affected_count cannot exceed total_plants"})
		return
	}

	diseaseNote := strings.TrimSpace(input.DiseaseNote)
	if disease != "lainnya" {
		diseaseNote = ""
	}

	photoURL := ""
	if isMultipart {
		uploadedPhotoURL, err := savePlantMonitoringImage(c, "photo")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Photo upload failed"})
			return
		}
		photoURL = uploadedPhotoURL
	}
	if photoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Photo is required"})
		return
	}

	entity := maintenancemodels.PlantMonitoring{
		PlantBatchID:  plantBatch.ID,
		Phase:         phaseKey,
		Progress:      progress,
		HealthStatus:  healthStatus,
		Disease:       disease,
		DiseaseNote:   diseaseNote,
		AffectedCount: input.AffectedCount,
		TotalPlants:   totalPlants,
		Note:          strings.TrimSpace(input.Note),
		CreatedBy:     &userID,
		PhotoURL:      photoURL,
	}

	if err := config.DB.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save plant monitoring"})
		return
	}

	createdBy := ""
	if entity.CreatedBy != nil {
		createdBy = *entity.CreatedBy
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Plant monitoring created", "data": plantMonitoringResponse{
		ID:            entity.ID,
		PlantBatchID:  entity.PlantBatchID,
		BatchCode:     plantBatch.BatchCode,
		PackageName:   plantBatch.Investment.Package.PackageName,
		Phase:         entity.Phase,
		Progress:      entity.Progress,
		HealthStatus:  entity.HealthStatus,
		Disease:       entity.Disease,
		DiseaseNote:   entity.DiseaseNote,
		AffectedCount: entity.AffectedCount,
		TotalPlants:   entity.TotalPlants,
		Note:          entity.Note,
		CreatedBy:     createdBy,
		PhotoURL:      entity.PhotoURL,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}})
}

func ListPlantMonitorings(c *gin.Context) {
	var monitorings []maintenancemodels.PlantMonitoring
	query := config.DB.
		Preload("PlantBatch").
		Preload("PlantBatch.Investment").
		Preload("PlantBatch.Investment.Package").
		Order("created_at DESC")

	if batchIDQuery := strings.TrimSpace(c.Query("plant_batch_id")); batchIDQuery != "" {
		batchID, err := strconv.ParseUint(batchIDQuery, 10, 64)
		if err != nil || batchID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant_batch_id query"})
			return
		}
		query = query.Where("plant_batch_id = ?", uint(batchID))
	}

	if userIDValue, exists := c.Get("user_id"); exists {
		userID, _ := userIDValue.(string)
		if roleValue, ok := c.Get("role"); ok {
			role, _ := roleValue.(string)
			if role == "investor" || role == "mitra" {
				query = query.Joins("JOIN plant_batches ON plant_batches.id = plant_monitorings.plant_batch_id").
					Joins("JOIN investments ON investments.id = plant_batches.investment_id").
					Where("investments.user_id = ?", userID)
			}
		}
	}

	if err := query.Find(&monitorings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plant monitoring data"})
		return
	}

	response := make([]plantMonitoringResponse, 0, len(monitorings))
	for _, item := range monitorings {
		createdBy := ""
		if item.CreatedBy != nil {
			createdBy = *item.CreatedBy
		}

		response = append(response, plantMonitoringResponse{
			ID:            item.ID,
			PlantBatchID:  item.PlantBatchID,
			BatchCode:     item.PlantBatch.BatchCode,
			PackageName:   item.PlantBatch.Investment.Package.PackageName,
			Phase:         item.Phase,
			Progress:      item.Progress,
			HealthStatus:  item.HealthStatus,
			Disease:       item.Disease,
			DiseaseNote:   item.DiseaseNote,
			AffectedCount: item.AffectedCount,
			TotalPlants:   item.TotalPlants,
			Note:          item.Note,
			CreatedBy:     createdBy,
			PhotoURL:      item.PhotoURL,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
