package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	notificationmodels "pondok-tani-backend/models/notification"

	"github.com/gin-gonic/gin"
)

type pushNotificationRequest struct {
	UserID string            `json:"user_id" binding:"required,uuid"`
	Title  string            `json:"title" binding:"required"`
	Body   string            `json:"body" binding:"required"`
	Data   map[string]string `json:"data,omitempty"`
}

type fcmNotificationPayload struct {
	To           string            `json:"to"`
	Notification *notificationBody `json:"notification,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Priority     string            `json:"priority,omitempty"`
}

type notificationBody struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func sendFCMNotification(fcmToken, title, body string, data map[string]string) error {
	serverKey := os.Getenv("FCM_SERVER_KEY")
	if serverKey == "" {
		return errors.New("FCM_SERVER_KEY belum dikonfigurasi")
	}

	payload := fcmNotificationPayload{
		To: fcmToken,
		Notification: &notificationBody{
			Title: title,
			Body:  body,
		},
		Data:     data,
		Priority: "high",
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+serverKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("FCM request failed")
	}

	return nil
}

func SendPushNotification(c *gin.Context) {
	var req pushNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user authmodels.User
	if err := config.DB.First(&user, "id = ?", req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.FCMToken == nil || *user.FCMToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User tidak memiliki FCM token"})
		return
	}

	if err := sendFCMNotification(*user.FCMToken, req.Title, req.Body, req.Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Push notification sent"})
}

func GetMyNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var notifications []notificationmodels.Notification
	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	// Format response to match frontend expectations (title and message)
	response := make([]map[string]interface{}, 0)
	for _, n := range notifications {
		response = append(response, map[string]interface{}{
			"id":         n.ID,
			"title":      n.Subject,
			"message":    n.Message,
			"type":       n.Type,
			"read":       n.IsRead,
			"created_at": n.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func MarkNotificationRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notificationID := c.Param("id")

	if err := config.DB.Model(&notificationmodels.Notification{}).Where("id = ? AND user_id = ?", notificationID, userID).Update("is_read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}
