package middleware

import (
	"net/http"
	"os"
	"strings"

	"osohub/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWTSecret() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

// AuthMiddleware valida el JWT y pone el user_id en el contexto
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header", "documentation": "https://docs.osohub.com/auth#jwt"})
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "documentation": "https://docs.osohub.com/auth#jwt"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims", "documentation": "https://docs.osohub.com/auth#jwt"})
			c.Abort()
			return
		}
		c.Set("user_id", claims["user_id"].(string))
		c.Next()
	}
}

// AdminOnly middleware ensures only admins can access the endpoint
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header", "documentation": "https://docs.osohub.com/auth"})
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "documentation": "https://docs.osohub.com/auth"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		role, okRole := claims["role"].(string)
		if !ok || !okRole || role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admins only", "documentation": "https://docs.osohub.com/auth#roles"})
			c.Abort()
			return
		}
		c.Next()
	}
}
