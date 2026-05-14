package model

import "time"

type Produk struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NamaPart  string    `gorm:"type:varchar(255);not null" json:"nama_part"`
	NomorPart string    `gorm:"type:varchar(100);unique;not null" json:"nomor_part"`
	CreatedAt time.Time `json:"dibuat_pada"`
	UpdatedAt time.Time `json:"diperbarui_pada"`
}
