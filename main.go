package main

import (
	"dukia-leverage-api/routes"
	"dukia-leverage-api/services"
	"log"
    "os"
	"reflect"

	"dukia-leverage-api/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"dukia-leverage-api/models"
	"fmt"
)

func main() {
	// Load the configuration
	config.ConnectDatabase()

    // Initialize Gin engine
	router := gin.Default()

	//Enable CORS
	router.Use(cors.Default())

	// Define API group
    api := router.Group("/api")

	api.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Welcome to Dukia API"})
    })

    // ✅ Pass `router`, not `api`, to route functions
    routes.UserRoutes(api)
    routes.LeveragingRoutes(api)
    routes.AdminRoutes(api)


    port := os.Getenv("PORT")
    if port == "" {
    port = "10000" // Default port for local development
    }
	router.Run(":" + port) 

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

	






}
