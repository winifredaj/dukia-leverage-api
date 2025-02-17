package models

import  "gorm.io/gorm"

type LeverageTransaction struct {
	gorm.Model   

	ID              uint    `gorm:"primaryKey"`
	UserID          uint    `gorm:"not null"`
	GoldHoldingID   uint    `gorm:"not null"`
	LeverageAmount  float64 `gorm:"not null"`
	InterestRate    float64 `gorm:"default:28.0"`
	TenureMonths    int     `gorm:"not null"`
	Status          string  `gorm:"type:leverage_status; default:'pending'"`
	
}


