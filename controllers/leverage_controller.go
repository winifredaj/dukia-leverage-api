package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"dukia-leverage-api/config"
	"dukia-leverage-api/models"
	"dukia-leverage-api/utils"
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

  
    //Check if gold holding exists
    var goldHolding models.GoldHolding
    if err := config.DB.Where("id = ?", request.GoldHoldingID).First(&goldHolding).Error; err != nil {
        c.JSON(http.StatusNotFound,gin.H{"error":"Gold Holding ID not found"})
        return
    }

    //Ensure user has atleast 50g of gold
    if goldHolding.Quantity < 50 {
        c.JSON(http.StatusBadRequest,gin.H{"error":"Insufficient gold balance. Minimum 50g required."})
        return
    }
    
    //Fetch current gold price
    goldPrice, err := utils.GetCurrentGoldPrice()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to retrieve market price"})
        return
    }

    //Calculate Total value and max leverage
    totalGoldValue := goldHolding.Quantity * goldPrice
    maxLeverage := totalGoldValue * 0.75 // 75% of total gold value
    
    if request.LeverageAmount > maxLeverage {
        c.JSON(http.StatusBadRequest,gin.H{"error":"Requested Leverage amount exceeds allowed limit."})
        return
    }

    //Calculate net amount and fees
    processingFee := request.LeverageAmount * 0.01 // 1% processing fee
    custodianFee := request.LeverageAmount * 0.025// 2.5% custodian fee
    NetDisbursed := request.LeverageAmount * (processingFee + custodianFee)

   
    //Save leverage requests
    leverage := models.LeverageTransaction{
        UserID:              request.UserID,
        GoldHoldingID:       request.GoldHoldingID,
        LeverageAmount:      request.LeverageAmount,
        TenureMonths:        request.TenureMonths,
        NetDisbursed:        NetDisbursed,
        ProcessingFee:       processingFee,
        CustodianFee:        custodianFee,
        InterestRate:        28.0, // Hardcoded for now, replace with actual interest rate
        Status:               "pending",
    }
    config.DB.Create(&leverage)
    

    c.JSON(http.StatusOK, gin.H{
        "message": "Leverage request submitted successfully!","leverage_id": leverage.ID})
}

//Placeholder controller functions to be implemented

func ApproveLeverage(c *gin.Context){
    c.JSON(http.StatusOK,gin.H{"message":"Leverage approved sucessfully"})


}

func GetPendingRequests(c *gin.Context){
    c.JSON(http.StatusOK,gin.H{"message":"Fetched pending requests"})

}

func ForceLiquidate(c *gin.Context){
    c.JSON(http.StatusOK,gin.H{"message":"Leverage liquidated successfully"})
}