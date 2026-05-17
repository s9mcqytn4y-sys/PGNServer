package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "go.uber.org/automaxprocs"

	_ "pgn-server/docs"
	"pgn-server/internal/analitik"
	"pgn-server/internal/infrastruktur"
	"pgn-server/internal/kualitas"
	"pgn-server/internal/manufaktur"
	"pgn-server/internal/media"
	"pgn-server/internal/otentikasi"
	"pgn-server/pkg/respon"
)

//go:embed landing.html
var landingHTML string

// @title PGNServer API
// @version 1.0
// @description REST API untuk ekosistem manufaktur dan kontrol kualitas PGNServer.
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Inisialisasi pengenalan cgroup limit di Go 1.25.x
	// Automaxprocs secara implisit menyesuaikan GOMAXPROCS tanpa konfigurasi manual

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
		sqlDB.SetConnMaxLifetime(5 * time.Minute) // Tambahan batas umur koneksi maksimal
	}

	// Menjalankan automigrasi entitas manufaktur dan otentikasi
	errMigrasi := infrastruktur.PelaksanaanAutoMigrasi(db)
	if errMigrasi != nil {
		log.Fatalf("Gagal melaksanakan automigrasi: %v", errMigrasi)
	}

	// Menjalankan Seeder Manufaktur
	errSeeder := manufaktur.JalankanSeeder(db)
	if errSeeder != nil {
		log.Printf("Gagal menjalankan seeder: %v", errSeeder)
	}

	// Setup Dependensi Otentikasi
	repoOtentikasi := otentikasi.EkstraksiRepositoriBaru(db)
	layananOtentikasi := otentikasi.KonstruksiLayananBaru(repoOtentikasi)
	handlerOtentikasi := otentikasi.KonstruksiPenangananBaru(layananOtentikasi)

	// Setup Dependensi Kualitas
	repoKualitas := kualitas.KonstruksiRepositoriBaru(db)
	layananKualitas := kualitas.KonstruksiLayananBaru(repoKualitas, db)
	handlerKualitas := kualitas.KonstruksiPenangananBaru(layananKualitas)

	// Setup Dependensi Analitik
	repoAnalitik := analitik.KonstruksiRepositoriBaru(db)
	layananAnalitik := analitik.KonstruksiLayananBaru(repoAnalitik, db)
	handlerAnalitik := analitik.KonstruksiPenangananBaru(layananAnalitik)

	// Setup Dependensi Media
	repoMedia := media.KonstruksiRepositoriBaru(db)
	layananMedia := media.KonstruksiLayananBaru(repoMedia, db)
	handlerMedia := media.KonstruksiPenangananBaru(layananMedia)

	// Setup Dependensi Manufaktur
	repoManufaktur := manufaktur.KonstruksiRepositoriBaru(db)
	layananManufaktur := manufaktur.KonstruksiLayananBaru(repoManufaktur)
	handlerManufaktur := manufaktur.KonstruksiPenangananBaru(layananManufaktur)

	// Konfigurasi layanan router web Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	rute := gin.Default()
	rute.SetTrustedProxies(nil) // Mengamankan peringatan 'trusted all proxies'

	// Middleware Keamanan Global
	rute.Use(infrastruktur.MiddlewareCORS())
	rute.Use(infrastruktur.MiddlewareSecureHeaders())
	rute.Use(infrastruktur.MiddlewareCorrelationID())
	rute.Use(infrastruktur.MiddlewarePenangkapPanic())
	rute.Use(infrastruktur.MiddlewareRateLimiter(10, 20)) // 10 rps, 20 kapasitas

	// Rute Publik untuk Swagger dan Landing Page
	rute.GET("/", tanganiLandingPage(db))
	rute.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Kumpulan Endpoint API
	api := rute.Group("/api/v1")
	{
		// Landing page pada root API
		api.GET("/", tanganiLandingPage(db))

		// Endpoint pemeriksaan sistem (Health Check & Readiness)
		api.GET("/health", TanganiHealth(db))
		api.GET("/kesehatan", TanganiKesehatan(db))
		api.GET("/readiness", TanganiReadiness(db))
		api.GET("/cek_sistem", TanganiCekSistem(db))

		// Endpoint Operasi Krusial (Perlu Whitelist IP)
		krusial := api.Group("/krusial")
		krusial.Use(infrastruktur.MiddlewareIPWhitelist())
		{
			krusial.GET("/status", TanganiStatusKrusial())
		}

		// Endpoint Autentikasi Publik
		auth := api.Group("/otentikasi")
		{
			auth.POST("/daftar", handlerOtentikasi.TanganiRegistrasi)
			auth.POST("/masuk", handlerOtentikasi.TanganiLogin)
			auth.POST("/lupa-sandi", handlerOtentikasi.TanganiLupaSandi)
			auth.POST("/keluar", handlerOtentikasi.TanganiLogout)
			auth.GET("/profil", infrastruktur.PenjagaSesiJWT(), handlerOtentikasi.TanganiProfil)
		}

		// Endpoint Operasi Internal (Perlu JWT)
		operasi := api.Group("/operasi")
		operasi.Use(infrastruktur.PenjagaSesiJWT())
		{
			operasi.POST("/rekam_lembar_periksa", handlerKualitas.TanganiRekamLembarPeriksa)
			operasi.GET("/riwayat_lembar_periksa", handlerKualitas.TanganiDaftarRiwayat)
			operasi.GET("/lembar_periksa/options", handlerKualitas.TanganiOpsiLembarPeriksa)
		}

		// Endpoint Analitik (Terbuka/JWT)
		analitikGrup := api.Group("/analitik")
		// analitikGrup.Use(infrastruktur.PenjagaSesiJWT()) // Aktifkan jika analitik butuh JWT
		{
			analitikGrup.GET("/metrik_pareto_bulanan", handlerAnalitik.TanganiParetoBulanan)
			analitikGrup.GET("/pareto", handlerAnalitik.TanganiPareto)
			analitikGrup.GET("/lacak", handlerAnalitik.TanganiLacakAkarMasalah)
			analitikGrup.GET("/ringkasan_ng", handlerAnalitik.TanganiRingkasanNG)
			analitikGrup.GET("/histogram_defect", handlerAnalitik.TanganiHistogramDefect)
			analitikGrup.GET("/trend_defect", handlerAnalitik.TanganiTrendDefect)
			analitikGrup.GET("/stratifikasi_defect", handlerAnalitik.TanganiStratifikasiDefect)
			analitikGrup.GET("/sinyal_kualitas", handlerAnalitik.TanganiSinyalKualitas)
		}

		// Endpoint Media
		mediaGrup := api.Group("/media")
		{
			mediaGrup.GET("/:id/pratinjau", handlerMedia.TanganiPratinjauMedia)
		}

		// Endpoint Manufaktur (Perlu JWT)
		pemasokGrup := api.Group("/suppliers")
		pemasokGrup.Use(infrastruktur.PenjagaSesiJWT())
		{
			pemasokGrup.POST("", handlerManufaktur.TanganiTambahPemasok)
			pemasokGrup.GET("", handlerManufaktur.TanganiAmbilSemuaPemasok)
			pemasokGrup.GET("/:id", handlerManufaktur.TanganiCariPemasokID)
			pemasokGrup.PUT("/:id", handlerManufaktur.TanganiUbahPemasok)
			pemasokGrup.DELETE("/:id", handlerManufaktur.TanganiHapusPemasok)
		}

		materialGrup := api.Group("/materials")
		materialGrup.Use(infrastruktur.PenjagaSesiJWT())
		{
			materialGrup.POST("", handlerManufaktur.TanganiTambahMaterial)
			materialGrup.GET("", handlerManufaktur.TanganiAmbilSemuaMaterial)
			materialGrup.GET("/:id", handlerManufaktur.TanganiCariMaterialID)
			materialGrup.PUT("/:id", handlerManufaktur.TanganiUbahMaterial)
			materialGrup.DELETE("/:id", handlerManufaktur.TanganiHapusMaterial)
			materialGrup.POST("/:id/media", handlerMedia.TanganiUnggahMedia)
		}

		customerGrup := api.Group("/customers")
		customerGrup.Use(infrastruktur.PenjagaSesiJWT())
		{
			customerGrup.POST("", handlerManufaktur.TanganiTambahCustomer)
			customerGrup.GET("", handlerManufaktur.TanganiAmbilSemuaCustomer)
			customerGrup.GET("/:id", handlerManufaktur.TanganiCariCustomerID)
			customerGrup.PUT("/:id", handlerManufaktur.TanganiUbahCustomer)
			customerGrup.DELETE("/:id", handlerManufaktur.TanganiHapusCustomer)
		}

		bomGrup := api.Group("/boms")
		bomGrup.Use(infrastruktur.PenjagaSesiJWT())
		{
			bomGrup.POST("", handlerManufaktur.TanganiTambahBOM)
			bomGrup.GET("", handlerManufaktur.TanganiAmbilSemuaBOM)
			bomGrup.GET("/:id", handlerManufaktur.TanganiCariBOMID)
			bomGrup.PUT("/:id", handlerManufaktur.TanganiUbahBOM)
			bomGrup.DELETE("/:id", handlerManufaktur.TanganiHapusBOM)
		}

		// Endpoint Master Data Snapshot (untuk sinkronisasi QControl)
		api.GET("/master-data/snapshot", infrastruktur.PenjagaSesiJWT(), handlerManufaktur.TanganiSnapshotMasterData)
		api.GET("/qcontrol/master-data", infrastruktur.PenjagaSesiJWT(), handlerManufaktur.TanganiSnapshotMasterData)
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

// LandingData menyimpan konfigurasi dinamis yang dirender ke landing.html
type LandingData struct {
	StatusApp       string
	DBConnected     bool
	DBHost          string
	DBName          string
	AllowedOrigins  string
	IPWhitelist     string
	WaktuServer     string
	SwaggerURL      string
	SupplierCount   int64
	MaterialCount   int64
	CustomerCount   int64
	BOMCount        int64
	InspectionCount int64
}

// tanganiLandingPage merestitusi antarmuka visual berbasis go:embed landing.html
func tanganiLandingPage(db *gorm.DB) gin.HandlerFunc {
	tmpl, err := template.New("landing").Parse(landingHTML)
	if err != nil {
		log.Fatalf("Gagal mem-parsing template landing.html: %v", err)
	}
	return func(k *gin.Context) {
		sqlDB, err := db.DB()
		dbConnected := false
		dbHost := os.Getenv("DB_HOST")
		dbName := os.Getenv("DB_NAME")

		if err == nil {
			errPing := sqlDB.Ping()
			dbConnected = (errPing == nil)
		}

		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "* (Dev Fallback)"
		}

		ipWhitelist := os.Getenv("IP_WHITELIST")
		if ipWhitelist == "" {
			ipWhitelist = "None (Fully Open)"
		}

		var supplierCount int64
		var materialCount int64
		var customerCount int64
		var bomCount int64
		var inspectionCount int64

		if dbConnected {
			db.Model(&manufaktur.Pemasok{}).Count(&supplierCount)
			db.Model(&manufaktur.Material{}).Count(&materialCount)
			db.Model(&manufaktur.Customer{}).Count(&customerCount)
			db.Model(&manufaktur.KomposisiMaterialBOM{}).Count(&bomCount)
			db.Model(&kualitas.LembarPeriksa{}).Count(&inspectionCount)
		}

		data := LandingData{
			StatusApp:       "BETA_ACTIVE",
			DBConnected:     dbConnected,
			DBHost:          dbHost,
			DBName:          dbName,
			AllowedOrigins:  allowedOrigins,
			IPWhitelist:     ipWhitelist,
			WaktuServer:     time.Now().Format("2006-01-02 15:04:05 MST"),
			SwaggerURL:      "/swagger/index.html",
			SupplierCount:   supplierCount,
			MaterialCount:   materialCount,
			CustomerCount:   customerCount,
			BOMCount:        bomCount,
			InspectionCount: inspectionCount,
		}

		k.Header("Content-Type", "text/html; charset=utf-8")
		errExec := tmpl.Execute(k.Writer, data)
		if errExec != nil {
			log.Printf("Gagal merender landing page: %v", errExec)
			k.String(http.StatusInternalServerError, "Gagal merender landing page")
		}
	}
}

