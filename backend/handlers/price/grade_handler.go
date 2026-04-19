package price

import (
	"net/http"
	"strings"

	"pondok-tani-backend/config"
	postharvestmodels "pondok-tani-backend/models/postharvest"

	"github.com/gin-gonic/gin"
)

type createOrUpdateGradeRequest struct {
	GradeName   string  `json:"grade_name" binding:"required"`
	Description *string `json:"description"`
}

func ListGrades(c *gin.Context) {
	var grades []postharvestmodels.Grade
	if err := config.DB.Order("id ASC").Find(&grades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch grades"})
		return
	}

	c.JSON(http.StatusOK, grades)
}

func CreateGrade(c *gin.Context) {
	var req createOrUpdateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	gradeName := strings.TrimSpace(req.GradeName)
	if gradeName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "grade_name is required"})
		return
	}

	grade := postharvestmodels.Grade{
		GradeName:   gradeName,
		Description: req.Description,
	}

	if err := config.DB.Create(&grade).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create grade. Grade name may already exist"})
		return
	}

	c.JSON(http.StatusCreated, grade)
}

func UpdateGrade(c *gin.Context) {
	var grade postharvestmodels.Grade
	if err := config.DB.First(&grade, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grade not found"})
		return
	}

	var req createOrUpdateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	gradeName := strings.TrimSpace(req.GradeName)
	if gradeName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "grade_name is required"})
		return
	}

	grade.GradeName = gradeName
	grade.Description = req.Description

	if err := config.DB.Save(&grade).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update grade"})
		return
	}

	c.JSON(http.StatusOK, grade)
}

func DeleteGrade(c *gin.Context) {
	if err := config.DB.Delete(&postharvestmodels.Grade{}, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete grade"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Grade deleted"})
}
