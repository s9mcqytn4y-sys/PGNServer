package kualitas

import "time"

// DTO_Detail_Inspeksi mendefinisikan detail inspeksi dalam JSON.
type DTO_Detail_Inspeksi struct {
	UnikPartID      uint    `json:"unik_part_id" binding:"required"`
	KodeCacat       string  `json:"kode_cacat" binding:"required"`
	WaktuPergeseran string  `json:"waktu_pergeseran" binding:"required"`
	RasioCacat      float64 `json:"rasio_cacat" binding:"gte=0"`
	RasioTotalOK    float64 `json:"rasio_total_ok" binding:"gte=0"`
}

// DTO_LembarPeriksa_Kirim adalah komposit majemuk payload masuk.
type DTO_LembarPeriksa_Kirim struct {
	Tanggal            string                `json:"tanggal" binding:"required"`
	ZonaLini           string                `json:"zona_lini" binding:"required"`
	PenggunaIDTercatat uint                  `json:"pengguna_id_tercatat" binding:"required"`
	Detail             []DTO_Detail_Inspeksi `json:"detail" binding:"required,dive"`
}

// LembarPeriksa merepresentasikan data utama lembar inspeksi fisik QC.
type LembarPeriksa struct {
	ID                 uint      `gorm:"primaryKey;column:id"`
	Tanggal            string    `gorm:"column:tanggal;type:date"`
	ZonaLini           string    `gorm:"column:zona_lini"`
	PenggunaIDTercatat uint      `gorm:"column:pengguna_id_tercatat"`
	DibuatPada         time.Time `gorm:"autoCreateTime"`
}

// DetailInspeksi merepresentasikan komponen rasio cacat harian.
type DetailInspeksi struct {
	ID              uint    `gorm:"primaryKey;column:id"`
	LembarPeriksaID uint    `gorm:"column:lembar_periksa_id"`
	UnikPartID      uint    `gorm:"column:unik_part_id"`
	KodeCacat       string  `gorm:"column:kode_cacat"`
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
