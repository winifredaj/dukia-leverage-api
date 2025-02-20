package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"testing"

	"dukia-leverage-api/config"
	"dukia-leverage-api/controllers"
	"dukia-leverage-api/controllers/mocks"
	"dukia-leverage-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	
)

func TestMain(m *testing.M) {
	//initiate a tesr database
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"),&gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database" +err.Error())
    }

	config.DB = testDB
	exitVal := m.Run()
	os.Exit(exitVal) // 
}

func setupMockRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	//req, _ := http.NewRequest("POST", "/leverage/apply", bytes.NewBuffer(jsonValue))

	r.POST("/leverage/apply",controllers.ApplyLeverage)
	return r
}

func TestApplyLeverage_Success(t *testing.T) {

	//Use MockDB instead  real DB

	mockDB := &mocks.MockDB{DB: config.DB}
	config.DB = mockDB.DB  //Override global DB instance

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

	mockDB := &mocks.MockDB{DB: config.DB}
	config.DB = mockDB.DB  //Override global DB instance

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
