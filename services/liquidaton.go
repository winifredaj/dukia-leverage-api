package services

import (
    "dukia-leverage-api/models"
	"dukia-leverage-api/utils"
    "log"
)

// MonitorLTVandLiquidate checks for loans that exceed the LTV threshold and liquidates them
func MonitorLTVandLiquidate(){
	//Fetch all active loans
	loans, err := models.GetActiveLoans()
    if err !=nil {
        log.Println("Error fetching active loans",err)
		return
    }

	//Fetch current Gold price
	goldPrice, err := utils.GetCurrentGoldPrice()
    if err !=nil {
        log.Println("Error fetching gold price",err)
        return
        }
        
    //Iterate over loans to check their current LTV    
    for _, loan := range loans {
        //Calculate the current LTV using the amount, required collateral, and gold price
        currentLTV := (loan.Amount / (loan.CollateralGold *goldPrice)) * 100

        if currentLTV > 75 {

        //Liquidate collateral
        err :=models.LiquidateCollateral(loan.ID)
        if err !=nil {
            log.Println("Error liquidating collateral for loan ID:", loan.ID, err)
            continue

        }

        //Notify user about liquidation
         err = utils.SendNotification(loan.UserID, "Your loan has reached its LTV limit. Your collateral has been liquidated.")
        if err !=nil {
            log.Println("Error sending notification for loan ID:", loan.ID, err)
        }
    }

}

}