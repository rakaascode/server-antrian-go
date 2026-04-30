package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakaascode/server-antrian-go.git/internal/auth"
)

// AuthMiddleware validasi JWT dan inject claims ke context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false, "message": "Authorization header diperlukan",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false, "message": "Format: Bearer <token>",
			})
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false, "message": "Token tidak valid atau kadaluarsa",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// Inject cabang_id jika ada (untuk admin)
		if claims.CabangID != nil {
			c.Set("cabang_id", *claims.CabangID)
		}

		c.Next()
	}
}

// OptionalAuth — inject claims jika ada token, tapi tidak wajib
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			if claims, err := auth.ParseToken(parts[1]); err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("role", claims.Role)
				if claims.CabangID != nil {
					c.Set("cabang_id", *claims.CabangID)
				}
			}
		}
		c.Next()
	}
}
