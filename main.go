package main

import (
	"dukia-leverage-api/routes"
	"dukia-leverage-api/services"
	"log"
    "os"
	"reflect"

	"dukia-leverage-api/config"
	"github.com/gin-gonic/gin"

	"dukia-leverage-api/models"
	"fmt"
)

func main() {
	// Load the configuration
	config.ConnectDatabase()

	// Print all registered models
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.GoldHolding{},
		&models.LeverageTransaction{},
		&models.MarginCall{},
		&models.Loan{},
	}
	fmt.Println("Checking all registered models:")
	for _, model := range modelsToMigrate {
		fmt.Println(" - ", reflect.TypeOf(model))
	}

	fmt.Println("Checking if AutoMigraion is running...")

	// Migrate the models

	err := config.DB.AutoMigrate(&models.User{}, &models.GoldHolding{}, &models.LeverageTransaction{}, &models.MarginCall{}, &models.Loan{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("AutoMigrate executed successfully!")

	//Trigger liqidation process
	services.MonitorLTVandLiquidate()

	// Initialize Gin engine
	router := gin.Default()

	// Register routes
	routes.UserRoutes(router)
	routes.LeveragingRoutes(router)
	routes.AdminRoutes(router)

    port := os.Getenv("PORT")
    if port == "" {
    port = "8080" // Default port for local development
    }
    router.Run(":" + port)


	// Start the server
	//router.Run(":8080")

}
