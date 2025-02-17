package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"dukia-leverage-api/config"
	"dukia-leverage-api/models"
    "dukia-leverage-api/services"
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
}
    
    //Check if user has sufficient balance
    func CheckEligibility(c *gin.Context){
        userID := c.Query("user_id")
        if userID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
            return
        }

        //Fetch user's gold holdings from the database and check if they are already in the database
        goldHoldings, err := models.GetGoldHoldingsByUserID(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to retrieve gold holdings"})
            return
            }

            totalGoldWeight:= goldHoldings.TotalWeight() // Implement this method

            if totalGoldWeight < 50 { // Implement this method
                c.JSON(http.StatusOK, gin.H{
                    "eligible": false,
                    "max_loan_amount" : 0,
                    "message": "Insufficient gold balance. Minimumo 50grams required.",
                })
                return
            }
            //Fetch current market price of gold
            goldPrice, err := services.GetCurrentGoldPrice()
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to retrieve market price"})
                return
            }

            //Calculate total value of user's gold holdings
            totalGoldValue := totalGoldWeight * goldPrice

            //Calculate maximum loan amountbased on LTV ratio

            ltvRatio := 0.75
            maxLoanAmount := totalGoldValue * ltvRatio

            c.JSON(http.StatusOK, gin.H{
                "eligible": true,
                "max_loan_amount" : maxLoanAmount,
                "message": "User is eligible for leverage application.",
            })
        
    

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