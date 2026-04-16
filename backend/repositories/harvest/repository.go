package harvest

import (
	"errors"

	harvestmodels "pondok-tani-backend/models/harvest"
	warehousemodels "pondok-tani-backend/models/warehouse"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HarvestRepository interface {
	Create(tx *gorm.DB, harvest *harvestmodels.Harvest) error
	Exists(tx *gorm.DB, id uint) (bool, error)
	FindByID(tx *gorm.DB, id uint) (*harvestmodels.Harvest, error)
}

type HarvestOutputRepository interface {
	Create(tx *gorm.DB, output *harvestmodels.HarvestOutput) error
}

type WarehouseStockRepository interface {
	AddByScope(tx *gorm.DB, plantBatchID uint, gradeID *uint, delta float64) (*warehousemodels.WarehouseStock, error)
}

type StockMovementRepository interface {
	Create(tx *gorm.DB, movement *warehousemodels.StockMovement) error
}

type GormHarvestRepository struct {
	db *gorm.DB
}

type GormHarvestOutputRepository struct {
	db *gorm.DB
}

type GormWarehouseStockRepository struct {
	db *gorm.DB
}

type GormStockMovementRepository struct {
	db *gorm.DB
}

func NewGormHarvestRepository(db *gorm.DB) *GormHarvestRepository {
	return &GormHarvestRepository{db: db}
}

func NewGormHarvestOutputRepository(db *gorm.DB) *GormHarvestOutputRepository {
	return &GormHarvestOutputRepository{db: db}
}

func NewGormWarehouseStockRepository(db *gorm.DB) *GormWarehouseStockRepository {
	return &GormWarehouseStockRepository{db: db}
}

func NewGormStockMovementRepository(db *gorm.DB) *GormStockMovementRepository {
	return &GormStockMovementRepository{db: db}
}

func dbOrTx(defaultDB *gorm.DB, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return defaultDB
}

func (r *GormHarvestRepository) Create(tx *gorm.DB, harvest *harvestmodels.Harvest) error {
	return dbOrTx(r.db, tx).Create(harvest).Error
}

func (r *GormHarvestRepository) Exists(tx *gorm.DB, id uint) (bool, error) {
	var count int64
	err := dbOrTx(r.db, tx).Model(&harvestmodels.Harvest{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormHarvestRepository) FindByID(tx *gorm.DB, id uint) (*harvestmodels.Harvest, error) {
	var entity harvestmodels.Harvest
	err := dbOrTx(r.db, tx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *GormHarvestOutputRepository) Create(tx *gorm.DB, output *harvestmodels.HarvestOutput) error {
	return dbOrTx(r.db, tx).Create(output).Error
}

func (r *GormStockMovementRepository) Create(tx *gorm.DB, movement *warehousemodels.StockMovement) error {
	return dbOrTx(r.db, tx).Create(movement).Error
}

func (r *GormWarehouseStockRepository) AddByScope(tx *gorm.DB, plantBatchID uint, gradeID *uint, delta float64) (*warehousemodels.WarehouseStock, error) {
	db := dbOrTx(r.db, tx)

	var stock warehousemodels.WarehouseStock

	query := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("plant_batch_id = ?", plantBatchID)
	if gradeID == nil {
		query = query.Where("grade_id IS NULL")
	} else {
		query = query.Where("grade_id = ?", *gradeID)
	}

	err := query.First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			stock = warehousemodels.WarehouseStock{
				PlantBatchID:  plantBatchID,
				GradeID:       gradeID,
				TotalQuantity: delta,
				Unit:          "kg",
			}
			if createErr := db.Create(&stock).Error; createErr != nil {
				return nil, createErr
			}
			return &stock, nil
		}
		return nil, err
	}

	stock.TotalQuantity += delta
	if err := db.Save(&stock).Error; err != nil {
		return nil, err
	}

	return &stock, nil
}
