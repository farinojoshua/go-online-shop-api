package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// terapkan key
		key := "123"

		// ambil auth dari request header
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "akses tidak diizinkan"})
			c.Abort()
			return
		}

		if auth != key {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "akses tidak diizinkan"})
			c.Abort()
			return
		}

		// lanjutkan ke handler
		c.Next()
	}
}
