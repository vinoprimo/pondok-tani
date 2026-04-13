package main

import (
	"fmt"
	"os"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	harvestmodels "pondok-tani-backend/models/harvest"
	"pondok-tani-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func migrateDB() {
	config.ConnectDB()
	err := config.DB.AutoMigrate(&authmodels.User{}, &harvestmodels.Harvest{}, &harvestmodels.HarvestOutput{}, &harvestmodels.StockMovement{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}
	fmt.Println("database migrated successfully")
}

func main() {
	godotenv.Load()

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		migrateDB()
		return
	}

	config.ConnectDB()
	config.DB.AutoMigrate(&authmodels.User{}, &harvestmodels.Harvest{}, &harvestmodels.HarvestOutput{}, &harvestmodels.StockMovement{})

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
		})
	})
	routes.AuthRoutes(r)
	r.Run(":8000")
}
