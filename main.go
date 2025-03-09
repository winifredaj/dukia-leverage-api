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

    // Initialize Gin engine
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Define API group
    //api := router.Group("/api")

    // ✅ Pass `router`, not `api`, to route functions
    routes.UserRoutes(router)
    routes.LeveragingRoutes(router)
    routes.AdminRoutes(router)


    port := os.Getenv("PORT")
    if port == "" {
    port = "10000" // Default port for local development
    }

	log.Println("Starting server on port " + port)
    
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Dukia Leverage API is running"})
    })

	router.GET("/debug/routes", func(c *gin.Context) {
		routes := router.Routes()
		for _, r := range routes {
			fmt.Println("Registered Route:", r.Method, r.Path)
		}
		c.JSON(200, gin.H{"routes": routes})
	})
	

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

   
    
	//router.Run(":" + port) 
	router.Run("0.0.0.0:" + port) 

	// Start the server
	// router.Run(":8080")

}
