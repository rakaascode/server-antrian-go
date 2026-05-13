package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireSuperAdmin memastikan hanya global admin (admin tanpa cabang_id) yang bisa akses.
// Admin cabang (yang punya cabang_id) tidak diizinkan.
func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Hanya super admin yang diizinkan",
			})
			return
		}

		// Jika cabang_id ada di context → ini admin cabang, bukan super admin
		if _, hasCabang := c.Get("cabang_id"); hasCabang {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Endpoint ini khusus super admin (global admin)",
			})
			return
		}

		c.Next()
	}
}
