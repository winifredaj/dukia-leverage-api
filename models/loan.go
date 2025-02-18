package models

import "gorm.io/gorm"

//Loan represents a loan in the system
type Loan struct {
	gorm.Model      //Use gorm.Model for auto-generated fields
	ID         		int			`gorm:"primaryKey"`
	UserID      	int 		`gorm:"not null"`
	Amount      	float64 	`gorm:"not null"`
	CollateralGold 	float64     `gorm:"not null"`
	Status          string      `gorm:"type:loan_status; default:'inactive'"`

}

//GetActiveLoans retrieves all active loans from datatbase

func GetActiveLoans() ([]Loan, error) {
	//Plceholder impleentation
	//Replace with actual database retrival logic
    return []Loan{
		{ID: 1, UserID: 1, Amount: 1000, CollateralGold: 500},
        {ID: 2, UserID: 2, Amount: 2000, CollateralGold: 1000},

	}, nil
	}

  //LiquidateCollateral processes the liquidation of a loan's collateral
  func LiquidateCollateral(loanID int) error{
	// Placeholder implementation
	// Replace with actual liquidation logic
	//Update the loan status in the database
    return nil
    }

  