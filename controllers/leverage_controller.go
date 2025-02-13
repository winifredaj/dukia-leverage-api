package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dukia-leverage-api/config"
	"dukia-leverage-api/models"
)

// GetLeverage(Placeholder function for now)
func GetLeverage(c *gin.Context) {
	
	c.JSON(http.StatusOK, gin.H{
        "message": "GetLeverage Controller is working",
    })
}

// ApplyLeverage(Placeholder function for now)
 func ApplyLeverage(c *gin.Context) {
    var request struct {
        UserID              uint    `json:"user_id"`
        GoldHoldingID       uint    `json:"gold_holding_id"`
        LeverageAmount      float64 `json:"leverage_amount"`
        TenureMonths        int     `json:"tenure_months"`     
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    //Ensure gold holding exists
    var goldHolding models.GoldHolding
    result := config.DB.Where("id = ?", request.GoldHoldingID).First(&goldHolding)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Gold Holding ID not found"})
        return
    }
    
    //Check if leverage limit(75% of gold holding) has been reached
    maxLeverage := goldHolding.CurrentValue * 0.75
    if request.LeverageAmount > maxLeverage {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Leverage amount exceeds allowed limit "})
        return
    }
    

    //Save leverage requests
    leverage := models.LeverageTransaction{
        UserID:              request.UserID,
        GoldHoldingID:       request.GoldHoldingID,
        LeverageAmount:      request.LeverageAmount,
        TenureMonths:        request.TenureMonths,
        Status:               "pending",
    }
    config.DB.Create(&leverage)
    

    c.JSON(http.StatusOK, gin.H{
        "message": "Leverage request submitted","leverage_id": leverage.ID})
}