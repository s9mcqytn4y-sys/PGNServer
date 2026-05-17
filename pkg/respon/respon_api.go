package respon

import (
	"net/http"

	"pgn-server/pkg/pencatatan_log"

	"github.com/gin-gonic/gin"
)

// ResponStandar adalah struktur data untuk standarisasi balasan API
type ResponStandar struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
	Errors  []string    `json:"errors"`
}

// Sukses mengembalikan format respons sukses standar.
func Sukses(k *gin.Context, pesan string, data interface{}) {
	k.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": pesan,
		"data":    data,
		"meta":    nil,
	})
}

// Galat_Server mengembalikan format respons untuk kesalahan internal server (ramah-antarmuka).
// Stack trace error sesungguhnya HANYA di-log di sisi backend (tidak bocor ke klien).
func Galat_Server(k *gin.Context, pesan string, err error) {
	if err != nil {
		pencatatan_log.Kritis(k, "500 Server Error: %v | Konteks: %s", err, pesan)
	}
	k.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": "Terjadi kendala internal pada sistem kami. Silakan hubungi dukungan teknis.",
		"errors":  []string{},
	})
}

// Galat_Validasi mengembalikan format respons untuk data masukan yang tidak valid (400 Bad Request).
func Galat_Validasi(k *gin.Context, pesan string, rincianGalat []string) {
	if rincianGalat == nil {
		rincianGalat = []string{}
	}
	k.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": pesan,
		"errors":  rincianGalat,
	})
}

// Galat_TidakSah mengembalikan format respons Unauthorized (401).
func Galat_TidakSah(k *gin.Context, pesan string) {
	k.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": pesan,
		"errors":  []string{},
	})
}

// Galat_Dilarang mengembalikan format respons Forbidden (403).
func Galat_Dilarang(k *gin.Context, pesan string) {
	k.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"message": pesan,
		"errors":  []string{},
	})
}

// Galat_TidakDitemukan mengembalikan format respons Not Found (404).
func Galat_TidakDitemukan(k *gin.Context, pesan string) {
	k.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": pesan,
		"errors":  []string{},
	})
}

// Galat_TerlaluBanyakPermintaan mengembalikan format respons Rate Limit Exceeded (429).
func Galat_TerlaluBanyakPermintaan(k *gin.Context, pesan string) {
	k.JSON(http.StatusTooManyRequests, gin.H{
		"success": false,
		"message": pesan,
		"errors":  []string{},
	})
}
