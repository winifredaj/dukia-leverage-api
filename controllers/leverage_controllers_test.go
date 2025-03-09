package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"dukia-leverage-api/config"
	"dukia-leverage-api/controllers"

	"dukia-leverage-api/controllers/mocks"
	"dukia-leverage-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

//Global variable for database transaction rollback

var testDB *gorm.DB
var tx *gorm.DB

func TestMain(m *testing.M) {
	os.Setenv("TEST_ENV", "true") // Ensure mock gold price is used

	//initiate a test database
	dsn := "host=localhost user=postgres password=postgres dbname=dukia_leverage port=5432 sslmode=disable TimeZone=Africa/Lagos"
	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to test database", err)
	}
	log.Println("Connected to Test database")

	//Run migrations for test database
	err = testDB.AutoMigrate(
		&models.User{},
		&models.GoldHolding{},
		&models.MarginCall{},
		&models.Loan{},
		&models.LeverageTransaction{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	//Set testDB as global daatabase for tests
	config.DB = testDB

	//Run tests
	exitVal := m.Run()

	//Clean up
	sqlDB, _ := testDB.DB()
	sqlDB.Close()

	os.Exit(exitVal)

}

func setupMockRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	//req, _ := http.NewRequest("POST", "/leverage/apply", bytes.NewBuffer(jsonValue))
	r.POST("/leverage/apply", controllers.ApplyLeverage)
	return r
}

func runTestWithTransaction(t *testing.T, testFunc func(*testing.T)) {
	t.Helper()

	//Ensure testDB is available
	if testDB == nil {
		t.Fatal("Test setup failed: test database connection is not initialized")
	}

	//Start a new transaction and rollback it on completion of the test
	tx = testDB.Begin()
	if tx.Error != nil {
		t.Fatalf("Failed to begin transction: %v", tx.Error)
	}

	defer tx.Rollback() //Ensures test data is cleaned up

	//Assign to global config.DB
	config.DB = tx
	testFunc(t) //Run the actual test
}

func TestApplyLeverage_Success(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T) {
		//Insert test gold holdings with unique values
		tx.Exec("TRUNCATE TABLE gold_holdings RESTART IDENTITY CASCADE")
		tx.Create(&models.GoldHolding{ID: 1, UserID: 1, Quantity: 100, GoldType: "Bars"})

		// Cleanup before inserting new data
		//config.DB.Exec("DELETE FROM gold_holdings WHERE id IN (1, 2)")

		//Use MockDB instead  real DB
		mockDB := &mocks.MockDB{DB: config.DB}
		config.DB = mockDB.DB //Override global DB instance

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
	})
}

func TestApplyLeverage_InsufficientGold(t *testing.T) {
	runTestWithTransaction(t, func(t *testing.T) {

		//Ensure the transaction is valid before executing SQL statement
		if tx.Error != nil {
			t.Fatalf("Invalid transction before inserting test data: %v", tx.Error)
		}
		//Truncate table with error handling
		if err := tx.Exec("TRUNCATE TABLE gold_holdings RESTART IDENTITY CASCADE").Error; err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}

		//Insert test data
		err := tx.Create(&models.GoldHolding{ID: 2, UserID: 2, Quantity: 10, GoldType: "Pool"}).Error
		if err != nil {
			t.Fatalf("Failedto insert test data: %v", err)
		}

		mockDB := &mocks.MockDB{DB: config.DB}
		config.DB = mockDB.DB //Override global DB instance

		// Setup the router
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

		//Debugging output
		if w.Code == http.StatusNotFound {
			t.Errorf("Gold Holding not found, Ensure test data exists.")
		}

		//Check expected response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
