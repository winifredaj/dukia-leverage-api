package routes

import (
	"dukia-leverage-api/controllers"
	"dukia-leverage-api/middleware"
  "github.com/gin-gonic/gin"
)

func LeveragingRoutes(router *gin.Engine) {
	leverageGroup := router.Group("/leverage").Use(middleware.JWTMiddleware())
	{
		leverageGroup.GET("/status/:user_id", controllers.GetLeverage)
        leverageGroup.POST("/apply", controllers.ApplyLeverage)
    }
}