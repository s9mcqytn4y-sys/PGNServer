package pengelola

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"pgn-server/pkg/utilitas"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

// AutentikasiMiddleware memverifikasi Bearer JWT dari Header
func AutentikasiMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, utilitas.FormatRespons(false, "Tidak terotorisasi, token tidak ditemukan", nil))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "pgn-secret-key-2026"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak terduga: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, utilitas.FormatRespons(false, "Token tidak valid atau kedaluwarsa", nil))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, utilitas.FormatRespons(false, "Gagal mengambil claims", nil))
			c.Abort()
			return
		}

		// Simpan info user ke context
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

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
