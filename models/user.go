package models

import "gorm.io/gorm"

type User struct {
	gorm.Model 		 //Use gorm.Model for auto-generated fields
	
	ID               uint    `gorm:"primaryKey"`
	Name             string  `gorm:"size:255"`
	Email            string  `gorm:"unique"`
	TransactionCode  string  `gorm:"size:6"`
	Password         string  `gorm:"size:255"`

}