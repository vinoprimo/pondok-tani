package models

type Harvest struct {
	ID        uint   `gorm:"primaryKey"`
	Tanggal   string
	BatchID   uint

	Outputs   []HarvestOutput `gorm:"foreignKey:HarvestID"`
}