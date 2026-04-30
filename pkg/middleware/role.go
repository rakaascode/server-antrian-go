package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole membatasi akses berdasarkan role user
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false, "message": "Role tidak ditemukan",
			})
			return
		}

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false, "message": "Akses ditolak. Anda tidak memiliki izin",
		})
	}
}
