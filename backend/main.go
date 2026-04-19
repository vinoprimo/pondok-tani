package main

import (
	"fmt"
	"os"

	"pondok-tani-backend/config"
	authmodels "pondok-tani-backend/models/auth"
	coremodels "pondok-tani-backend/models/core"
	harvestmodels "pondok-tani-backend/models/harvest"
	maintenancemodels "pondok-tani-backend/models/maintenance"
	notificationmodels "pondok-tani-backend/models/notification"
	postharvestmodels "pondok-tani-backend/models/postharvest"
	pricemodels "pondok-tani-backend/models/price"
	salesmodels "pondok-tani-backend/models/sales"
	warehousemodels "pondok-tani-backend/models/warehouse"
	"pondok-tani-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func migrateDB() {
	config.ConnectDB()
	err := config.DB.AutoMigrate(
		&authmodels.User{},
		&coremodels.InvestmentPackage{},
		&coremodels.Investment{},
		&coremodels.PlantBatch{},
		&maintenancemodels.MaintenanceSchedule{},
		&maintenancemodels.MaintenanceActivity{},
		&harvestmodels.Harvest{},
		&harvestmodels.HarvestOutput{},
		&postharvestmodels.DryingProcess{},
		&postharvestmodels.GradingBatch{},
		&postharvestmodels.Grade{},
		&postharvestmodels.GradingDetail{},
		&warehousemodels.WarehouseStock{},
		&warehousemodels.StockMovement{},
		&salesmodels.SalesOrder{},
		&salesmodels.SalesDetail{},
		&notificationmodels.Notification{},
		&pricemodels.NationalPrice{},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	if err := config.DB.Exec(`ALTER TABLE plant_batches ADD COLUMN IF NOT EXISTS location varchar(150) NOT NULL DEFAULT ''`).Error; err != nil {
		panic(fmt.Sprintf("failed to ensure plant_batches.location column: %v", err))
	}
	if err := config.DB.Exec(`ALTER TABLE plant_batches ADD COLUMN IF NOT EXISTS seed_count bigint NOT NULL DEFAULT 0`).Error; err != nil {
		panic(fmt.Sprintf("failed to ensure plant_batches.seed_count column: %v", err))
	}

	if config.DB.Migrator().HasColumn(&coremodels.Investment{}, "expected_return_at") {
		if err := config.DB.Migrator().DropColumn(&coremodels.Investment{}, "expected_return_at"); err != nil {
			panic(fmt.Sprintf("failed to drop expected_return_at column: %v", err))
		}
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
	err := config.DB.AutoMigrate(
		&authmodels.User{},
		&coremodels.InvestmentPackage{},
		&coremodels.Investment{},
		&coremodels.PlantBatch{},
		&maintenancemodels.MaintenanceSchedule{},
		&maintenancemodels.MaintenanceActivity{},
		&harvestmodels.Harvest{},
		&harvestmodels.HarvestOutput{},
		&postharvestmodels.DryingProcess{},
		&postharvestmodels.GradingBatch{},
		&postharvestmodels.Grade{},
		&postharvestmodels.GradingDetail{},
		&warehousemodels.WarehouseStock{},
		&warehousemodels.StockMovement{},
		&salesmodels.SalesOrder{},
		&salesmodels.SalesDetail{},
		&notificationmodels.Notification{},
		&pricemodels.NationalPrice{},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	if err := config.DB.Exec(`ALTER TABLE plant_batches ADD COLUMN IF NOT EXISTS location varchar(150) NOT NULL DEFAULT ''`).Error; err != nil {
		panic(fmt.Sprintf("failed to ensure plant_batches.location column: %v", err))
	}
	if err := config.DB.Exec(`ALTER TABLE plant_batches ADD COLUMN IF NOT EXISTS seed_count bigint NOT NULL DEFAULT 0`).Error; err != nil {
		panic(fmt.Sprintf("failed to ensure plant_batches.seed_count column: %v", err))
	}

	if config.DB.Migrator().HasColumn(&coremodels.Investment{}, "expected_return_at") {
		if err := config.DB.Migrator().DropColumn(&coremodels.Investment{}, "expected_return_at"); err != nil {
			panic(fmt.Sprintf("failed to drop expected_return_at column: %v", err))
		}
	}

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
	routes.UserRoutes(r)
	routes.AdminRoutes(r)
	routes.HarvestRoutes(r)
	routes.InvestmentPackageRoutes(r)
	r.Run(":8000")
}
