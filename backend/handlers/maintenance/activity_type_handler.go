package maintenance

import (
	"math"
	"net/http"
	"pondok-tani-backend/models/maintenance"
	repo "pondok-tani-backend/repositories/maintenance"
	"pondok-tani-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllActivityTypes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	activityTypes, total, err := repo.GetAllActivityTypes(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity types"})
		return
	}

	pagination := &utils.PaginationData{
		Rows:       activityTypes,
		TotalRows:  total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}

	response := utils.FormatPaginationResponse(pagination)
	c.JSON(http.StatusOK, response)
}

func GetActivityTypeByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	activityType, err := repo.GetActivityTypeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activityType)
}

func CreateActivityType(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activityType := maintenance.ActivityType{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := repo.CreateActivityType(&activityType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity type, may be duplicate"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Activity type created successfully", "data": activityType})
}

func UpdateActivityType(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	updates["description"] = input.Description

	if err := repo.UpdateActivityType(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity type updated successfully"})
}

func DeleteActivityType(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := repo.DeleteActivityType(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity type deleted successfully"})
}
