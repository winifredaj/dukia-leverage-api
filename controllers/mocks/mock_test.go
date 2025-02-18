package mocks

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"dukia-leverage-api/models"
)

//mockDB struct to simulate a database

type MockDB struct {
    mock.Mock
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	mockGoldHolding := models.GoldHolding{
		ID: 1,
		UserID: 1,
		Quantity: 100,    // 100g of gold for testing
		CurrentValue: 5000,  // Example value
		GoldType: "24k",
	}
	//Check if destination is of type GoldHoldng and assign mock data
	if goldHolding, ok := dest.(*models.GoldHolding); ok {
		*goldHolding = mockGoldHolding

	}
    return &gorm.DB{}
}