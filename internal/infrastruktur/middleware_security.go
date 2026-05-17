package infrastruktur

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"pgn-server/pkg/respon"
)

// rateLimiterClient membungkus *rate.Limiter dan penanda akses terakhir
type rateLimiterClient struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*rateLimiterClient)
)

// MulaiPembersihRateLimiter berjalan di background untuk membersihkan memori klien yang sudah tidak aktif
func init() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

// MiddlewareRateLimiter menggunakan Token Bucket Algorithm untuk membatasi request per IP
func MiddlewareRateLimiter(batasRequestPerDetik float64, ukuranBucket int) gin.HandlerFunc {
	return func(k *gin.Context) {
		ipClient := k.ClientIP()

		mu.Lock()
		if _, ditemukan := clients[ipClient]; !ditemukan {
			clients[ipClient] = &rateLimiterClient{
				limiter: rate.NewLimiter(rate.Limit(batasRequestPerDetik), ukuranBucket),
			}
		}
		clients[ipClient].lastSeen = time.Now()
		limiterClient := clients[ipClient].limiter
		mu.Unlock()

		if !limiterClient.Allow() {
			respon.Galat_TerlaluBanyakPermintaan(k, "Woops, kamu request terlalu cepat. Coba lagi beberapa saat ya (Rate Limit Exceeded).")
			k.Abort()
			return
		}

		k.Next()
	}
}

// MiddlewarePenangkapPanic mencegah crash aplikasi dan mengembalikan standar JSON 500
func MiddlewarePenangkapPanic() gin.HandlerFunc {
	return func(k *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				respon.Galat_Server(k, "Panic terdeteksi di Middleware", fmt.Errorf("%v", err))
				k.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		k.Next()
	}
}
