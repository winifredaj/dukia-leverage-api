package models

import "gorm.io/gorm"

type MarginCall struct {
	gorm.Model

	ID                    uint `gorm:"primaryKey"`
	LeverageTransactionID uint `gorm:"not null"`
	UserID                uint
	RequiredCollateral    float64 `gorm:"not null"`
	Status                string  `gorm:"type:margincall_status; default:'pending'"`
}
