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

/// GetLeverageDetails retrieves leverage transaction details for a user
 func GetLeverageDetails(c *gin.Context){
    userID, exists := c.Get("userID")
    if!exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error":"Unauthorized"})
        return
    }

    leverageID := c.Param("id")
    var leverage models.LeverageTransaction

    if err := config.DB.Where("id =? AND user_id =?", leverageID, userID).First(&leverage).Error; err != nil {
        c.JSON(http.StatusNotFound,gin.H{"error":"Leverage transaction not found"})
        return
    }

    c.JSON(http.StatusOK, leverage)
}


// Apply for Leverage
 func ApplyLeverage(c *gin.Context) {
    var request struct {
        UserID              uint    `json:"user_id"`
        GoldHoldingID       uint    `json:"gold_holding_id"`
        LeverageAmount      float64 `json:"leverage_amount"`
        TenureMonths        int     `json:"tenure_months"`     
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

    //Validate that gold_holding_id is not zero
    if request.GoldHoldingID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gold_holding_id, must be greater than zero"})
        return
    }

    if err := config.DB.Where("id = ?", request.GoldHoldingID).First(&goldHolding).Error; err != nil {
        fmt.Println("DEBUG:GoldHolding not found for ID:", request.GoldHoldingID, "Error:", err)
        c.JSON(http.StatusNotFound,gin.H{"error":"Gold Holding not found"})
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
    netDisbursed := request.LeverageAmount - (processingFee + custodianFee)

    if netDisbursed <= 0{
        c.JSON(http.StatusBadRequest,gin.H{"error":"Insuffiient leverage net disbursed amount "})
        return
    }

    //Save leverage requests
    leverage := models.LeverageTransaction{
        UserID:              request.UserID,
        GoldHoldingID:       request.GoldHoldingID,
        LeverageAmount:      request.LeverageAmount,
        TenureMonths:        request.TenureMonths,
        NetDisbursed:        netDisbursed,
        ProcessingFee:       processingFee,
        CustodianFee:        custodianFee,
        InterestRate:        28.0, // Hardcoded for now, replace with actual interest rate
        Status:               "pending",
    }
    
    if err := config.DB.Create(&leverage).Error; err != nil {
        fmt.Println(" Failed to save leverage transaction", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to save leverage request"})
        return
    }
    

    c.JSON(http.StatusOK, gin.H{
        "message": "Leverage request submitted successfully!",
        "leverage_id": leverage.ID,
    })
}

//CancelLeverageRequest allows a user to cancel a pending leverage request 
 func CancelLeverageRequest(c *gin.Context){
    userID, exists := c.Get("userID")
    if!exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error":"Unauthorized"})
        return
    }

    leverageID := c.Param("id")
    var leverage models.LeverageTransaction

    if err := config.DB.Where("id =? AND user_id =?", leverageID, userID).First(&leverage).Error; err != nil {
        c.JSON(http.StatusNotFound,gin.H{"error":"Leverage request not found"})
        return
    }

    if leverage.Status != "pending"{
    c.JSON(http.StatusBadRequest, gin.H{"message":"Leverage request cannot be cancelled as it is already processed"})
    return
    }

    if err := config.DB.Delete(&leverage).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to cancel leverage request"})
        return
    }
    c.JSON(http.StatusOK,gin.H{"message":"Leverage request cancelled successfully"})
}


// SimulateMarginCall checks if a leverage transaction requires a margin call
 func SimulateMarginCall(c *gin.Context){
    leverageID := c.Param("id")
    var leverage models.LeverageTransaction

    // Fetch leverage transactionn
    if err := config.DB.Where("id =?", leverageID).First(&leverage).Error; err != nil {
        c.JSON(http.StatusNotFound,gin.H{"error":"Leverage transaction not found"})
        return
    }

    //Fetch related margin call
    var marginCall models.MarginCall
    if err := config.DB.Where("leverage_transaction_id =?", leverageID).First(&marginCall).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error":"No margin call associated with this leverage transaction"})
        return
    }

    currentMarketPrice, err := utils.GetCurrentGoldPrice() //"https://api.dukiapreciousmetals.co/api/price/products7"
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to retrieve market price"})
        return
    }

    marginCallThreshold :=marginCall.RequiredCollateral

    if currentMarketPrice > marginCallThreshold {
        c.JSON(http.StatusOK,gin.H{"message":"Margin call triggered", "current_price":currentMarketPrice})
        return
    } else{
    c.JSON(http.StatusOK,gin.H{"message":"No Margin call triggered", "current_price": currentMarketPrice})
    }
}



