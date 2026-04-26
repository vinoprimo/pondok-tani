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

func applyNationalPriceMigrations() error {
	if config.DB.Migrator().HasColumn(&pricemodels.NationalPrice{}, "id") && !config.DB.Migrator().HasColumn(&pricemodels.NationalPrice{}, "national_price_id") {
		if err := config.DB.Exec(`ALTER TABLE national_prices RENAME COLUMN id TO national_price_id`).Error; err != nil {
			return fmt.Errorf("failed to rename national_prices.id column: %w", err)
		}
	}

	if !config.DB.Migrator().HasColumn(&pricemodels.NationalPrice{}, "harvest_type") {
		if err := config.DB.Exec(`ALTER TABLE national_prices ADD COLUMN IF NOT EXISTS harvest_type varchar(20) NOT NULL DEFAULT 'basah'`).Error; err != nil {
			return fmt.Errorf("failed to ensure national_prices.harvest_type column: %w", err)
		}
	}

	if config.DB.Migrator().HasColumn(&pricemodels.NationalPrice{}, "source") {
		if err := config.DB.Migrator().DropColumn(&pricemodels.NationalPrice{}, "source"); err != nil {
			return fmt.Errorf("failed to drop national_prices.source column: %w", err)
		}
	}

	if err := config.DB.Exec(`UPDATE national_prices SET harvest_type = lower(trim(harvest_type)) WHERE harvest_type IS NOT NULL`).Error; err != nil {
		return fmt.Errorf("failed to normalize national_prices.harvest_type values: %w", err)
	}
	if err := config.DB.Exec(`UPDATE national_prices SET harvest_type = 'basah' WHERE harvest_type IS NULL OR harvest_type NOT IN ('basah','kering')`).Error; err != nil {
		return fmt.Errorf("failed to repair invalid national_prices.harvest_type values: %w", err)
	}
	if err := config.DB.Exec(`UPDATE national_prices SET grade_id = NULL WHERE harvest_type = 'basah'`).Error; err != nil {
		return fmt.Errorf("failed to clear grade_id for basah harvest_type rows: %w", err)
	}

	if err := config.DB.Exec(`ALTER TABLE national_prices DROP CONSTRAINT IF EXISTS chk_national_prices_harvest_type`).Error; err != nil {
		return fmt.Errorf("failed to drop old national_prices harvest_type constraint: %w", err)
	}
	if err := config.DB.Exec(`ALTER TABLE national_prices ADD CONSTRAINT chk_national_prices_harvest_type CHECK (harvest_type IN ('basah','kering'))`).Error; err != nil {
		return fmt.Errorf("failed to add national_prices harvest_type constraint: %w", err)
	}

	if err := config.DB.Exec(`ALTER TABLE national_prices DROP CONSTRAINT IF EXISTS chk_national_prices_basah_no_grade`).Error; err != nil {
		return fmt.Errorf("failed to drop old national_prices basah-grade constraint: %w", err)
	}
	if err := config.DB.Exec(`ALTER TABLE national_prices ADD CONSTRAINT chk_national_prices_basah_no_grade CHECK (harvest_type <> 'basah' OR grade_id IS NULL)`).Error; err != nil {
		return fmt.Errorf("failed to add national_prices basah-grade constraint: %w", err)
	}

	return nil
}

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

	if err := applyNationalPriceMigrations(); err != nil {
		panic(err.Error())
	}

	if err := config.DB.Exec(`ALTER TABLE plant_batches ADD COLUMN IF NOT EXISTS location varchar(150) NOT NULL DEFAULT ''`).Error; err != nil {
		panic(fmt.Sprintf("failed to ensure plant_batches.location column: %v", err))
	}
	if err := config.DB.Exec(`ALTER TABLE plant_batches ADD COLUMN IF NOT EXISTS seed_count bigint NOT NULL DEFAULT 0`).Error; err != nil {
		panic(fmt.Sprintf("failed to ensure plant_batches.seed_count column: %v", err))
	}

	if err := config.DB.Exec(`ALTER TABLE maintenance_schedules ADD COLUMN IF NOT EXISTS user_id uuid`).Error; err != nil {
		panic(fmt.Sprintf("failed to ensure maintenance_schedules.user_id column: %v", err))
	}
	if err := config.DB.Exec(`
		UPDATE maintenance_schedules AS ms
		SET user_id = inv.user_id
		FROM plant_batches AS pb
		JOIN investments AS inv ON inv.id = pb.investment_id
		WHERE ms.plant_batch_id = pb.id
		  AND ms.user_id IS NULL
	`).Error; err != nil {
		panic(fmt.Sprintf("failed to backfill maintenance_schedules.user_id values: %v", err))
	}
	if err := config.DB.Exec(`ALTER TABLE maintenance_schedules DROP CONSTRAINT IF EXISTS fk_maintenance_schedules_user`).Error; err != nil {
		panic(fmt.Sprintf("failed to reset maintenance_schedules.user foreign key: %v", err))
	}
	if err := config.DB.Exec(`
		ALTER TABLE maintenance_schedules
		ADD CONSTRAINT fk_maintenance_schedules_user
		FOREIGN KEY (user_id) REFERENCES users(id)
		ON UPDATE CASCADE ON DELETE RESTRICT
	`).Error; err != nil {
		panic(fmt.Sprintf("failed to apply maintenance_schedules.user foreign key: %v", err))
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

	if err := applyNationalPriceMigrations(); err != nil {
		panic(err.Error())
	}

	if err := config.DB.Exec(`ALTER TABLE plant_batches ADD COLUMN IF NOT EXISTS location varchar(150) NOT NULL DEFAULT ''`).Error; err != nil {
		panic(fmt.Sprintf("failed to ensure plant_batches.location column: %v", err))
	}
	if err := config.DB.Exec(`ALTER TABLE plant_batches ADD COLUMN IF NOT EXISTS seed_count bigint NOT NULL DEFAULT 0`).Error; err != nil {
		panic(fmt.Sprintf("failed to ensure plant_batches.seed_count column: %v", err))
	}

	if err := config.DB.Exec(`ALTER TABLE maintenance_schedules ADD COLUMN IF NOT EXISTS user_id uuid`).Error; err != nil {
		panic(fmt.Sprintf("failed to ensure maintenance_schedules.user_id column: %v", err))
	}
	if err := config.DB.Exec(`
		UPDATE maintenance_schedules AS ms
		SET user_id = inv.user_id
		FROM plant_batches AS pb
		JOIN investments AS inv ON inv.id = pb.investment_id
		WHERE ms.plant_batch_id = pb.id
		  AND ms.user_id IS NULL
	`).Error; err != nil {
		panic(fmt.Sprintf("failed to backfill maintenance_schedules.user_id values: %v", err))
	}
	if err := config.DB.Exec(`ALTER TABLE maintenance_schedules DROP CONSTRAINT IF EXISTS fk_maintenance_schedules_user`).Error; err != nil {
		panic(fmt.Sprintf("failed to reset maintenance_schedules.user foreign key: %v", err))
	}
	if err := config.DB.Exec(`
		ALTER TABLE maintenance_schedules
		ADD CONSTRAINT fk_maintenance_schedules_user
		FOREIGN KEY (user_id) REFERENCES users(id)
		ON UPDATE CASCADE ON DELETE RESTRICT
	`).Error; err != nil {
		panic(fmt.Sprintf("failed to apply maintenance_schedules.user foreign key: %v", err))
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
