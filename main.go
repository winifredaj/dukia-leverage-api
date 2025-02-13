package main

import(
	"dukia-leverage-api/config"
	"dukia-leverage-api/routes"
	"github.com/gin-gonic/gin"

    "dukia-leverage-api/models"
    "fmt"

)

func main() {
	// Load the configuration
    config.ConnectDatabase()

    // Migrate the models
    config.DB.AutoMigrate(&models.User{}, &models.GoldHolding{}, &models.LeverageTransaction{},  &models.MarginCall{})
    fmt.Println("Database migrations completed!")

    // Initialize Gin engine
    router := gin.Default()

    // Register routes
    routes.UserRoutes(router)
	routes.LeveragingRoutes(router)

    // Start the server
    router.Run(":8080")

    
}