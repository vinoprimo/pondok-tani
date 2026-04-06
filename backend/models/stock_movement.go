package models

type StockMovement struct {
	ID               uint    `gorm:"primaryKey"`
	WarehouseStockID uint
	ReferenceType    string
	ReferenceID      uint
	MovementType     string
	Quantity         float64
}