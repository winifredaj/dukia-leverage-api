package controllers

import (
   "dukia-leverage-api/controllers/mocks"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	
	"dukia-leverage-api/models"
	"dukia-leverage-api/config"
	

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupMockRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/leverage/apply", controllers.ApplyLeverage)
	return r
}

func TestApplyLeverage_Success(t *testing.T) {

	//Use MockDB instead  real DB
	mockDB := &MockDB{}
	config.DB = mockDB  //Override global DB instance

	r := setupMockRouter()

	mockRequest := models.LeverageTransaction{
		UserID:         1,
		GoldHoldingID:  1,
		LeverageAmount: 5000,
		TenureMonths:   12,
	}

	jsonValue, _ := json.Marshal(mockRequest)
	req, _ := http.NewRequest("POST", "/leverage/apply", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestApplyLeverage_InsufficientGold(t *testing.T) {

	mockDB := &MockDB{}
	config.DB = mockDB  //Override global DB instance

	r := setupMockRouter()

	mockRequest := models.LeverageTransaction{
		UserID:         2,
		GoldHoldingID:  2,
		LeverageAmount: 5000,
		TenureMonths:   12,
	}

	jsonValue, _ := json.Marshal(mockRequest)
	req, _ := http.NewRequest("POST", "/leverage/apply", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
