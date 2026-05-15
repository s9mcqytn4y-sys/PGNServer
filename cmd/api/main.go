package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"pgn-server/internal/konfigurasi"
	"pgn-server/internal/pengelola"
	"pgn-server/pkg/utilitas"
	_ "pgn-server/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

// @title PGNServer API
// @version 1.0
// @description API Server untuk PGN (Produk Gagal & NG) Intelligence System.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.pgn.co.id/support
// @contact.email support@pgn.co.id

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func init() {
	// Inisialisasi slog sebagai default logger
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))
}

func main() {
	// Memuat file konfigurasi .env
	if err := godotenv.Load(); err != nil {
		slog.Warn("File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	// Hubungkan ke database
	konfigurasi.HubungkanDatabase()

	// Inisialisasi router Gin
	r := gin.New()

	// --- Middleware Global ---
	r.Use(gin.Recovery())                 // Recovery dari panic
	r.Use(pengelola.LogSlogMiddleware())  // Structured logging
	r.Use(pengelola.TimeoutMiddleware(30 * time.Second)) // Global timeout
	r.Use(pengelola.LimitRequestMiddleware(rate.Limit(100), 200)) // Rate limit 100 req/s

	// Statis & Gambar
	r.Static("/penyimpanan", "./penyimpanan")

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routing API v1
	v1 := r.Group("/api/v1")
	{
		// Health check (Tanpa Auth)
		v1.GET("/kesehatan", func(c *gin.Context) {
			c.JSON(http.StatusOK, utilitas.FormatRespons(true, "Server PGN berjalan dengan baik", gin.H{
				"waktu_server": time.Now().Format(time.RFC3339),
				"status":       "stabil",
			}))
		})

		// Produk (Dengan Auth)
		produk := v1.Group("/produk")
		produk.Use(pengelola.AutentikasiMiddleware())
		{
			produk.GET("/", pengelola.AmbilSemuaProduk)
			produk.POST("/", pengelola.TambahProduk)
		}
	}

	// Menjalankan server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server PGNServer dimulai", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("Gagal menjalankan server", "error", err)
		os.Exit(1)
	}
}