// TanganiCekSistem menyajikan status operasional aplikasi dan pangkalan data.
// @Summary Pemeriksaan Kesehatan Sistem
// @Description Memvalidasi kesiapan server API dan konektivitas live database
// @Tags Sistem
// @Produce json
// @Success 200 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/cek_sistem [get]
func TanganiCekSistem(db *gorm.DB) gin.HandlerFunc {
	return func(k *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			respon.Galat_Server(k, "Gagal mendapatkan koneksi pangkalan data.", err)
			return
		}
		errPing := sqlDB.Ping()
		if errPing != nil {
			respon.Galat_Server(k, "Pangkalan data tidak dapat dijangkau.", errPing)
			return
		}
		respon.Sukses(k, "Sistem PGNServer beroperasi secara optimal dan terhubung ke pangkalan data.", nil)
	}
}

// TanganiStatusKrusial menyajikan status akses ter-whitelist IP.
// @Summary Status Operasi Krusial (Whitelisted IP)
// @Description Menyajikan data status sensitif jika IP pengirim lolos Whitelisting korporasi
// @Tags Sistem
// @Produce json
// @Security BearerAuth
// @Success 200 {object} respon.ResponStandar
// @Failure 401 {object} respon.ResponStandar
// @Router /api/v1/krusial/status [get]
func TanganiStatusKrusial() gin.HandlerFunc {
	return func(k *gin.Context) {
		respon.Sukses(k, "Akses diizinkan, IP kamu terdaftar dalam whitelist pangkalan data.", gin.H{
			"ip_whitelist_status": "authorized",
			"timestamp":           time.Now(),
		})
	}
}

