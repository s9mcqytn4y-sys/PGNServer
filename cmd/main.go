package main

import (
	"log"
	"pgn-server/internal/konfigurasi"
	"pgn-server/internal/pengelola"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Muat .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan")
	}

	// 2. Inisialisasi DB
	konfigurasi.HubungkanDatabase()

	// 3. Setup Router dengan Gin
	r := gin.Default()

	// Keamanan Dasar: Recovery dari Panic & Logger
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Grup API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/produk", pengelola.AmbilSemuaProduk)
		v1.POST("/produk", pengelola.SimpanProduk)
	}

	log.Println("Server PGNServer berjalan di port 8080")
	r.Run(":8080")
}
