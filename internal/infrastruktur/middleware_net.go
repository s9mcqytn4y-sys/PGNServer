package infrastruktur

import (
	"net"
	"net/http"
	"os"
	"strings"

	"pgn-server/pkg/pencatatan_log"
	"pgn-server/pkg/respon"

	"github.com/gin-gonic/gin"
)

// MiddlewareCORS mengoptimalkan Cross-Origin Resource Sharing secara dinamis sesuai spesifikasi
func MiddlewareCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		asalPermintaan := c.GetHeader("Origin")
		if asalPermintaan == "" {
			// Jika request tidak memiliki header Origin, maka bukan CORS request
			c.Next()
			return
		}

		// Memuat daftar origin yang diizinkan dari environment variable
		originsEnv := os.Getenv("ALLOWED_ORIGINS")
		originsDiizinkan := []string{}
		if originsEnv != "" {
			for _, o := range strings.Split(originsEnv, ",") {
				originsDiizinkan = append(originsDiizinkan, strings.TrimSpace(o))
			}
		}

		asalDiterima := ""
		for _, o := range originsDiizinkan {
			if o == "*" || o == asalPermintaan {
				asalDiterima = asalPermintaan
				break
			}
		}

		// Development fallback jika ALLOWED_ORIGINS kosong demi mempermudah integrasi tim frontend
		if asalDiterima == "" && len(originsDiizinkan) == 0 {
			asalDiterima = asalPermintaan
		}

		if asalDiterima == "" {
			respon.Galat_Dilarang(c, "Akses ditolak oleh kebijakan CORS")
			c.Abort()
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", asalDiterima)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Methods & Headers Hardening
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Max-Age", "43200") // 12 Jam cache preflight

		// Preflight Optimization
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ekstrakRealIP mengambil IP asli client secara proxy-aware dan aman dari spoofing
func ekstrakRealIP(c *gin.Context) string {
	// 1. Cek X-Forwarded-For (leftmost IP)
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		bagian := strings.Split(xff, ",")
		ipKiri := strings.TrimSpace(bagian[0])
		if ipKiri != "" {
			// Hapus port jika ada
			if host, _, err := net.SplitHostPort(ipKiri); err == nil {
				return host
			}
			return ipKiri
		}
	}

	// 2. Cek X-Real-IP
	xri := c.GetHeader("X-Real-IP")
	if xri != "" {
		ipReal := strings.TrimSpace(xri)
		if host, _, err := net.SplitHostPort(ipReal); err == nil {
			return host
		}
		return ipReal
	}

	// 3. Fallback ke RemoteAddr
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err == nil {
		return host
	}
	return c.Request.RemoteAddr
}

// MiddlewareIPWhitelist membatasi akses rute khusus hanya untuk daftar IP tertentu
func MiddlewareIPWhitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		ipWhitelistsEnv := os.Getenv("IP_WHITELIST")
		
		// Jika konfigurasi kosong, kita tolak akses secara default demi keamanan enterprise
		if ipWhitelistsEnv == "" {
			pencatatan_log.Peringatan(c, "IP_WHITELIST kosong di .env. Akses ditolak secara default.")
			respon.Galat_Dilarang(c, "Akses ditolak, IP kamu tidak terdaftar dalam sistem whitelist kami")
			c.Abort()
			return
		}

		daftarIP := []string{}
		for _, ipRaw := range strings.Split(ipWhitelistsEnv, ",") {
			daftarIP = append(daftarIP, strings.TrimSpace(ipRaw))
		}

		ipKlien := ekstrakRealIP(c)
		
		// Validasi pencocokan IP klien dengan whitelist
		ipDiizinkan := false
		for _, ipWhitelisted := range daftarIP {
			if ipKlien == ipWhitelisted {
				ipDiizinkan = true
				break
			}
		}

		if !ipDiizinkan {
			pencatatan_log.Peringatan(c, "Akses ditolak untuk IP: %s (tidak terdaftar di whitelist: %s)", ipKlien, ipWhitelistsEnv)
			respon.Galat_Dilarang(c, "Akses ditolak, IP kamu tidak terdaftar dalam sistem whitelist kami")
			c.Abort()
			return
		}

		pencatatan_log.Info(c, "Akses diizinkan untuk IP whitelisted: %s", ipKlien)
		c.Next()
	}
}

// MiddlewareSecureHeaders menyuntikkan header keamanan standar korporat secara global
func MiddlewareSecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		
		c.Next()
	}
}
