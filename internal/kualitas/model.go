// Package kualitas menangani modul pencatatan inspeksi kontrol kualitas.
package kualitas

import "time"

// DTODetailInspeksi mendefinisikan detail inspeksi dalam JSON dengan struktur TPS.
type DTODetailInspeksi struct {
	UnikPartID      uint    `json:"unikPartId" example:"1" binding:"required"`          // ID unik part material yang diperiksa
	KodeCacat       string  `json:"kodeCacat" example:"A" binding:"required"`           // Kode cacat defect (misal: "A", "B", dll)
	WaktuPergeseran string  `json:"waktuPergeseran" example:"08:00" binding:"required"` // Waktu/Shift pergeseran jam pemeriksaan
	TotalProduksi   float64 `json:"totalProduksi" example:"100" binding:"gte=0"`        // Total kuantitas barang yang diproduksi
	RasioTotalOK    float64 `json:"rasioTotalOK" example:"98" binding:"gte=0"`          // Kuantitas barang dengan status OK (Lolos QC)
	RasioCacat      float64 `json:"rasioCacat" example:"2" binding:"gte=0"`             // Kuantitas barang reject / NG (Not Good)
}

// DTOLembarPeriksaKirim adalah komposit majemuk payload masuk.
type DTOLembarPeriksaKirim struct {
	Tanggal            string              `json:"tanggal" example:"2026-05-17T00:00:00Z" binding:"required"` // Tanggal pemeriksaan (format ISO 8601)
	ZonaLini           string              `json:"zonaLini" example:"Lini 1" binding:"required"`              // Zona atau Lini produksi tempat inspeksi
	Shift              string              `json:"shift,omitempty" example:"NORMAL"`                          // Shift kerja (hanya NORMAL yang didukung)
	PenggunaIDTercatat uint                `json:"penggunaIdTercatat" example:"1" binding:"required"`         // ID pengguna/staf QC yang merekam data
	Detail             []DTODetailInspeksi `json:"detail" binding:"required,dive"`                            // Himpunan baris rincian inspeksi part
}

// OpsiLembarPeriksaDto mendefinisikan opsi dinamis untuk pengisian lembar periksa.
type OpsiLembarPeriksaDto struct {
	Shifts    []string `json:"shifts"`
	ZonaLini  []string `json:"zonaLini"`
	TimeSlots []string `json:"timeSlots"`
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
