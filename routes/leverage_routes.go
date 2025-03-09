package routes

import (
	"dukia-leverage-api/controllers"
	"dukia-leverage-api/middleware"
	"github.com/gin-gonic/gin"
)

func LeveragingRoutes(router *gin.Engine) {
	leverageGroup := router.Group("/leverage").Use(middleware.JWTMiddleware("use"))
	{
		//leverageGroup.GET("/status/:user_id", controllers.GetLeverage)
		leverageGroup.GET("/status/:user_id", controllers.GetLeverageDetails)
		leverageGroup.POST("/apply", controllers.ApplyLeverage)
		leverageGroup.DELETE("cancel/:id", controllers.CancelLeverageRequest)
		leverageGroup.GET("/simulate-margin-call/:id", controllers.SimulateMarginCall)

	}
}
