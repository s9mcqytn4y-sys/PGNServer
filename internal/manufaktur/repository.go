package manufaktur

import (
	"gorm.io/gorm"
)

// RepositoriManufaktur mendefinisikan kontrak operasi CRUD data master manufaktur.
type RepositoriManufaktur interface {
	// CRUD Pemasok (Supplier)
	BuatPemasok(p *Pemasok) error
	CariPemasokBerdasarkanID(id uint) (*Pemasok, error)
	DaftarPemasok() ([]Pemasok, error)
	PerbaruiPemasok(p *Pemasok) error
	HapusPemasok(id uint) error

	// CRUD Material
	BuatMaterial(m *Material) error
	CariMaterialBerdasarkanID(id uint) (*Material, error)
	CariMaterialBerdasarkanSKU(sku string) (*Material, error)
	DaftarMaterial() ([]Material, error)
	PerbaruiMaterial(m *Material) error
	HapusMaterial(id uint) error

	// CRUD Customer
	BuatCustomer(c *Customer) error
	CariCustomerBerdasarkanID(id uint) (*Customer, error)
	DaftarCustomer() ([]Customer, error)
	PerbaruiCustomer(c *Customer) error
	HapusCustomer(id uint) error

	// CRUD BOM
	BuatBOM(b *KomposisiMaterialBOM) error
	CariBOMBerdasarkanID(id uint) (*KomposisiMaterialBOM, error)
	DaftarBOM() ([]KomposisiMaterialBOM, error)
	PerbaruiBOM(b *KomposisiMaterialBOM) error
	HapusBOM(id uint) error
}

type repositoriManufaktur struct {
	db *gorm.DB
}

// KonstruksiRepositoriBaru membuat instance baru RepositoriManufaktur.
func KonstruksiRepositoriBaru(db *gorm.DB) RepositoriManufaktur {
	return &repositoriManufaktur{db: db}
}

// === IMPLEMENTASI CRUD PEMASOK ===

func (r *repositoriManufaktur) BuatPemasok(p *Pemasok) error {
	return r.db.Create(p).Error
}

func (r *repositoriManufaktur) CariPemasokBerdasarkanID(id uint) (*Pemasok, error) {
	var p Pemasok
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repositoriManufaktur) DaftarPemasok() ([]Pemasok, error) {
	var list []Pemasok
	if err := r.db.Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repositoriManufaktur) PerbaruiPemasok(p *Pemasok) error {
	return r.db.Save(p).Error
}

func (r *repositoriManufaktur) HapusPemasok(id uint) error {
	return r.db.Delete(&Pemasok{}, id).Error
}

// === IMPLEMENTASI CRUD MATERIAL ===

func (r *repositoriManufaktur) BuatMaterial(m *Material) error {
	return r.db.Create(m).Error
}

func (r *repositoriManufaktur) CariMaterialBerdasarkanID(id uint) (*Material, error) {
	var m Material
	if err := r.db.Preload("Pemasok").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repositoriManufaktur) CariMaterialBerdasarkanSKU(sku string) (*Material, error) {
	var m Material
	if err := r.db.Preload("Pemasok").Where("kode_sku = ?", sku).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repositoriManufaktur) DaftarMaterial() ([]Material, error) {
	var list []Material
	if err := r.db.Preload("Pemasok").Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repositoriManufaktur) PerbaruiMaterial(m *Material) error {
	return r.db.Save(m).Error
}

func (r *repositoriManufaktur) HapusMaterial(id uint) error {
	return r.db.Delete(&Material{}, id).Error
}

// === IMPLEMENTASI CRUD CUSTOMER ===

func (r *repositoriManufaktur) BuatCustomer(c *Customer) error {
	return r.db.Create(c).Error
}

func (r *repositoriManufaktur) CariCustomerBerdasarkanID(id uint) (*Customer, error) {
	var c Customer
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repositoriManufaktur) DaftarCustomer() ([]Customer, error) {
	var list []Customer
	if err := r.db.Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repositoriManufaktur) PerbaruiCustomer(c *Customer) error {
	return r.db.Save(c).Error
}

func (r *repositoriManufaktur) HapusCustomer(id uint) error {
	return r.db.Delete(&Customer{}, id).Error
}

// === IMPLEMENTASI CRUD BOM ===

func (r *repositoriManufaktur) BuatBOM(b *KomposisiMaterialBOM) error {
	return r.db.Create(b).Error
}

func (r *repositoriManufaktur) CariBOMBerdasarkanID(id uint) (*KomposisiMaterialBOM, error) {
	var b KomposisiMaterialBOM
	if err := r.db.Preload("ParentMaterial").Preload("MaterialBaku").First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *repositoriManufaktur) DaftarBOM() ([]KomposisiMaterialBOM, error) {
	var list []KomposisiMaterialBOM
	if err := r.db.Preload("ParentMaterial").Preload("MaterialBaku").Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repositoriManufaktur) PerbaruiBOM(b *KomposisiMaterialBOM) error {
	return r.db.Save(b).Error
}

func (r *repositoriManufaktur) HapusBOM(id uint) error {
	return r.db.Delete(&KomposisiMaterialBOM{}, id).Error
}
