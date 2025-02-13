package models

import "gorm.io/gorm"

type GoldHolding struct {
	gorm.Model
	UserID      	uint 		`gorm:"not null"`
	Weight				uint 		`gorm:"not null"`
	CurrentValue 	float64 `gorm:"not null"`
	Status  			string 	`gorm:"not null"`
}