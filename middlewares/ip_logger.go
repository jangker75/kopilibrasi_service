package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func IPLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ipAddress := c.ClientIP()
		log.Printf("IP Address: %s", ipAddress)

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		// Optionally, you can store the IP address in the context
		c.Set("ClientIP", ipAddress)

		// Continue to next middleware/handler
		c.Next()
	}
}
