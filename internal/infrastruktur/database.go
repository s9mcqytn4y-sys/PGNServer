// Package infrastruktur memuat fungsi dasar infrastruktur seperti koneksi DB.
package infrastruktur

import (
	"log"

	"pgn-server/internal/kualitas"
	"pgn-server/internal/manufaktur"
	"pgn-server/internal/otentikasi"
	"pgn-server/internal/media"

	"gorm.io/gorm"
)

// PelaksanaanAutoMigrasi mengatur proses translasi tabel pangkalan data
// berdasarkan abstraksi entitas di dalam subsistem aplikasi.
func PelaksanaanAutoMigrasi(db *gorm.DB) error {
	log.Println("Memulai sinkronisasi abstraksi automigrasi GORM...")

	err := db.AutoMigrate(
		&otentikasi.Pengguna{},
		&manufaktur.Customer{},
		&manufaktur.Pemasok{},
		&manufaktur.Material{},
		&manufaktur.KomposisiMaterialBOM{},
		&kualitas.LembarPeriksa{},
		&kualitas.DetailInspeksi{},
		&kualitas.BukuBesarCacat{},
		&media.AsetDigital{},
	)

	if err != nil {
		log.Printf("Galat sewaktu automigrasi: %v\n", err)
		return err
	}

	log.Println("Penyatuan skema pangkalan data Postgres 17.x berhasil.")
	return nil
}
