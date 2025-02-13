package middleware

import (
    "fmt"
    "net/http"
    "strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWTKey is the secret key used to sign and verify JWT tokens

var JWTKey = []byte("secretkey")

// JWTMiddleware is a middleware function that checks for a valid JWT token in the request header	
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {	
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()

		    return
        }
		
        // Extract the token from the Authorization header and parse it
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
                return nil, fmt.Errorf("Unexpected signing method")
            }
			return JWTKey, nil
			})
		
		// If the token is invalid, return an unauthorized response
        if err != nil ||!token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			    return

		}//proceed to next middleware
				c.Next()
	}

}
