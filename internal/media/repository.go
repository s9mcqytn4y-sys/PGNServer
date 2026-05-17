package media

import (
	"gorm.io/gorm"
)

// RepositoriMedia mengatur akses data persisten untuk aset media digital.
type RepositoriMedia interface {
	Simpan(aset *AsetDigital) error
	HapusBerdasarkanID(id uint) error
	CariBerdasarkanID(id uint) (*AsetDigital, error)
	DapatkanBerdasarkanIDMaterial(idMaterial uint) ([]*AsetDigital, error)
}

type repositoriMedia struct {
	db *gorm.DB
}

// KonstruksiRepositoriBaru membuat instance baru RepositoriMedia.
func KonstruksiRepositoriBaru(db *gorm.DB) RepositoriMedia {
	return &repositoriMedia{db: db}
}

func (r *repositoriMedia) Simpan(aset *AsetDigital) error {
	return r.db.Save(aset).Error
}

func (r *repositoriMedia) HapusBerdasarkanID(id uint) error {
	return r.db.Delete(&AsetDigital{}, id).Error
}

func (r *repositoriMedia) CariBerdasarkanID(id uint) (*AsetDigital, error) {
	var aset AsetDigital
	err := r.db.First(&aset, id).Error
	if err != nil {
		return nil, err
	}
	return &aset, nil
}

func (r *repositoriMedia) DapatkanBerdasarkanIDMaterial(idMaterial uint) ([]*AsetDigital, error) {
	var daftarAset []*AsetDigital
	err := r.db.Where("id_material = ?", idMaterial).Find(&daftarAset).Error
	if err != nil {
		return nil, err
	}
	return daftarAset, nil
}
