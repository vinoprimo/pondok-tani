package harvest

import (
	"context"

	harvestmodels "pondok-tani-backend/models/harvest"
	harvestrepo "pondok-tani-backend/repositories/harvest"
	"time"
)

type CreateHarvestInput struct {
	PlantBatchID uint
	HarvestDate  time.Time
	CreatedBy    string
	Notes        *string
}

type HarvestService interface {
	CreateHarvest(ctx context.Context, input CreateHarvestInput) (*harvestmodels.Harvest, error)
}

type harvestService struct {
	harvestRepo harvestrepo.HarvestRepository
}

func NewHarvestService(harvestRepo harvestrepo.HarvestRepository) HarvestService {
	return &harvestService{harvestRepo: harvestRepo}
}

func (s *harvestService) CreateHarvest(_ context.Context, input CreateHarvestInput) (*harvestmodels.Harvest, error) {
	if input.HarvestDate.IsZero() {
		return nil, ErrInvalidHarvestDate
	}
	if input.PlantBatchID == 0 {
		return nil, ErrInvalidBatchID
	}

	entity := &harvestmodels.Harvest{
		PlantBatchID: input.PlantBatchID,
		HarvestDate:  input.HarvestDate,
		CreatedBy:    input.CreatedBy,
		Notes:        input.Notes,
	}

	if err := s.harvestRepo.Create(nil, entity); err != nil {
		return nil, err
	}

	return entity, nil
}
