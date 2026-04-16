package harvest

import "errors"

var (
	ErrInvalidHarvestDate = errors.New("harvest_date is required")
	ErrInvalidBatchID     = errors.New("plant_batch_id must be greater than 0")
	ErrHarvestNotFound    = errors.New("harvest not found")
	ErrInvalidJenis       = errors.New("jenis must be either basah or kering")
	ErrInvalidBerat       = errors.New("quantity must be greater than 0")
)
