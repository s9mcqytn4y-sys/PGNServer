// Package kualitas memuat fungsionalitas repository untuk inspeksi kualitas.
package kualitas

import (
	"encoding/json"
	"gorm.io/gorm"
)

// RepositoriKualitas menyediakan lapisan operasi pangkalan data spesifik kualitas.
type RepositoriKualitas interface {
	SimpanLembarPeriksaMassal(dto *DTOLembarPeriksaKirim, tx *gorm.DB) error
}

type repositoriKualitas struct{}

func KonstruksiRepositoriBaru() RepositoriKualitas {
	return &repositoriKualitas{}
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

	// 3. Transmisi tunggal seketika memanfaatkan JSON_TABLE (Fitur Terobosan Postgres 17)
	// Kita map secara spesifik path tiap variabel untuk penyisipan atomik.
	kueriSQL := `
		INSERT INTO detail_inspeksis (lembar_periksa_id, unik_part_id, kode_cacat, waktu_pergeseran, rasio_cacat, rasio_total_ok)
		SELECT ?, unik_part_id, kode_cacat, waktu_pergeseran, rasio_cacat, rasio_total_ok
		FROM JSON_TABLE(
			?::jsonb,
			'$[*]' COLUMNS (
				unik_part_id bigint PATH '$.unik_part_id',
				kode_cacat text PATH '$.kode_cacat',
				waktu_pergeseran text PATH '$.waktu_pergeseran',
				rasio_cacat numeric PATH '$.rasio_cacat',
				rasio_total_ok numeric PATH '$.rasio_total_ok'
			)
		)
	`
	errEksekusi := tx.Exec(kueriSQL, lembar.ID, string(dataJSON)).Error
	return errEksekusi
}
