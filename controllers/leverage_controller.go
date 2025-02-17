package controllers

import (
	"fmt"
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

// Apply for Leverage
 func ApplyLeverage(c *gin.Context) {
    var request struct {
        UserID              uint    `json:"UserID"`
        GoldHoldingID       uint    `json:"GoldHoldingID"`
        LeverageAmount      float64 `json:"LeverageAmount"`
        TenureMonths        int     `json:"TenureMonths"`     
    }

    //Debugging: Print the incoming JSON request
    if err := c.ShouldBindJSON(&request); err != nil {
        fmt.Println("JSON Binding Error:",err) 
        c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid request format"})
        return
    }
    fmt.Printf("Received Request: %+v\n", request) // ✅ Debugging print


    //Validate GoldHoldingID exists before querying
    //To-Do: Implement a real-world validation for gold holding existence
    if request.GoldHoldingID == 0 {
        fmt.Println("Invalid GoldHoldingID:", request.GoldHoldingID)
        c.JSON(http.StatusBadRequest, gin.H{"error": "GoldHoldingID cannot be 0 "})
        return  
    }
    
    //Check if user has sufficient balance
    //To-Do: Implement a real-world validation for user balance
    //if user.Balance < request.LeverageAmount {
    //    c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Insufficient balance"})
    //    return  
    //}


    //Ensure gold holding exists
    var goldHolding models.GoldHolding
    result := config.DB.Where("id = ?", request.GoldHoldingID).First(&goldHolding)

    if result.Error != nil {
        fmt.Println("Gold Holding Not Found:", result.Error)  // Debugging print
        c.JSON(http.StatusNotFound, gin.H{"error": "Gold Holding ID not found"})
        return
    }
    //Print the retrieved GoldHolding
    fmt.Printf("Found Gold Holding: %+v\n", goldHolding)
    
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