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
	RecordIncomingFromGradingDetail(ctx context.Context, tx *gorm.DB, gradingDetailID uint, plantBatchID uint, gradeID uint, qty float64) error
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

func (s *stockService) RecordIncomingFromGradingDetail(
	_ context.Context,
	tx *gorm.DB,
	gradingDetailID uint,
	plantBatchID uint,
	gradeID uint,
	qty float64,
) error {
	if qty <= 0 {
		return ErrInvalidBerat
	}

	grade := gradeID
	stock, err := s.stockRepo.AddByScope(tx, plantBatchID, &grade, qty)
	if err != nil {
		return err
	}

	movement := &warehousemodels.StockMovement{
		WarehouseStockID: stock.ID,
		ReferenceType:    "grading_detail",
		ReferenceID:      gradingDetailID,
		MovementType:     "in",
		Quantity:         qty,
		MovementDate:     time.Now(),
	}

	if err := s.movementRepo.Create(tx, movement); err != nil {
		return err
	}

	return nil
}
