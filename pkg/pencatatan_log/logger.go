package pencatatan_log

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Level mendefinisikan tingkatan log sistem
type Level string

const (
	INFO     Level = "INFO"
	WARN     Level = "WARN"
	ERROR    Level = "ERROR"
	CRITICAL Level = "CRITICAL"
)

// HasilkanUUIDv4 membuat UUID unik versi 4 secara kriptografis tanpa dependensi eksternal
func HasilkanUUIDv4() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback sederhana jika terjadi kegagalan pembacaan entropi
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	
	// Sesuai dengan spesifikasi RFC 4122
	b[6] = (b[6] & 0x0f) | 0x40 // Set versi ke 4
	b[8] = (b[8] & 0x3f) | 0x80 // Set varian ke 10xx

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// Catat mencetak pesan log terformat lengkap dengan Correlation ID (Request ID)
func Catat(k *gin.Context, level Level, format string, v ...interface{}) {
	reqID := "-"
	if k != nil {
		if id, ada := k.Get("RequestID"); ada {
			if strID, ok := id.(string); ok {
				reqID = strID
			}
		} else {
			// Coba ambil dari header respons jika context belum memuat
			headerID := k.Writer.Header().Get("X-Request-ID")
			if headerID != "" {
				reqID = headerID
			}
		}
	}

	waktu := time.Now().Format("2006-05-02 15:04:05.000")
	pesan := fmt.Sprintf(format, v...)
	
	// Format: [WAKTU] [LEVEL] [Correlation-ID: X-Request-ID] Pesan
	log.Printf("[%s] [%s] [Correlation-ID: %s] %s\n", waktu, level, reqID, pesan)
}

// Info mencetak log informatif umum
func Info(k *gin.Context, format string, v ...interface{}) {
	Catat(k, INFO, format, v...)
}

// Peringatan mencetak log peringatan non-fatal
func Peringatan(k *gin.Context, format string, v ...interface{}) {
	Catat(k, WARN, format, v...)
}

// Galat mencetak log kegagalan operasional sistem
func Galat(k *gin.Context, format string, v ...interface{}) {
	Catat(k, ERROR, format, v...)
}

// Kritis mencetak log kegagalan fatal yang butuh penanganan segera
func Kritis(k *gin.Context, format string, v ...interface{}) {
	Catat(k, CRITICAL, format, v...)
}
