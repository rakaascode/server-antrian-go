package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedOrigins = map[string]bool{
	"https://crm-lautan-teduh.vercel.app": true,
	"https://sistem-antrean.vercel.app":   true,
	"http://localhost:3000":               true,
	"http://127.0.0.1:3000":               true,
}

// isAllowedOrigin memeriksa apakah origin termasuk domain yang diizinkan
func isAllowedOrigin(origin string) bool {
	if allowedOrigins[origin] {
		return true
	}

	// Dukungan untuk vercel preview deployments
	if strings.HasSuffix(origin, "-rakaascodes-projects.vercel.app") ||
		strings.HasPrefix(origin, "https://crm-lautan-teduh-") ||
		strings.HasPrefix(origin, "https://sistem-antrean-") {
		return true
	}

	return false
}

// CORSMiddleware mengelola header Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" && isAllowedOrigin(origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		if c.Request.Method == http.MethodOptions {
			if origin != "" && !isAllowedOrigin(origin) {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
