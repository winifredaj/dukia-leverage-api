package main

import (
	"dukia-leverage-api/config"
	"dukia-leverage-api/routes"
    "dukia-leverage-api/services"
    "os"
	"log"
	"reflect"

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
    }
    fmt.Println("Checking all registered models:")
    for _, model := range modelsToMigrate {
        fmt.Println(" - ", reflect.TypeOf(model))
        }

    fmt.Println("Checking if AutoMigraion is running...")

    // Migrate the models

    err:= config.DB.AutoMigrate(&models.User{}, &models.GoldHolding{}, &models.LeverageTransaction{},  &models.MarginCall{})
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

    // Start the server
    router.Run(":8080")
       
}