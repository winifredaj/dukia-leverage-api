package routes

import (
    "dukia-leverage-api/controllers"
	"github.com/gin-gonic/gin"
    
)

func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/auth")
	{
        userGroup.POST("/register", controllers.RegisterUser)
        userGroup.POST("/login", controllers.LoginUser)
    }
}