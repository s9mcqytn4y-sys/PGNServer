// Package kualitas menangani modul pencatatan inspeksi kontrol kualitas.
package kualitas

import "time"

// DTODetailInspeksi mendefinisikan detail inspeksi dalam JSON dengan struktur TPS.
type DTODetailInspeksi struct {
	UnikPartID      uint    `json:"unikPartId,omitempty" binding:"required"`
	KodeCacat       string  `json:"kodeCacat,omitempty" binding:"required"`
	WaktuPergeseran string  `json:"waktuPergeseran,omitempty" binding:"required"` // Shift
	TotalProduksi   float64 `json:"totalProduksi,omitempty" binding:"gte=0"`
	RasioTotalOK    float64 `json:"rasioTotalOK,omitempty" binding:"gte=0"`
	RasioCacat      float64 `json:"rasioCacat,omitempty" binding:"gte=0"` // NG
}

// DTOLembarPeriksaKirim adalah komposit majemuk payload masuk.
type DTOLembarPeriksaKirim struct {
	Tanggal            string              `json:"tanggal,omitempty" binding:"required"`
	ZonaLini           string              `json:"zonaLini,omitempty" binding:"required"`
	PenggunaIDTercatat uint                `json:"penggunaIdTercatat,omitempty" binding:"required"`
	Detail             []DTODetailInspeksi `json:"detail,omitempty" binding:"required,dive"`
}

// LembarPeriksa merepresentasikan data utama lembar inspeksi fisik QC.
type LembarPeriksa struct {
	ID                 uint      `gorm:"primaryKey;column:id"`
	Tanggal            string    `gorm:"column:tanggal;type:date;index:idx_lembar_periksa_tanggal_lini,priority:1"`
	ZonaLini           string    `gorm:"column:zona_lini;index:idx_lembar_periksa_tanggal_lini,priority:2"`
	PenggunaIDTercatat uint      `gorm:"column:pengguna_id_tercatat"`
	DibuatPada         time.Time `gorm:"autoCreateTime"`
}

// DetailInspeksi merepresentasikan komponen rasio cacat harian.
type DetailInspeksi struct {
	ID              uint    `gorm:"primaryKey;column:id"`
	LembarPeriksaID uint    `gorm:"column:lembar_periksa_id;index:idx_detail_inspeksi_lookup,priority:1"`
	UnikPartID      uint    `gorm:"column:unik_part_id;index:idx_detail_inspeksi_lookup,priority:2"`
	KodeCacat       string  `gorm:"column:kode_cacat;index:idx_detail_inspeksi_lookup,priority:3"`
	WaktuPergeseran string  `gorm:"column:waktu_pergeseran"`
	RasioCacat      float64 `gorm:"column:rasio_cacat"`
	RasioTotalOK    float64 `gorm:"column:rasio_total_ok"`
}

// BukuBesarCacat merepresentasikan entri pencatatan rasio penyusutan.
type BukuBesarCacat struct {
	ID              uint      `gorm:"primaryKey;column:id"`
	IDMaterial      uint      `gorm:"column:id_material"`
	TotalPenyusutan float64   `gorm:"column:total_penyusutan"`
	DibuatPada      time.Time `gorm:"autoCreateTime"`
}
