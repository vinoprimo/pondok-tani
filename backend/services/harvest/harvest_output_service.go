package harvest

import (
	"context"
	"errors"
	"strings"

	harvestmodels "pondok-tani-backend/models/harvest"
	harvestrepo "pondok-tani-backend/repositories/harvest"

	"gorm.io/gorm"
)

type CreateHarvestOutputInput struct {
	Jenis    string
	Quantity float64
}

type HarvestOutputService interface {
	CreateHarvestOutput(ctx context.Context, harvestID uint, input CreateHarvestOutputInput) (*harvestmodels.HarvestOutput, error)
}

type harvestOutputService struct {
	db          *gorm.DB
	harvestRepo harvestrepo.HarvestRepository
	outputRepo  harvestrepo.HarvestOutputRepository
	stockSvc    StockService
}

func NewHarvestOutputService(
	db *gorm.DB,
	harvestRepo harvestrepo.HarvestRepository,
	outputRepo harvestrepo.HarvestOutputRepository,
	stockSvc StockService,
) HarvestOutputService {
	return &harvestOutputService{
		db:          db,
		harvestRepo: harvestRepo,
		outputRepo:  outputRepo,
		stockSvc:    stockSvc,
	}
}

func (s *harvestOutputService) CreateHarvestOutput(
	ctx context.Context,
	harvestID uint,
	input CreateHarvestOutputInput,
) (*harvestmodels.HarvestOutput, error) {
	if input.Quantity <= 0 {
		return nil, ErrInvalidBerat
	}

	jenis := strings.ToLower(strings.TrimSpace(input.Jenis))
	if jenis != "basah" && jenis != "kering" {
		return nil, ErrInvalidJenis
	}

	harvestEntity, err := s.harvestRepo.FindByID(nil, harvestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHarvestNotFound
		}
		return nil, err
	}

	status := "pending_drying"
	if jenis == "basah" {
		status = "stocked"
	}

	var createdOutput harvestmodels.HarvestOutput

	err = s.db.Transaction(func(tx *gorm.DB) error {
		output := &harvestmodels.HarvestOutput{
			HarvestID: harvestID,
			Jenis:     jenis,
			Quantity:  input.Quantity,
			Status:    status,
		}

		if err := s.outputRepo.Create(tx, output); err != nil {
			return err
		}

		if jenis == "basah" {
			if err := s.stockSvc.RecordIncomingFromHarvestOutput(ctx, tx, output.ID, harvestEntity.PlantBatchID, nil, input.Quantity); err != nil {
				return err
			}
		}

		createdOutput = *output
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &createdOutput, nil
}
