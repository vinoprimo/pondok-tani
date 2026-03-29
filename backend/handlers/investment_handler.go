package handlers

import (
	"context"
	"net/http"
	"time"

	"pondok-tani-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateInvestmentRequest struct {
	PackageID string  `json:"package_id"`
}

func CreateInvestment(c *gin.Context) {
	var req CreateInvestmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	// 🔥 ambil data paket
	var price float64

	err := config.DB.QueryRow(context.Background(),
		"SELECT price FROM investment_packages WHERE id=$1",
		req.PackageID,
	).Scan(&price)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid package"})
		return
	}

	investmentID := uuid.New().String()

	_, err = config.DB.Exec(context.Background(),
		`INSERT INTO investments 
		(id, user_id, investment_date, total_investment_amount, status, package_id)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		investmentID,
		userID,
		time.Now(),
		price, // 🔥 dari package
		"active",
		req.PackageID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Investment created",
	})
}

func GetMyInvestments(c *gin.Context) {
	userID, _ := c.Get("user_id")

	rows, err := config.DB.Query(context.Background(),
		`SELECT 
			i.id, 
			i.investment_date, 
			i.total_investment_amount, 
			i.status,
			p.package_name
		FROM investments i
		JOIN investment_packages p ON i.package_id = p.id
		WHERE i.user_id=$1
		ORDER BY i.created_at DESC`,
		userID,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var investments []map[string]interface{}

	for rows.Next() {
		var id, status, packageName string
		var date string
		var amount float64

		err := rows.Scan(&id, &date, &amount, &status, &packageName)
		if err != nil {
			continue
		}

		investments = append(investments, gin.H{
			"id": id,
			"date": date,
			"amount": amount,
			"status": status,
			"package_name": packageName,
		})
	}

	c.JSON(200, investments)
}

func GetInvestmentDetail(c *gin.Context) {
	id := c.Param("id")

	var result map[string]interface{}

	var idResult, status, packageName string
	var amount float64
	var plantQuantity int

	err := config.DB.QueryRow(context.Background(),
		`SELECT 
			i.id,
			i.total_investment_amount,
			i.status,
			p.package_name,
			p.plant_quantity
		FROM investments i
		JOIN investment_packages p ON i.package_id = p.id
		WHERE i.id=$1`,
		id,
	).Scan(
		&idResult,
		&amount,
		&status,
		&packageName,
		&plantQuantity,
	)

	result = map[string]interface{}{
		"id":             idResult,
		"amount":         amount,
		"status":         status,
		"package_name":   packageName,
		"plant_quantity": plantQuantity,
	}

	if err != nil {
		c.JSON(404, gin.H{"error": "Not found"})
		return
	}

	c.JSON(200, result)
}