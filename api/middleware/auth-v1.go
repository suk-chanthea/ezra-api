package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware checks a simple header token or database user
func AuthMiddleware(db *gorm.DB, token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example: simple token check
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || authHeader != token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Optional: check user in DB
		// userID := c.GetHeader("X-User-ID")
		// db.First(&user, userID)

		c.Next()
	}
}
