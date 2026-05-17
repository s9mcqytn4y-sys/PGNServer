package otentikasi

import "gorm.io/gorm"

// RepositoriOtentikasi merepresentasikan kontrak akses data otentikasi.
type RepositoriOtentikasi interface {
	Daftar(pengguna *Pengguna) error
	CariBerdasarkanSurel(surel string) (*Pengguna, error)
	PerbaruiSandi(pengguna *Pengguna) error
	PerbaruiTenggatSesi(id uint, tenggatSesi *gorm.DeletedAt) error // atau tipe time opsional
}

type repositoriOtentikasi struct {
	db *gorm.DB
}

func EkstraksiRepositoriBaru(db *gorm.DB) RepositoriOtentikasi {
	return &repositoriOtentikasi{db}
}

func (r *repositoriOtentikasi) Daftar(pengguna *Pengguna) error {
	return r.db.Create(pengguna).Error
}

func (r *repositoriOtentikasi) CariBerdasarkanSurel(surel string) (*Pengguna, error) {
	var pengguna Pengguna
	err := r.db.Where("surel_kredensial = ?", surel).First(&pengguna).Error
	if err != nil {
		return nil, err
	}
	return &pengguna, nil
}

func (r *repositoriOtentikasi) PerbaruiSandi(pengguna *Pengguna) error {
	return r.db.Model(pengguna).Update("kata_sandi_terenskripsi", pengguna.KataSandiTerenskripsi).Error
}

func (r *repositoriOtentikasi) PerbaruiTenggatSesi(id uint, tenggatSesi *gorm.DeletedAt) error {
	return r.db.Model(&Pengguna{}).Where("id = ?", id).Update("tenggat_sesi_aktif", tenggatSesi).Error
}
