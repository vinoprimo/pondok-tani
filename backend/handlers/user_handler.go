package handlers

import (
	"context"
	"net/http"

	"pondok-tani-backend/config"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := config.DB.QueryRow(context.Background(),
		"SELECT id, name, email FROM users WHERE id=$1",
		userID,
	).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}