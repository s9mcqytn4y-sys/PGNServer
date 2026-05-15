package pengelola

import (
	"context"
	"log/slog"
	"net/http"
	"pgn-server/pkg/utilitas"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
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

// LogSlogMiddleware mengimplementasikan structured logging menggunakan slog
func LogSlogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mulai := time.Now()
		path := c.Request.URL.Path
		kueri := c.Request.URL.RawQuery

		c.Next()

		akhir := time.Since(mulai)
		status := c.Writer.Status()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				slog.Error("Kesalahan request",
					"jalur", path,
					"status", status,
					"error", e,
				)
			}
		} else {
			slog.Info("Request selesai",
				"metode", c.Request.Method,
				"jalur", path,
				"kueri", kueri,
				"status", status,
				"durasi", akhir,
				"ip", c.ClientIP(),
			)
		}
	}
}

// LimitRequestMiddleware membatasi jumlah request per detik (Rate Limiting)
func LimitRequestMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, utilitas.FormatRespons(false, "Terlalu banyak permintaan, silakan coba lagi nanti", nil))
			c.Abort()
			return
		}
		c.Next()
	}
}

// TimeoutMiddleware membatasi waktu eksekusi setiap request
func TimeoutMiddleware(durasi time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), durasi)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		selesai := make(chan struct{}, 1)
		go func() {
			c.Next()
			selesai <- struct{}{}
		}()

		select {
		case <-selesai:
			return
		case <-ctx.Done():
			c.JSON(http.StatusRequestTimeout, utilitas.FormatRespons(false, "Permintaan melewati batas waktu (timeout)", nil))
			c.Abort()
		}
	}
}
