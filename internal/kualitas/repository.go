// Package kualitas memuat fungsionalitas repository untuk inspeksi kualitas.
package kualitas

import (
	"encoding/json"
	"gorm.io/gorm"
)

// RepositoriKualitas menyediakan lapisan operasi pangkalan data spesifik kualitas.
type RepositoriKualitas interface {
	SimpanLembarPeriksaMassal(dto *DTOLembarPeriksaKirim, tx *gorm.DB) error
	DaftarRiwayat(limit int, offset int, tanggalMulai string, tanggalSelesai string, zonaLini string) ([]LembarPeriksa, error)
}

type repositoriKualitas struct {
	db *gorm.DB
}

func KonstruksiRepositoriBaru(db *gorm.DB) RepositoriKualitas {
	return &repositoriKualitas{db: db}
}

// SimpanLembarPeriksaMassal meniadakan pola N+1 rekursif ORM dan menggunakan SQL JSON_TABLE.
func (r *repositoriKualitas) SimpanLembarPeriksaMassal(dto *DTOLembarPeriksaKirim, tx *gorm.DB) error {
	// 1. Simpan baris induk ke tabel lembar_periksas
	lembar := LembarPeriksa{
		Tanggal:            dto.Tanggal,
		ZonaLini:           dto.ZonaLini,
		PenggunaIDTercatat: dto.PenggunaIDTercatat,
	}
	if err := tx.Create(&lembar).Error; err != nil {
		return err
	}

	// 2. Persiapkan array komposit menjadi JSON String murni
	dataJSON, errJSON := json.Marshal(dto.Detail)
	if errJSON != nil {
		return errJSON
	}

	// 3. Transmisi tunggal seketika memanfaatkan jsonb_to_recordset (Native PostgreSQL JSON)
	// Kita map secara spesifik path tiap variabel untuk penyisipan atomik sesuai dengan nama JSON tag.
	kueriSQL := `
		INSERT INTO detail_inspeksis (lembar_periksa_id, unik_part_id, kode_cacat, waktu_pergeseran, rasio_cacat, rasio_total_ok)
		SELECT ?, "unikPartId", "kodeCacat", "waktuPergeseran", "rasioCacat", "rasioTotalOK"
		FROM jsonb_to_recordset(?::jsonb) AS x(
			"unikPartId" bigint,
			"kodeCacat" text,
			"waktuPergeseran" text,
			"rasioCacat" numeric,
			"rasioTotalOK" numeric,
			"totalProduksi" numeric
		)
	`
	errEksekusi := tx.Exec(kueriSQL, lembar.ID, string(dataJSON)).Error
	return errEksekusi
}

// DaftarRiwayat melakukan paginasi dan filter historis lembar periksa.
func (r *repositoriKualitas) DaftarRiwayat(limit int, offset int, tanggalMulai string, tanggalSelesai string, zonaLini string) ([]LembarPeriksa, error) {
	var riwayat []LembarPeriksa

	query := r.db.Model(&LembarPeriksa{})

	if tanggalMulai != "" && tanggalSelesai != "" {
		query = query.Where("tanggal >= ? AND tanggal <= ?", tanggalMulai, tanggalSelesai)
	}
	if zonaLini != "" {
		query = query.Where("zona_lini = ?", zonaLini)
	}

	err := query.Limit(limit).Offset(offset).Order("tanggal DESC").Find(&riwayat).Error
	return riwayat, err
}
