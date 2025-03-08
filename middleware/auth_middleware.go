package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"dukia-leverage-api/utils"

	"github.com/gin-gonic/gin"
	//"github.com/golang-jwt/jwt/v4"
)

// Load the JWT secret key
var JWTKey = []byte(getJWTKey())

func getJWTKey() string {
	secret := os.Getenv("USER_JWT_SECRET")
	if secret == "" {
		fmt.Println("Error: USER_JWT_SECRET environment variable not set.")
	}
	return secret
}

// JWTMiddleware enforces authentication and authorization
func JWTMiddleware(expectedRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the Authorization header exists
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("JWTMiddleware: Missing Authorization header for", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
			c.Abort()
			return
		}

		// Ensure the Authorization header follows the correct format (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Println("JWTMiddleware: Invalid token format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
	 // Load the JWT secret key from environment variable USER_JWT_SECRET
		userSecretKey := []byte(os.Getenv("USER_JWT_SECRET"))

		// Validate token using utils.ValidateToken
		token, claims, err := utils.ValidateToken(tokenString, userSecretKey)
		if err != nil ||token == nil || !token.Valid {
			log.Println("JWTMiddleware: Invalid token for", c.Request.URL.Path, "-", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims properly
		//claimsMap, ok := claims.(jwt.MapClaims)
		//if !ok {
			//log.Println("JWTMiddleware: Invalid token claims for", c.Request.URL.Path)
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			//c.Abort()
			//return
		//}

		// Extract user role
		role, roleExists := claims["role"].(string)
		if !roleExists {
			log.Println("JWTMiddleware: Role not found in token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}

		// Enforce role-based access
		if role != expectedRole {
			log.Println("JWTMiddleware: Access denied - expected role:", expectedRole, "but found:", role)
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}

		// Extract user ID safely
		if userIDFloat, exists := claims["user_id"].(float64); exists {
			userID := int(userIDFloat)
			c.Set("user_id", userID)
		}

		// Extract email safely
		if email, exists := claims["email"].(string); exists {
			c.Set("email", email)
		
		}

		// Set user role in context
		c.Set("role", role)

		// Proceed to next handler
		c.Next()
	}
}
