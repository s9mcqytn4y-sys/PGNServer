package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"pgn-server/internal/infrastruktur"
	"pgn-server/internal/otentikasi"
	"pgn-server/pkg/respon"
)

func main() {
	// Inisialisasi pengenalan cgroup limit di Go 1.25.x
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Memuat konfigurasi
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Berkas .env tidak ditemukan, menggunakan variabel lingkungan sistem.")
	}

	// Mempersiapkan string koneksi ke PostgreSQL
	inangDb := os.Getenv("DB_HOST")
	penggunaDb := os.Getenv("DB_USER")
	sandiDb := os.Getenv("DB_PASSWORD")
	namaDb := os.Getenv("DB_NAME")
	pelabuhanDb := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		inangDb, penggunaDb, sandiDb, namaDb, pelabuhanDb)

	// Inisialisasi GORM ke PostgreSQL
	db, errDb := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDb != nil {
		log.Fatalf("Gagal terhubung ke pangkalan data: %v", errDb)
	}

	// Mendapatkan objek underlying *sql.DB untuk penyetelan lebih lanjut
	sqlDb, errSql := db.DB()
	if errSql == nil {
		sqlDb.SetMaxIdleConns(10)
		sqlDb.SetMaxOpenConns(100)
	}

	// Menjalankan automigrasi entitas manufaktur dan otentikasi
	errMigrasi := infrastruktur.PelaksanaanAutoMigrasi(db)
	if errMigrasi != nil {
		log.Fatalf("Gagal melaksanakan automigrasi: %v", errMigrasi)
	}

	// Setup Dependensi Otentikasi
	repoOtentikasi := otentikasi.EkstraksiRepositoriBaru(db)
	layananOtentikasi := otentikasi.KonstruksiLayananBaru(repoOtentikasi)
	handlerOtentikasi := otentikasi.KonstruksiPenangananBaru(layananOtentikasi)

	// Konfigurasi layanan router web Gin
	rute := gin.Default()

	// Kumpulan Endpoint API
	api := rute.Group("/api/v1")
	{
		// Endpoint pemeriksaan sistem (Health Check)
		api.GET("/cek_sistem", func(k *gin.Context) {
			errPing := sqlDb.Ping()
			if errPing != nil {
				respon.Galat_Server(k, "Pangkalan data tidak dapat dijangkau.")
				return
			}
			respon.Sukses(k, "Sistem PGNServer beroperasi secara optimal dan terhubung ke pangkalan data.", nil)
		})

		// Endpoint Autentikasi Publik
		auth := api.Group("/otentikasi")
		{
			auth.POST("/daftar", handlerOtentikasi.TanganiRegistrasi)
			auth.POST("/masuk", handlerOtentikasi.TanganiLogin)
			auth.POST("/lupa-sandi", handlerOtentikasi.TanganiLupaSandi)
		}
	}

	// Menjalankan server
	pelabuhanAplikasi := os.Getenv("APP_PORT")
	if pelabuhanAplikasi == "" {
		pelabuhanAplikasi = "8080"
	}

	log.Printf("Memulai layanan di pelabuhan %s...", pelabuhanAplikasi)
	if errJalan := rute.Run(":" + pelabuhanAplikasi); errJalan != nil {
		log.Fatalf("Gagal menjalankan server: %v", errJalan)
	}
}
