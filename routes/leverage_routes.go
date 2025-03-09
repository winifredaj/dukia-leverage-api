package routes

import (
	"dukia-leverage-api/controllers"
	"dukia-leverage-api/middleware"
	"github.com/gin-gonic/gin"
)

func LeveragingRoutes(router *gin.Engine) {
	leverageGroup := router.Group("/api").Use(middleware.JWTMiddleware("use"))
	{
		//leverageGroup.GET("/status/:user_id", controllers.GetLeverage)
		leverageGroup.GET("/leverage/status/:user_id", controllers.GetLeverageDetails)
		leverageGroup.POST("/leverage/apply", controllers.ApplyLeverage)
		leverageGroup.DELETE("/leverage/cancel/:id", controllers.CancelLeverageRequest)
		leverageGroup.GET("/leverage/simulate-margin-call/:id", controllers.SimulateMarginCall)

	}
}
