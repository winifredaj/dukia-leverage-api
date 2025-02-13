package models

import "gorm.io/gorm"

type GoldHolding struct {
	gorm.Model
	UserID      	uint 		`gorm:"not null"`
	Weight		float64 	`gorm:"not null"`
	CurrentValue 	float64 	`gorm:"not null"`
	GoldType 	string	 	`gorm:"not null"`
	Status  	string 		`gorm:"type:enum('active','pledged','liquidated'); default:'active'"`

}
