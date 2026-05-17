package media

import (
	"time"
)

// AsetDigital merepresentasikan berkas media (foto cacat/komponen) yang terhubung ke material.
type AsetDigital struct {
	ID              uint      `gorm:"primaryKey;column:id" json:"id"`
	IDMaterial      uint      `gorm:"column:id_material;index;not null" json:"idMaterial"`
	TipeMIME        string    `gorm:"column:tipe_mime;type:varchar(50);not null" json:"tipeMime"`
	UkuranBerkas    int64     `gorm:"column:ukuran_berkas;not null" json:"ukuranBerkas"`
	Ekstensi        string    `gorm:"column:ekstensi;type:varchar(10);not null" json:"ekstensi"`
	DirektoriLokal  string    `gorm:"column:direktori_lokal;type:varchar(500)" json:"direktoriLokal"`
	TautanEksternal string    `gorm:"column:tautan_eksternal;type:varchar(1000)" json:"tautanEksternal"`
	TipePenyimpanan string    `gorm:"column:tipe_penyimpanan;type:varchar(20);not null" json:"tipePenyimpanan"` // "LOKAL" atau "EKSTERNAL"
	DibuatPada      time.Time `gorm:"autoCreateTime" json:"dibuatPada"`
	DiubahPada      time.Time `gorm:"autoUpdateTime" json:"diubahPada"`
}
