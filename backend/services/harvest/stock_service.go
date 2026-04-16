package harvest

import (
	"context"

	warehousemodels "pondok-tani-backend/models/warehouse"
	harvestrepo "pondok-tani-backend/repositories/harvest"
	"time"

	"gorm.io/gorm"
)

type StockService interface {
	RecordIncomingFromHarvestOutput(ctx context.Context, tx *gorm.DB, harvestOutputID uint, plantBatchID uint, gradeID *uint, qty float64) error
}

type stockService struct {
	stockRepo    harvestrepo.WarehouseStockRepository
	movementRepo harvestrepo.StockMovementRepository
}

func NewStockService(
	stockRepo harvestrepo.WarehouseStockRepository,
	movementRepo harvestrepo.StockMovementRepository,
) StockService {
	return &stockService{
		stockRepo:    stockRepo,
		movementRepo: movementRepo,
	}
}

func (s *stockService) RecordIncomingFromHarvestOutput(
	_ context.Context,
	tx *gorm.DB,
	harvestOutputID uint,
	plantBatchID uint,
	gradeID *uint,
	qty float64,
) error {
	if qty <= 0 {
		return ErrInvalidBerat
	}

	stock, err := s.stockRepo.AddByScope(tx, plantBatchID, gradeID, qty)
	if err != nil {
		return err
	}

	movement := &warehousemodels.StockMovement{
		WarehouseStockID: stock.ID,
		ReferenceType:    "harvest_output",
		ReferenceID:      harvestOutputID,
		MovementType:     "in",
		Quantity:         qty,
		MovementDate:     time.Now(),
	}

	if err := s.movementRepo.Create(tx, movement); err != nil {
		return err
	}

	return nil
}
