package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"fmt"

	"gorm.io/driver/postgres"
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
	os.Setenv("TEST_ENV", "true")// Ensure mock gold price is used
	
	//initiate a test database
	dsn:= "host=localhost user=postgres password=postgres dbname=dukia_leverage port=5432 sslmode=disable TimeZone=Africa/Lagos"
	testDB, err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to test database" ,err)
		os.Exit(1)
    }
	//Wrap everying in a transaction
	config.DB = testDB
	tx := config.DB.Begin()

	//Truncate and reset IDs	
	tx.Exec("TRUNCATE TABLE gold_holdings RESTART IDENTITY CASCADE")

	//Run migrations for test database
	err = config.DB.AutoMigrate(
		&models.User{},
		&models.GoldHolding{},
		&models.MarginCall{},
		&models.Loan{},
		&models.LeverageTransaction{},
	)
	if err != nil {
		fmt.Println("Migration failed: ",err)
        os.Exit(1)
    }
	//Insert fresh test data
	testGoldHoldings := []models.GoldHolding{
	{ID: 1, UserID: 1, Quantity: 100, CurrentValue: 5000, GoldType: "24k"},
	{ID: 2, UserID: 2, Quantity: 10, CurrentValue: 500, GoldType: "24k"},
	}

	for _, gold := range testGoldHoldings {
    tx.Create(&gold)
    }

	exitVal := m.Run()
	tx.Rollback() // Rollback transaction on exit

	//Drop tables after tests to keep database clean	
	//config.DB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	os.Exit(exitVal)
	
}

func setupMockRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	//req, _ := http.NewRequest("POST", "/leverage/apply", bytes.NewBuffer(jsonValue))

	r.POST("/leverage/apply",controllers.ApplyLeverage)
	return r
}

func TestApplyLeverage_Success(t *testing.T) {
	// Cleanup before inserting new data
config.DB.Exec("DELETE FROM gold_holdings WHERE id IN (1, 2)")

// Insert test gold holdings
goldHolding1 := models.GoldHolding{
    ID:       1,
    UserID:   1,
    Quantity: 100,
    GoldType: "24k",
}
config.DB.Create(&goldHolding1)

goldHolding2 := models.GoldHolding{
    ID:       2,
    UserID:   2,
    Quantity: 10,
    GoldType: "24k",
}
config.DB.Create(&goldHolding2)


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
	config.DB.Create(&models.GoldHolding{ID: 2, UserID: 2, Quantity: 10})

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

	if w.Code == http.StatusNotFound {
		t.Errorf("Gold Holding not found, Ensure test data exists.")
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
