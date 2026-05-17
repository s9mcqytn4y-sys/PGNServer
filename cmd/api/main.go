package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "pgn-server/docs"
	"pgn-server/internal/analitik"
	"pgn-server/internal/infrastruktur"
	"pgn-server/internal/kualitas"
	"pgn-server/internal/otentikasi"
	"pgn-server/pkg/respon"
)

// @title PGNServer API
// @version 1.0
// @description REST API untuk ekosistem manufaktur dan kontrol kualitas PGNServer.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Inisialisasi pengenalan cgroup limit di Go 1.25.x
	// Meskipun Go 1.25.x memiliki peningkatan otomatisasi, menetapkan GOMAXPROCS
	// sesuai jumlah CPU sistem yang dialokasikan di dalam kontainer masih merupakan praktik yang baik.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Memuat konfigurasi
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Berkas .env tidak ditemukan, menggunakan variabel lingkungan sistem.")
	}

	// Mempersiapkan string koneksi ke PostgreSQL
	inangDB := os.Getenv("DB_HOST")
	if inangDB == "" {
		inangDB = "localhost"
	}
	penggunaDB := os.Getenv("DB_USER")
	sandiDB := os.Getenv("DB_PASSWORD")
	namaDB := os.Getenv("DB_NAME")
	pelabuhanDB := os.Getenv("DB_PORT")
	if pelabuhanDB == "" {
		pelabuhanDB = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		inangDB, penggunaDB, sandiDB, namaDB, pelabuhanDB)

	// Inisialisasi GORM ke PostgreSQL
	db, errDB := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDB != nil {
		log.Fatalf("Gagal terhubung ke pangkalan data: %v", errDB)
	}

	// Mendapatkan objek underlying *sql.DB untuk penyetelan lebih lanjut
	sqlDB, errSQL := db.DB()
	if errSQL == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
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

	// Setup Dependensi Kualitas
	repoKualitas := kualitas.KonstruksiRepositoriBaru()
	layananKualitas := kualitas.KonstruksiLayananBaru(repoKualitas, db)
	handlerKualitas := kualitas.KonstruksiPenangananBaru(layananKualitas)

	// Setup Dependensi Analitik
	repoAnalitik := analitik.KonstruksiRepositoriBaru(db)
	layananAnalitik := analitik.KonstruksiLayananBaru(repoAnalitik)
	handlerAnalitik := analitik.KonstruksiPenangananBaru(layananAnalitik)

	// Konfigurasi layanan router web Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	rute := gin.Default()
	rute.SetTrustedProxies(nil) // Mengamankan peringatan 'trusted all proxies'

	// Rute Publik untuk Swagger
	rute.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Kumpulan Endpoint API
	api := rute.Group("/api/v1")
	{
		// Endpoint pemeriksaan sistem (Health Check)
		api.GET("/cek_sistem", func(k *gin.Context) {
			errPing := sqlDB.Ping()
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

		// Endpoint Operasi Internal (Perlu JWT)
		operasi := api.Group("/operasi")
		operasi.Use(infrastruktur.PenjagaSesiJWT())
		{
			operasi.POST("/rekam_lembar_periksa", handlerKualitas.TanganiRekamLembarPeriksa)
		}

		// Endpoint Analitik (Terbuka/JWT)
		analitikGrup := api.Group("/analitik")
		// analitikGrup.Use(infrastruktur.PenjagaSesiJWT()) // Aktifkan jika analitik butuh JWT
		{
			analitikGrup.GET("/metrik_pareto_bulanan", handlerAnalitik.TanganiParetoBulanan)
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
