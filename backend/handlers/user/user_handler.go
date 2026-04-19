package user

import (
	"net/http"
	"strings"
	"time"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type adminCreateUserRequest struct {
	Name        string  `json:"name" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	Role        string  `json:"role"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
}

type adminUpdateUserRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	Role        string  `json:"role"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
}

func normalizeRole(role string) string {
	normalized := strings.ToLower(strings.TrimSpace(role))
	if normalized == "" {
		return "investor"
	}
	if normalized != "investor" && normalized != "mitra" && normalized != "admin" {
		return "investor"
	}
	return normalized
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user authmodels.User

	err := config.DB.
		Preload("SelectedPackage").
		Preload("Investments").
		Preload("Investments.Package").
		First(&user, "id = ?", userID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func ListUsers(c *gin.Context) {
	var users []authmodels.User
	if err := config.DB.
		Preload("SelectedPackage").
		Preload("Investments").
		Preload("Investments.Package").
		Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var req adminCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare password"})
		return
	}

	now := time.Now()
	user := authmodels.User{
		ID:          uuid.New().String(),
		Name:        strings.TrimSpace(req.Name),
		Email:       strings.TrimSpace(strings.ToLower(req.Email)),
		Password:    string(hashedPassword),
		Role:        normalizeRole(req.Role),
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user. Email might already be used"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": user.ID})
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req adminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user authmodels.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	updates := map[string]interface{}{}

	if strings.TrimSpace(req.Name) != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}
	if strings.TrimSpace(req.Email) != "" {
		updates["email"] = strings.TrimSpace(strings.ToLower(req.Email))
	}
	if strings.TrimSpace(req.Role) != "" {
		updates["role"] = normalizeRole(req.Role)
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = req.PhoneNumber
	}
	if req.Address != nil {
		updates["address"] = req.Address
	}
	if strings.TrimSpace(req.Password) != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare password"})
			return
		}
		updates["password"] = string(hashedPassword)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No changes provided"})
		return
	}

	if err := config.DB.Model(&authmodels.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	requestUserID, _ := c.Get("user_id")

	if requestUserID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start deletion transaction"})
		return
	}

	var investmentIDs []uint
	if err := tx.Model(&coremodels.Investment{}).
		Where("user_id = ?", userID).
		Pluck("id", &investmentIDs).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load user investments"})
		return
	}

	if len(investmentIDs) > 0 {
		if err := tx.Where("investment_id IN ?", investmentIDs).Delete(&coremodels.PlantBatch{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user plant batches"})
			return
		}

		if err := tx.Where("id IN ?", investmentIDs).Delete(&coremodels.Investment{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user investments"})
			return
		}
	}

	result := tx.Delete(&authmodels.User{}, "id = ?", userID)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize user deletion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
