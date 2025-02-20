package mocks

import (
	"dukia-leverage-api/models"
	//"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	
)

//mockDB struct to simulate a database

type MockDB struct {
	DB *gorm.DB
}

//Override Where() to return a valid *gorm.DB
func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.DB    //Always return a valid *gorm.DB instance
}
    
//Override First() to return a mock GoldHolding data and prevent nil ointer derefrences
func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	if goldHolding, ok := dest.(*models.GoldHolding); 
	ok {
		*goldHolding = models.GoldHolding{
		ID: 			1,
		UserID: 		1,
		Quantity:		100,    // 100g of gold for testing
		CurrentValue: 	5000,  // Example value
		GoldType: 		"24k",
	}
}
	
    return &gorm.DB{}
}