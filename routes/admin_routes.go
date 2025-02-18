package routes 

import(
	"dukia-leverage-api/controllers"
	"github.com/gin-gonic/gin"
    "dukia-leverage-api/utils"
)

func AdminRoutes(router *gin.Engine) {
	admin :=router.Group("/admin")
	admin.Use(AdminAuthMiddleware())
	{
		admin.POST("/approve-leverage", controllers.ApproveLeverage)
        admin.GET("/leverage-requests", controllers.GetPendingRequests)
        admin.POST("/liquidate/:id", controllers.ForceLiquidate)
		}
        
	}

	///AdminAuthMiddleware checks if the user is an admin
	func AdminAuthMiddleware() gin.HandlerFunc {
		return func(c *gin.Context) {
            token := c.GetHeader("Authorization")
            if token == "" {
                c.JSON(401, gin.H{"error": "Authorization token is missing"})
                c.Abort()
                return
            }

		    // Verify the token and extract the user's role from it
            // If the user's role is not "admin", return a 403 Forbidden response
            // else, continue with the request
            claims, err := utils.ParseToken(token)
			if err != nil {
                c.JSON(403, gin.H{"error": "Invalid token"})
                c.Abort()
                return
            }
            role, ok := (*claims)["role"].(string)
            if !ok || role != "admin" {
                c.JSON(403, gin.H{"error": "Unauthorized access"})
                c.Abort()
                return
            }

			c.Next()  // If the user is an admin, proceed to the next middleware or route	
        }
	}
