// Package infrastruktur memuat penjaga sesi dan otentikasi lapis middleware.
package infrastruktur

import (
	"os"
	"strings"

	"pgn-server/pkg/pencatatan_log"
	"pgn-server/pkg/respon"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// PenjagaSesiJWT merupakan middleware (penjebak perlintasan) otorisasi eksklusif.
func PenjagaSesiJWT() gin.HandlerFunc {
	return func(k *gin.Context) {
		tajukOtorisasi := k.GetHeader("Authorization")
		if tajukOtorisasi == "" {
			respon.Galat_TidakSah(k, "Akses ditolak: Tajuk otorisasi tidak ditemukan")
			k.Abort()
			return
		}

		// Ekspektasi: "Bearer <token>"
		bagianTajuk := strings.Split(tajukOtorisasi, " ")
		if len(bagianTajuk) != 2 || strings.ToLower(bagianTajuk[0]) != "bearer" {
			respon.Galat_TidakSah(k, "Akses ditolak: Format tajuk otorisasi tidak valid")
			k.Abort()
			return
		}

		tokenString := bagianTajuk[1]
		kunciRahasia := os.Getenv("JWT_SECRET")
		if kunciRahasia == "" {
			kunciRahasia = "rahasia-default"
		}

		token, errVerifikasi := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validasi algoritma sandi
			if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(kunciRahasia), nil
		})

		if errVerifikasi != nil || !token.Valid {
			respon.Galat_TidakSah(k, "Akses ditolak: Sesi tidak sah atau telah kedaluwarsa")
			k.Abort()
			return
		}

		// Menyisipkan data klaim pengguna ke dalam konteks permintaan
		if klaim, terkonfirmasi := token.Claims.(jwt.MapClaims); terkonfirmasi {
			k.Set("id_pengguna", klaim["id"])
			k.Set("surel", klaim["surel"])
			k.Set("peran", klaim["peran"])
		}

		k.Next()
	}
}

// PenjagaRole bertindak sebagai authorization gatekeeper berbasis RBAC
func PenjagaRole(rolesLolos ...string) gin.HandlerFunc {
	return func(k *gin.Context) {
		peranAktif, ada := k.Get("peran")
		if !ada {
			respon.Galat_TidakSah(k, "Unauthorized: Sesi tidak memiliki konteks role (Missing Context)")
			k.Abort()
			return
		}

		peranString, ok := peranAktif.(string)
		if !ok {
			respon.Galat_TidakSah(k, "Unauthorized: Role korup atau tidak terbaca")
			k.Abort()
			return
		}

		roleValid := false
		for _, r := range rolesLolos {
			if peranString == r {
				roleValid = true
				break
			}
		}

		if !roleValid {
			respon.Galat_Dilarang(k, "Forbidden: Role kamu saat ini tidak memiliki privilege untuk mengakses resource ini")
			k.Abort()
			return
		}

		k.Next()
	}
}

// MiddlewareCorrelationID menginjeksi X-Request-ID (Correlation ID) di header request dan response
func MiddlewareCorrelationID() gin.HandlerFunc {
	return func(k *gin.Context) {
		requestID := k.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = pencatatan_log.HasilkanUUIDv4()
		}

		// Set di konteks Gin agar bisa diakses oleh logger/handler
		k.Set("RequestID", requestID)

		// Set di header respons agar klien bisa memverifikasi
		k.Writer.Header().Set("X-Request-ID", requestID)

		// Log kedatangan request
		pencatatan_log.Info(k, "Permintaan masuk: %s %s dari %s", k.Request.Method, k.Request.URL.Path, k.ClientIP())

		k.Next()

		// Log keluar request
		pencatatan_log.Info(k, "Permintaan selesai: %s %s -> Status %d", k.Request.Method, k.Request.URL.Path, k.Writer.Status())
	}
}