///////////////////////////////////////////////ADMIN FUNCTIONS///////////////////////////////////////////////////////////

//ApproveLeverageRequest approves a pending leverage 
func ApproveLeverageRequest(c *gin.Context){
    leverageID := c.Param("id")

    var leverage models.LeverageTransaction
    if err := config.DB.Where("id =?", leverageID).First(&leverage).Error; err != nil {
        c.JSON(http.StatusNotFound,gin.H{"error":"Leverage request not found"})
        return
    }

    if leverage.Status != "pending" {
        c.JSON(http.StatusBadRequest, gin.H{"message":"Only pending Leverage request can be approved."})
        return
    }

    //Update leverage statusto "approved"
    leverage.Status = "approved"
    if err := config.DB.Save(&leverage).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to approve leverage request"})
        return
    }

    c.JSON(http.StatusOK,gin.H{"message":"Leverage request approved sucessfully"})

}

//// RejectLeverageRequest allows an admin to reject a pending leverage request
 func RejectLeverageRequest(c *gin.Context){
    leverageID := c.Param("id")
    var leverage models.LeverageTransaction

    if err := config.DB.Where("id =?", leverageID).First(&leverage).Error; err != nil {
        c.JSON(http.StatusNotFound,gin.H{"error":"Leverage request not found"})
        return
    }

    if leverage.Status != "pending" {
        c.JSON(http.StatusBadRequest, gin.H{"message":"Only pending Leverage request can be rejected."})
        return
    }

    leverage.Status = "rejected"
    if err := config.DB.Save(&leverage).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to reject leverage request"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"messagr": "Leverage request rejected successfully"})
}


func GetPendingRequests(c *gin.Context){
    var pendingLeverages []models.LeverageTransaction
    if err:= config.DB.Where("status = ?", "pending").Find(&pendingLeverages).Error; err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pending requests"})
        return
    }
    c.JSON(http.StatusOK,gin.H{"pending_requests":pendingLeverages})

}

// ManageMarginCall allows an admin to resolve margin call issues
 func ManageMarginCall(c *gin.Context){
    marginCallID := c.Param("id")
    var marginCall models.MarginCall

    if err := config.DB.Where("id =?", marginCallID).First(&marginCall).Error; err != nil {
        c.JSON(http.StatusNotFound,gin.H{"error":"Margin call not found"})
        return
    }

    if marginCall.Status != "pending"{
        c.JSON(http.StatusBadRequest, gin.H{"error":"Margin call already resolved or defaulted."})
        return
    }

    marginCall.Status = "resolved"
    if  err := config.DB.Save(&marginCall).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"erro": "Failed to update manage margin call status"})
        return  
    }

    c.JSON(http.StatusOK, gin.H{"message":"Margin call resolved successfully"})
}


// CheckDefaultedLoans retrieves all defaulted leverage transactions
func CheckDefaultedLoans(c *gin.Context){
    var defaultedLeverages []models.LeverageTransaction

    if err:= config.DB.Where("status = ?", "defaulted").Find(&defaultedLeverages).Error; err!= nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve defaulted loans"})
        return
    }

    c.JSON(http.StatusOK, defaultedLeverages)
}


func ForceLiquidate(c *gin.Context){

    leverageID := c.Param("id")
    var leverage models.LeverageTransaction
    if err := config.DB.Where("id =?", leverageID).First(&leverage).Error; err != nil {
        c.JSON(http.StatusNotFound,gin.H{"error":"Leverage transaction not found"})
        return
    }
    if leverage.Status != "active"{
        c.JSON(http.StatusBadRequest, gin.H{"error":"Only ative leverage transaction can be liquidated."})
        return
    }

    //Mark leverage as liquidated
    leverage.Status = "liquidated"
    if  err := config.DB.Save(&leverage).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to liquidate leverage transaction"})
        return  
    }

    c.JSON(http.StatusOK,gin.H{"message":"Leverage liquidated successfully", "leverage":leverage})
}