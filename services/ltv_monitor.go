package services

import (
	"dukia-leverage-api/config"
	"dukia-leverage-api/models"
	"dukia-leverage-api/utils"
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

// MonitorLTVandLiquidate checks for loans that exceed the LTV threshold and liquidates them
func MonitorLTVandLiquidate() {
	c := cron.New()

	// Check  LTV every 5 minutes
	_, err := c.AddFunc("*/5 * * * *", func() {
		fmt.Println("Running LTV check...")

		// Fetch all active leverage transactions with gold holdings
		var transacations []models.LeverageTransaction
		config.DB.Find(&transacations)

		for _, loan := range transacations {
			//Calculate the current LTV using the amount, required collateral, and gold price
			var goldHolding models.GoldHolding

			if err := config.DB.Where("id= ?", loan.GoldHoldingID).First(&goldHolding).Error; err != nil {
				log.Println("Error fetching gold holding for loan ID:", loan.ID, err)
				continue
			}
			goldPrice, err := utils.GetCurrentGoldPrice()
			if err != nil {
				log.Println("Error fetching gold price", err)
				continue
			}

			//Update gold value and recalculate LTV
			totalGoldValue := goldHolding.Quantity * goldPrice
			currentLTV := (loan.LeverageAmount / totalGoldValue) * 100

			//Trigger margin call at 85%
			if currentLTV > 85 {
				handleMarginCall(loan, int(loan.UserID))

			}

			//Liqidate at 90%
			if currentLTV > 90 {
				handleLiquidation(loan, int(loan.UserID))
			}
		}

	})

	if err != nil {
		log.Fatalf("Failed to start LTV monitor:%v ", err)
	}

	c.Start()
}

func handleMarginCall(loan models.LeverageTransaction, userID int) {
	//Implementation of margin call
	err := utils.SendNotification(userID, "Your leverage transaction has reached the 85% LTV threshold. Margin call initiated.")
	if err != nil {
		log.Println("Error sending notification for transaction ID:", loan.ID, err)
	}
}

func handleLiquidation(loan models.LeverageTransaction, userID int) {

	err := utils.SendNotification(userID, "Your leverage transaction has reached the 90% LTV threshold. Your collateral has been liquidated.")
	if err != nil {
		log.Println("Error sending notification for transaction ID:", loan.ID, err)
	}

	//Update loan status to liquidated
	loan.Status = "Liquidated"
	config.DB.Save(&loan)
}
