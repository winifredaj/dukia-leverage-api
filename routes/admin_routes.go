package routes

import (
	"log"
	"net/http"
	"os"
	"strings"

	"dukia-leverage-api/controllers"
	"dukia-leverage-api/utils"
	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v4"
)

func AdminRoutes(router *gin.Engine) {
	//Login route without middleware
	router.POST("/admin/auth/login", controllers.LoginAdmin)

	//Admin routes with middleware
	admin := router.Group("/admin")
	admin.Use(AdminAuthMiddleware())
	{

		admin.POST("/approve-leverage/:id", controllers.ApproveLeverageRequest)
		admin.POST("/reject-leverage/:id", controllers.RejectLeverageRequest)
		admin.GET("/leverage-requests", controllers.GetPendingRequests)
		admin.POST("/liquidate-leverage", controllers.ForceLiquidate)
		admin.POST("/resolve-margin-call/:id", controllers.ManageMarginCall)
		admin.GET("/defaulted-leverages", controllers.CheckDefaultedLoans)

	}

}

///AdminAuthMiddleware checks if the user is an admin
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("AdminAuthMiddleware:Authorization header missing for", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			log.Println("AdminAuthMiddleware: Invalid token format for", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		tokenString := tokenParts[1]

		// Verify the token and extract the user's role from it
		// If the user's role is not "admin", return a 403 Forbidden response
		// else, continue with the request
		adminSecretKey := []byte(os.Getenv("ADMIN_JWT_SECRET"))

		token, claims, err := utils.ValidateToken(tokenString, adminSecretKey)
		if err != nil || token == nil || !token.Valid {
			log.Println("AdminAuthMiddleware: Error validating token for", c.Request.URL.Path, "error", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			log.Println("AdminAuthMiddleware: Unauthorized role ", role, c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}

		c.Next() // If the user is an admin, proceed to the next middleware or route
	}
}
