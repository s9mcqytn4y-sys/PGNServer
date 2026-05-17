package otentikasi

import "time"

// Pengguna merepresentasikan entitas akun untuk autentikasi dan otorisasi.
type Pengguna struct {
	ID                    uint       `gorm:"primaryKey;column:id" json:"id"`
	SurelKredensial       string     `gorm:"column:surel_kredensial;type:varchar(255);uniqueIndex;not null" json:"surelKredensial"`
	KataSandiTerenskripsi string     `gorm:"column:kata_sandi_terenskripsi;type:varchar(255);not null" json:"-"` // Tidak dirender dalam JSON v2
	PeranOtorisasi        string     `gorm:"column:peran_otorisasi;type:varchar(50);default:'OPERATOR'" json:"peranOtorisasi"`
	TenggatSesiAktif      *time.Time `gorm:"column:tenggat_sesi_aktif" json:"tenggatSesiAktif,omitempty"`
	DibuatPada            time.Time  `gorm:"autoCreateTime" json:"dibuatPada"`
	DiubahPada            time.Time  `gorm:"autoUpdateTime" json:"diubahPada"`
}
