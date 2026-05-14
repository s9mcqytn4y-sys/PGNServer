package model

import "time"

// Produk merepresentasikan data produk di aplikasi
type Produk struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	KodePart   string    `gorm:"uniqueIndex;not null" json:"kode_part"`
	Nama       string    `gorm:"not null" json:"nama"`
	DibuatPada time.Time `gorm:"autoCreateTime" json:"dibuat_pada"`
	DiubahPada time.Time `gorm:"autoUpdateTime" json:"diubah_pada"`
}

// Material merepresentasikan data material/bahan baku
type Material struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Kode       string    `gorm:"uniqueIndex;not null" json:"kode"`
	Nama       string    `gorm:"not null" json:"nama"`
	Pemasok    string    `json:"pemasok"`
	DibuatPada time.Time `gorm:"autoCreateTime" json:"dibuat_pada"`
	DiubahPada time.Time `gorm:"autoUpdateTime" json:"diubah_pada"`
}

// BillOfMaterial merepresentasikan relasi produk dan material
type BillOfMaterial struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProdukID   uint      `gorm:"not null;index" json:"produk_id"`
	MaterialID uint      `gorm:"not null;index" json:"material_id"`
	Jumlah     int       `json:"jumlah"`
	Produk     Produk    `gorm:"foreignKey:ProdukID" json:"produk"`
	Material   Material  `gorm:"foreignKey:MaterialID" json:"material"`
}

// RiwayatDefect merepresentasikan log kerusakan atau cacat
type RiwayatDefect struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProdukID      uint      `gorm:"not null;index" json:"produk_id"`
	KategoriCacat string    `gorm:"index;not null" json:"kategori_cacat"`
	Jumlah        int       `json:"jumlah"`
	Keterangan    string    `json:"keterangan"`
	Tanggal       time.Time `json:"tanggal"`
	DibuatPada    time.Time `gorm:"autoCreateTime" json:"dibuat_pada"`
	Produk        Produk    `gorm:"foreignKey:ProdukID" json:"produk"`
}
