package main

import (
	"log"
	"os"
	"pgn-server/internal/konfigurasi"
	"pgn-server/internal/pengelola"
	"pgn-server/pkg/utilitas"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
)

func main() {
	// 1. Memuat variabel environment dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment OS")
	}

	// 2. Hubungkan ke database PostgreSQL dan jalankan migrasi
	konfigurasi.HubungkanDatabase()

	// 3. Konfigurasi router Gin
	r := gin.Default()

	// 4. Rute Statis untuk penyajian gambar (folder sudah dipersiapkan)
	r.Static("/gambar", "./penyimpanan/gambar")

	// 5. RESTful API Routes
	// Health check (tidak perlu autentikasi)
	r.GET("/api/v1/kesehatan", func(c *gin.Context) {
		c.JSON(http.StatusOK, utilitas.FormatRespons(true, "Server PGN berjalan dengan baik", nil))
	})

	// Grup API v1 dengan Middleware Autentikasi
	api := r.Group("/api/v1")
	// Middleware dinonaktifkan sementara untuk kemudahan testing, buka komentar untuk mengaktifkan
	// api.Use(pengelola.AutentikasiMiddleware())
	{
		api.GET("/produk", pengelola.AmbilSemuaProduk)
		api.POST("/produk", pengelola.TambahProduk)
	}

	// 6. Jalankan Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server mulai berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Kesalahan saat menjalankan server: %v", err)
	}
}
