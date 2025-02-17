package models

//import "errors"

//Loan represents a loan in the system
type Loan struct {
	ID         		int			`json:"id"`
	UserID      	int 		`json:"user_id"`
	Amount      	float64 	`json:"amount"`
	CollateralGold 	float64

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

  