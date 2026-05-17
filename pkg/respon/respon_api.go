package respon

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponStandar adalah struktur data untuk standarisasi balasan API.
type ResponStandar struct {
	Sukses bool        `json:"sukses"`
	Pesan  string      `json:"pesan"`
	Data   interface{} `json:"data,omitempty"`
}

// Sukses mengembalikan format respons sukses standar.
func Sukses(k *gin.Context, pesan string, data interface{}) {
	k.JSON(http.StatusOK, ResponStandar{
		Sukses: true,
		Pesan:  pesan,
		Data:   data,
	})
}

// Galat_Server mengembalikan format respons untuk kesalahan internal server (ramah-antarmuka).
func Galat_Server(k *gin.Context, pesan string) {
	k.JSON(http.StatusInternalServerError, ResponStandar{
		Sukses: false,
		Pesan:  "Terjadi kendala pada sistem kami. Silakan coba beberapa saat lagi: " + pesan,
	})
}

// Galat_Validasi mengembalikan format respons untuk data masukan yang tidak valid.
func Galat_Validasi(k *gin.Context, pesan string) {
	k.JSON(http.StatusBadRequest, ResponStandar{
		Sukses: false,
		Pesan:  "Data yang Anda masukkan tidak sesuai kriteria: " + pesan,
	})
}
