package respon

import (
	"net/http"

	"pgn-server/pkg/pencatatan_log"

	"github.com/gin-gonic/gin"
)

// DetailKesalahan merepresentasikan rincian galat validasi field
type DetailKesalahan struct {
	Field string `json:"field"`
	Pesan string `json:"pesan"`
}

// ResponStandar adalah struktur data untuk standarisasi balasan API (Kompatibel penuh dengan QControl)
type ResponStandar struct {
	Sukses    bool              `json:"sukses"`
	Pesan     string            `json:"pesan"`
	Data      interface{}       `json:"data"`
	Metadata  interface{}       `json:"metadata"`
	Kesalahan []DetailKesalahan `json:"kesalahan"`

	// Legacy fields untuk retro-kompatibilitas
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Meta    interface{} `json:"meta"`
	Errors  []string    `json:"errors"`
}

// Sukses mengembalikan format respons sukses standar.
func Sukses(k *gin.Context, pesan string, data interface{}) {
	k.JSON(http.StatusOK, gin.H{
		"sukses":    true,
		"pesan":     pesan,
		"data":      data,
		"metadata":  nil,
		"kesalahan": nil,
		// Legacy fields
		"success": true,
		"message": pesan,
		"status":  "success",
		"meta":    nil,
		"errors":  nil,
	})
}

// SuksesDenganMetadata mengembalikan format respons sukses standar beserta informasi metadata.
func SuksesDenganMetadata(k *gin.Context, pesan string, data interface{}, metadata interface{}) {
	k.JSON(http.StatusOK, gin.H{
		"sukses":    true,
		"pesan":     pesan,
		"data":      data,
		"metadata":  metadata,
		"kesalahan": nil,
		// Legacy fields
		"success": true,
		"message": pesan,
		"status":  "success",
		"meta":    metadata,
		"errors":  nil,
	})
}

// Galat_Server mengembalikan format respons untuk kesalahan internal server (ramah-antarmuka).
// Stack trace error sesungguhnya HANYA di-log di sisi backend (tidak bocor ke klien).
func Galat_Server(k *gin.Context, pesan string, err error) {
	if err != nil {
		pencatatan_log.Kritis(k, "500 Server Error: %v | Konteks: %s", err, pesan)
	}
	k.JSON(http.StatusInternalServerError, gin.H{
		"sukses":    false,
		"pesan":     "Terjadi kendala internal pada sistem kami. Silakan hubungi dukungan teknis.",
		"data":      nil,
		"metadata":  nil,
		"kesalahan": nil,
		// Legacy fields
		"success": false,
		"message": "Terjadi kendala internal pada sistem kami. Silakan hubungi dukungan teknis.",
		"status":  "error",
		"meta":    nil,
		"errors":  []string{},
	})
}

// Galat_Validasi mengembalikan format respons untuk data masukan yang tidak valid (400 Bad Request).
func Galat_Validasi(k *gin.Context, pesan string, rincianGalat []string) {
	var listKesalahan []DetailKesalahan
	if len(rincianGalat) > 0 {
		for _, g := range rincianGalat {
			// Parsing sederhana jika string berupa "field: pesan"
			found := false
			for i := 0; i < len(g); i++ {
				if g[i] == ':' && i > 0 && i < len(g)-1 {
					field := g[:i]
					msg := g[i+1:]
					// trim spaces
					for len(field) > 0 && field[0] == ' ' {
						field = field[1:]
					}
					for len(field) > 0 && field[len(field)-1] == ' ' {
						field = field[:len(field)-1]
					}
					for len(msg) > 0 && msg[0] == ' ' {
						msg = msg[1:]
					}
					for len(msg) > 0 && msg[len(msg)-1] == ' ' {
						msg = msg[:len(msg)-1]
					}
					listKesalahan = append(listKesalahan, DetailKesalahan{
						Field: field,
						Pesan: msg,
					})
					found = true
					break
				}
			}
			if !found {
				listKesalahan = append(listKesalahan, DetailKesalahan{
					Field: "global",
					Pesan: g,
				})
			}
		}
	} else {
		listKesalahan = append(listKesalahan, DetailKesalahan{
			Field: "global",
			Pesan: pesan,
		})
	}

	if rincianGalat == nil {
		rincianGalat = []string{}
	}

	k.JSON(http.StatusBadRequest, gin.H{
		"sukses":    false,
		"pesan":     pesan,
		"data":      nil,
		"metadata":  nil,
		"kesalahan": listKesalahan,
		// Legacy fields
		"success": false,
		"message": pesan,
		"status":  "error",
		"meta":    nil,
		"errors":  rincianGalat,
	})
}

// Galat_TidakSah mengembalikan format respons Unauthorized (401).
func Galat_TidakSah(k *gin.Context, pesan string) {
	k.JSON(http.StatusUnauthorized, gin.H{
		"sukses":    false,
		"pesan":     pesan,
		"data":      nil,
		"metadata":  nil,
		"kesalahan": []DetailKesalahan{{Field: "global", Pesan: pesan}},
		// Legacy fields
		"success": false,
		"message": pesan,
		"status":  "error",
		"meta":    nil,
		"errors":  []string{pesan},
	})
}

// Galat_Dilarang mengembalikan format respons Forbidden (403).
func Galat_Dilarang(k *gin.Context, pesan string) {
	k.JSON(http.StatusForbidden, gin.H{
		"sukses":    false,
		"pesan":     pesan,
		"data":      nil,
		"metadata":  nil,
		"kesalahan": []DetailKesalahan{{Field: "global", Pesan: pesan}},
		// Legacy fields
		"success": false,
		"message": pesan,
		"status":  "error",
		"meta":    nil,
		"errors":  []string{pesan},
	})
}

// Galat_TidakDitemukan mengembalikan format respons Not Found (404).
func Galat_TidakDitemukan(k *gin.Context, pesan string) {
	k.JSON(http.StatusNotFound, gin.H{
		"sukses":    false,
		"pesan":     pesan,
		"data":      nil,
		"metadata":  nil,
		"kesalahan": []DetailKesalahan{{Field: "global", Pesan: pesan}},
		// Legacy fields
		"success": false,
		"message": pesan,
		"status":  "error",
		"meta":    nil,
		"errors":  []string{pesan},
	})
}

// Galat_TerlaluBanyakPermintaan mengembalikan format respons Rate Limit Exceeded (429).
func Galat_TerlaluBanyakPermintaan(k *gin.Context, pesan string) {
	k.JSON(http.StatusTooManyRequests, gin.H{
		"sukses":    false,
		"pesan":     pesan,
		"data":      nil,
		"metadata":  nil,
		"kesalahan": []DetailKesalahan{{Field: "global", Pesan: pesan}},
		// Legacy fields
		"success": false,
		"message": pesan,
		"status":  "error",
		"meta":    nil,
		"errors":  []string{pesan},
	})
}