// TanganiHealth menyajikan pemeriksaan liveness sistem.
// @Summary Pemeriksaan Liveness Aplikasi
// @Description Memvalidasi bahwa instansi server API aktif dan berjalan
// @Tags Sistem
// @Produce json
// @Success 200 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/health [get]
func TanganiHealth(db *gorm.DB) gin.HandlerFunc {
	return TanganiCekSistem(db)
}

// TanganiKesehatan menyajikan status operasional aplikasi dan pangkalan data.
// @Summary Pemeriksaan Kesehatan Sistem KMP
// @Description Memvalidasi kesiapan server API dan konektivitas live database untuk klien KMP
// @Tags Sistem
// @Produce json
// @Success 200 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/kesehatan [get]
func TanganiKesehatan(db *gorm.DB) gin.HandlerFunc {
	return TanganiCekSistem(db)
}

// TanganiReadiness menyajikan pemeriksaan kesiapan sistem.
// @Summary Pemeriksaan Kesiapan Aplikasi
// @Description Memvalidasi bahwa server API siap menerima trafik dengan memverifikasi koneksi database
// @Tags Sistem
// @Produce json
// @Success 200 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/readiness [get]
func TanganiReadiness(db *gorm.DB) gin.HandlerFunc {
	return func(k *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			respon.Galat_Server(k, "Gagal mendapatkan koneksi pangkalan data.", err)
			return
		}
		errPing := sqlDB.Ping()
		if errPing != nil {
			respon.Galat_Server(k, "Pangkalan data tidak dapat dijangkau.", errPing)
			return
		}
		respon.Sukses(k, "Aplikasi PGNServer siap menerima koneksi trafik.", nil)
	}
}
