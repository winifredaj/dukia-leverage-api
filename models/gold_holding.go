package models

import "gorm.io/gorm"

type GoldHolding struct {
	gorm.Model
	ID           uint    `gorm:"primaryKey"`
	UserID       uint    `gorm:"not null"`
	Quantity     float64 `gorm:"not null"`
	CurrentValue float64 `gorm:"not null"`
	GoldType     string  `gorm:"not null"`
}
