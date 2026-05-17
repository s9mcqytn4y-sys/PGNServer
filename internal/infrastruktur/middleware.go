// Package infrastruktur memuat penjaga sesi dan otentikasi lapis middleware.
package infrastruktur

import (
	"os"
	"strings"

	"pgn-server/pkg/respon"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// PenjagaSesiJWT merupakan middleware (penjebak perlintasan) otorisasi eksklusif.
func PenjagaSesiJWT() gin.HandlerFunc {
	return func(k *gin.Context) {
		tajukOtorisasi := k.GetHeader("Authorization")
		if tajukOtorisasi == "" {
			respon.Galat_Validasi(k, "Akses ditolak: Tajuk otorisasi tidak ditemukan")
			k.Abort()
			return
		}

		// Ekspektasi: "Bearer <token>"
		bagianTajuk := strings.Split(tajukOtorisasi, " ")
		if len(bagianTajuk) != 2 || strings.ToLower(bagianTajuk[0]) != "bearer" {
			respon.Galat_Validasi(k, "Akses ditolak: Format tajuk otorisasi tidak valid")
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
			respon.Galat_Validasi(k, "Akses ditolak: Sesi tidak sah atau telah kedaluwarsa")
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
