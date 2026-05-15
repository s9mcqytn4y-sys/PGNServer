package konfigurasi

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"pgn-server/internal/model"
)

// DB adalah instance global untuk database
var DB *gorm.DB

// HubungkanDatabase menginisialisasi koneksi dan pool ke PostgreSQL
func HubungkanDatabase() {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan objek sqlDB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Melakukan migrasi database berdasarkan skema baru
	err = db.AutoMigrate(
		&model.User{},
		&model.Customer{},
		&model.Supplier{},
		&model.LiniProduksi{},
		&model.Material{},
		&model.Produk{},
		&model.BillOfMaterial{},
		&model.DefectMaster{},
		&model.InspeksiHarian{},
		&model.LogInspeksi{},
		&model.BukuBesarDefectMaterial{},
		&model.ChecksheetHeader{},
		&model.ChecksheetDetail{},
	)
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	DB = db
	log.Println("Database berhasil terhubung dan dimigrasi.")

	// Jalankan Seeding
	JalankanSeeder(db)
}
