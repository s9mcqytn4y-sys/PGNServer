package respon

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponStandar adalah struktur data untuk standarisasi balasan API (Tipe Sukses)
type ResponStandar struct {
	Status  string      `json:"status"` // "success" atau "fail" atau "error"
	Pesan   string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Galat   []string    `json:"errors,omitempty"`
}

// Sukses mengembalikan format respons sukses standar.
func Sukses(k *gin.Context, pesan string, data interface{}) {
	k.JSON(http.StatusOK, ResponStandar{
		Status: "success",
		Pesan:  pesan,
		Data:   data,
	})
}

// Galat_Server mengembalikan format respons untuk kesalahan internal server (ramah-antarmuka).
// Stack trace error sesungguhnya HANYA di-log di sisi backend (tidak bocor ke klien).
func Galat_Server(k *gin.Context, pesan string, err error) {
	if err != nil {
		log.Printf("[CRITICAL ERROR] 500 Server Error: %v | Konteks: %s\n", err, pesan)
	}
	k.JSON(http.StatusInternalServerError, ResponStandar{
		Status: "error",
		Pesan:  "Terjadi kendala internal pada sistem kami. Silakan hubungi dukungan teknis.",
	})
}

// Galat_Validasi mengembalikan format respons untuk data masukan yang tidak valid (400 Bad Request).
func Galat_Validasi(k *gin.Context, pesan string, rincianGalat []string) {
	k.JSON(http.StatusBadRequest, ResponStandar{
		Status: "fail",
		Pesan:  pesan,
		Galat:  rincianGalat,
	})
}

// Galat_TidakSah mengembalikan format respons Unauthorized (401).
func Galat_TidakSah(k *gin.Context, pesan string) {
	k.JSON(http.StatusUnauthorized, ResponStandar{
		Status: "fail",
		Pesan:  pesan,
	})
}

// Galat_Dilarang mengembalikan format respons Forbidden (403).
func Galat_Dilarang(k *gin.Context, pesan string) {
	k.JSON(http.StatusForbidden, ResponStandar{
		Status: "fail",
		Pesan:  pesan,
	})
}

// Galat_TidakDitemukan mengembalikan format respons Not Found (404).
func Galat_TidakDitemukan(k *gin.Context, pesan string) {
	k.JSON(http.StatusNotFound, ResponStandar{
		Status: "fail",
		Pesan:  pesan,
	})
}

// Galat_TerlaluBanyakPermintaan mengembalikan format respons Rate Limit Exceeded (429).
func Galat_TerlaluBanyakPermintaan(k *gin.Context, pesan string) {
	k.JSON(http.StatusTooManyRequests, ResponStandar{
		Status: "fail",
		Pesan:  pesan,
	})
}
