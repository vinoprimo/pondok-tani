package models

type HarvestOutput struct {
	ID        uint    `gorm:"primaryKey"`
	HarvestID uint
	Jenis     string
	Berat     float64
	Status    string
}