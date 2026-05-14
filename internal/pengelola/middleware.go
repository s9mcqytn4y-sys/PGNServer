package pengelola

import (
	"net/http"
	"pgn-server/pkg/utilitas"
	"strings"

	"github.com/gin-gonic/gin"
)

// AutentikasiMiddleware memverifikasi Bearer token dari Header
func AutentikasiMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, utilitas.FormatRespons(false, "Tidak terotorisasi, token tidak ditemukan atau format salah", nil))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		// Implementasikan validasi token sebenarnya di sini (misal: verifikasi JWT)
		if token != "token-rahasia-pgn" {
			c.JSON(http.StatusUnauthorized, utilitas.FormatRespons(false, "Tidak terotorisasi, token tidak valid", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}
