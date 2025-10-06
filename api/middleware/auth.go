package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/suk-chanthea/ezra/service"
)

// JWTMiddleware validates JWT token from Authorization header
func JWTMiddleware(authService service.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }

        tokenStr := parts[1]
        token, err := authService.ValidateToken(tokenStr)
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            sub, ok := claims["sub"].(float64) // JWT numbers are float64
            if !ok {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
                c.Abort()
                return
            }
            c.Set("user_id", uint(sub)) // now guaranteed to be >0
            c.Set("username", claims["username"])
            c.Set("role", claims["role"])
        }

        c.Next()
    }
}

