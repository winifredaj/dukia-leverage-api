package routes

import (
	"dukia-leverage-api/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api")
	{
		userGroup.POST("/auth/register", controllers.RegisterUser)
		userGroup.POST("/auth/login", controllers.LoginUser)
	}
}
