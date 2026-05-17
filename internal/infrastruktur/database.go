package infrastruktur

import (
	"log"

	"pgn-server/internal/manufaktur"
	"pgn-server/internal/otentikasi"

	"gorm.io/gorm"
)

// PelaksanaanAutoMigrasi mengatur proses translasi tabel pangkalan data
// berdasarkan abstraksi entitas di dalam subsistem aplikasi.
func PelaksanaanAutoMigrasi(db *gorm.DB) error {
	log.Println("Memulai sinkronisasi abstraksi automigrasi GORM...")

	err := db.AutoMigrate(
		&otentikasi.Pengguna{},
		&manufaktur.Pemasok{},
		&manufaktur.Material{},
		&manufaktur.KomposisiMaterialBOM{},
	)

	if err != nil {
		log.Printf("Galat sewaktu automigrasi: %v\n", err)
		return err
	}

	log.Println("Penyatuan skema pangkalan data Postgres 17.x berhasil.")
	return nil
}
