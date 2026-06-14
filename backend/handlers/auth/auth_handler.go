package auth

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
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type mitraRegisterRequest struct {
	Name          string
	Email         string
	Password      string
	BatchCode     string
	PlantingDate  string
	Location      string
	SeedCount     uint
	LandArea      float64
	HealthStatus  string
	AffectedCount int
	Disease       string
	Summary       string
}

func normalizeRole(role string) string {
	normalized := strings.ToLower(strings.TrimSpace(role))
	if normalized == "" {
		return "investor"
	}
	if normalized != "investor" && normalized != "mitra" {
		return "investor"
	}
	return normalized
}

func normalizeMitraHealthStatus(input string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "sehat":
		return "sehat", nil
	case "sebagian terdampak", "sebagian_terdampak":
		return "sebagian_terdampak", nil
	case "mati":
		return "mati", nil
	default:
		return "", errors.New("invalid health_status")
	}
}

func normalizeMitraDisease(input string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "", "null", "tidak ada", "tidak_ada":
		return "", nil
	case "batang busuk", "batang_busuk":
		return "batang_busuk", nil
	case "lainnya":
		return "lainnya", nil
	default:
		return "", errors.New("invalid disease")
	}
}

func parseMitraRegisterForm(c *gin.Context) (mitraRegisterRequest, error) {
	seedCount, err := strconv.ParseUint(strings.TrimSpace(c.PostForm("seed_count")), 10, 64)
	if err != nil || seedCount == 0 {
		return mitraRegisterRequest{}, errors.New("invalid seed_count")
	}

	landArea, err := strconv.ParseFloat(strings.TrimSpace(c.PostForm("land_area")), 64)
	if err != nil || landArea <= 0 {
		return mitraRegisterRequest{}, errors.New("invalid land_area")
	}

	affectedCount := 0
	if value := strings.TrimSpace(c.PostForm("affected_count")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			return mitraRegisterRequest{}, errors.New("invalid affected_count")
		}
		affectedCount = parsed
	}

	req := mitraRegisterRequest{
		Name:          strings.TrimSpace(c.PostForm("name")),
		Email:         strings.TrimSpace(strings.ToLower(c.PostForm("email"))),
		Password:      c.PostForm("password"),
		BatchCode:     strings.TrimSpace(c.PostForm("batch_code")),
		PlantingDate:  strings.TrimSpace(c.PostForm("planting_date")),
		Location:      strings.TrimSpace(c.PostForm("location")),
		SeedCount:     uint(seedCount),
		LandArea:      landArea,
		HealthStatus:  c.PostForm("health_status"),
		AffectedCount: affectedCount,
		Disease:       c.PostForm("disease"),
		Summary:       strings.TrimSpace(c.PostForm("summary")),
	}

	if req.Name == "" || req.Email == "" || len(req.Password) < 6 || req.BatchCode == "" || req.PlantingDate == "" || req.Location == "" {
		return mitraRegisterRequest{}, errors.New("missing required field")
	}

	return req, nil
}

func saveMitraBatchImage(c *gin.Context) (string, error) {
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		return "", errors.New("unsupported file extension")
	}

	if err := os.MkdirAll("uploads/batches", 0o755); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("mitra_batch_%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join("uploads", "batches", fileName)
	if err := c.SaveUploadedFile(fileHeader, fullPath); err != nil {
		return "", err
	}

	return "/uploads/batches/" + fileName, nil
}

func ensureMitraPackage() (coremodels.InvestmentPackage, error) {
	var packageItem coremodels.InvestmentPackage
	if err := config.DB.Where("package_name = ?", "Kemitraan Mitra").First(&packageItem).Error; err == nil {
		return packageItem, nil
	}

	description := "Paket internal untuk pendaftaran dan validasi awal mitra."
	packageItem = coremodels.InvestmentPackage{
		PackageName: "Kemitraan Mitra",
		Description: &description,
		MinQuantity: 1,
		Price:       0,
		Status:      "inactive",
		Roi:         "",
		Duration:    "",
		Benefits:    []string{},
		IsPopular:   false,
	}

	if err := config.DB.Create(&packageItem).Error; err != nil {
		return coremodels.InvestmentPackage{}, err
	}

	return packageItem, nil
}

func Register(c *gin.Context) {
	var input registerRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	user := authmodels.User{
		ID:       uuid.New().String(),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     normalizeRole(input.Role),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func RegisterMitra(c *gin.Context) {
	req, err := parseMitraRegisterForm(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mitra registration form"})
		return
	}

	plantingDate, err := time.Parse("2006-01-02", req.PlantingDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting_date format. Use YYYY-MM-DD"})
		return
	}

	healthStatus, err := normalizeMitraHealthStatus(req.HealthStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid health_status"})
		return
	}

	disease, err := normalizeMitraDisease(req.Disease)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid disease"})
		return
	}

	if req.AffectedCount > int(req.SeedCount) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "affected_count cannot exceed seed_count"})
		return
	}

	photoURL, err := saveMitraBatchImage(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Photo upload failed"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare password"})
		return
	}

	packageItem, err := ensureMitraPackage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare mitra package"})
		return
	}

	now := time.Now()
	userID := uuid.New().String()
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start registration"})
		return
	}

	user := authmodels.User{
		ID:                userID,
		Name:              req.Name,
		Email:             req.Email,
		Password:          string(hashedPassword),
		Role:              "mitra",
		SelectedPackageID: &packageItem.ID,
		PackageStatus:     "pending_validation",
		PackageSelectedAt: &now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to register. Email might already be used"})
		return
	}

	investment := coremodels.Investment{
		UserID:         userID,
		PackageID:      packageItem.ID,
		Amount:         0,
		Status:         "pending_validation",
		InvestmentDate: now,
	}
	if err := tx.Create(&investment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save mitra registration"})
		return
	}

	plantBatch := coremodels.PlantBatch{
		InvestmentID: investment.ID,
		BatchCode:    req.BatchCode,
		PlantingDate: plantingDate,
		Location:     req.Location,
		SeedCount:    req.SeedCount,
		LandArea:     req.LandArea,
		PhotoURL:     photoURL,
		Status:       "pending_validation",
	}
	if err := tx.Create(&plantBatch).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save plant batch. Ensure batch_code is unique"})
		return
	}

	creatorID := userID
	monitoring := maintenancemodels.PlantMonitoring{
		PlantBatchID:  plantBatch.ID,
		Phase:         "panen",
		Progress:      100,
		HealthStatus:  healthStatus,
		Disease:       disease,
		DiseaseNote:   "",
		AffectedCount: req.AffectedCount,
		TotalPlants:   int(req.SeedCount),
		Note:          req.Summary,
		CreatedBy:     &creatorID,
		PhotoURL:      photoURL,
	}
	if err := tx.Create(&monitoring).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save plant monitoring"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish registration"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Mitra registered and waiting for admin validation"})
}

func Login(c *gin.Context) {
	var input loginRequest
	var user authmodels.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := config.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
