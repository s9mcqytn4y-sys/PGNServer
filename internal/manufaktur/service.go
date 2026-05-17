package manufaktur

import (
	"errors"
	"strconv"
	"time"

	"pgn-server/pkg/cache"
)

// LayananManufaktur membungkus logika bisnis manipulasi data master manufaktur.
type LayananManufaktur interface {
	// Pemasok (Supplier)
	TambahPemasok(dto *DTOPemasokSimpan) (*Pemasok, error)
	CariPemasokID(id uint) (*Pemasok, error)
	AmbilSemuaPemasok() ([]Pemasok, error)
	UbahPemasok(id uint, dto *DTOPemasokSimpan) (*Pemasok, error)
	HapusPemasokID(id uint) error

	// Material
	TambahMaterial(dto *DTOMaterialSimpan) (*Material, error)
	CariMaterialID(id uint) (*Material, error)
	AmbilSemuaMaterial() ([]Material, error)
	UbahMaterial(id uint, dto *DTOMaterialSimpan) (*Material, error)
	HapusMaterialID(id uint) error

	// Customer
	TambahCustomer(dto *DTOCustomerSimpan) (*Customer, error)
	CariCustomerID(id uint) (*Customer, error)
	AmbilSemuaCustomer() ([]Customer, error)
	UbahCustomer(id uint, dto *DTOCustomerSimpan) (*Customer, error)
	HapusCustomerID(id uint) error

	// BOM
	TambahBOM(dto *DTOBOMSimpan) (*KomposisiMaterialBOM, error)
	CariBOMID(id uint) (*KomposisiMaterialBOM, error)
	AmbilSemuaBOM() ([]KomposisiMaterialBOM, error)
	UbahBOM(id uint, dto *DTOBOMSimpan) (*KomposisiMaterialBOM, error)
	HapusBOMID(id uint) error

	// Snapshot
	AmbilSnapshotMasterData() (*MasterDataSnapshotDto, *MetadataSnapshotDto, error)
}

type layananManufaktur struct {
	repo RepositoriManufaktur
}

// KonstruksiLayananBaru membuat instance baru LayananManufaktur.
func KonstruksiLayananBaru(repo RepositoriManufaktur) LayananManufaktur {
	return &layananManufaktur{repo: repo}
}

// DTO structs for payload input validation
type DTOPemasokSimpan struct {
	SupplierCode string `json:"supplierCode" binding:"required"`
	NamaEntitas  string `json:"namaEntitas" binding:"required"`
	Kontak       string `json:"kontak"`
}

type DTOMaterialSimpan struct {
	KodeSKU      string  `json:"kodeSKU" binding:"required"`
	NamaMaterial string  `json:"namaMaterial" binding:"required"`
	TebalCM      float64 `json:"tebalCM"`
	BeratGSM     int     `json:"beratGSM"`
	LebarCM      float64 `json:"lebarCM"`
	PanjangCM    float64 `json:"panjangCM"`
	UnitSatuan   string  `json:"unitSatuan" binding:"required"`
	IDPemasok    uint    `json:"idPemasok" binding:"required"`
}

type DTOCustomerSimpan struct {
	CustomerCode string `json:"customerCode" binding:"required"`
	Nama         string `json:"nama" binding:"required"`
	Kontak       string `json:"kontak"`
}

type DTOBOMSimpan struct {
	IDParentMaterial            *uint   `json:"idParentMaterial"`
	IDRawMaterial               uint    `json:"idRawMaterial" binding:"required"`
	ParameterKuantitasPembentuk float64 `json:"parameterKuantitasPembentuk" binding:"required"`
}

// helper to invalidate caches when master data changes
func (l *layananManufaktur) invalidateCache() {
	cache.GlobalCache.Clear()
}

// === LAYANAN PEMASOK ===

func (l *layananManufaktur) TambahPemasok(dto *DTOPemasokSimpan) (*Pemasok, error) {
	p := &Pemasok{
		SupplierCode: dto.SupplierCode,
		NamaEntitas:  dto.NamaEntitas,
		Kontak:       dto.Kontak,
	}
	if err := l.repo.BuatPemasok(p); err != nil {
		return nil, err
	}
	l.invalidateCache()
	return p, nil
}

func (l *layananManufaktur) CariPemasokID(id uint) (*Pemasok, error) {
	return l.repo.CariPemasokBerdasarkanID(id)
}

func (l *layananManufaktur) AmbilSemuaPemasok() ([]Pemasok, error) {
	cacheKey := "manufaktur_semua_pemasok"
	if cachedVal, found := cache.GlobalCache.Get(cacheKey); found {
		return cachedVal.([]Pemasok), nil
	}
	res, err := l.repo.DaftarPemasok()
	if err != nil {
		return nil, err
	}
	cache.GlobalCache.Set(cacheKey, res, 15*time.Minute)
	return res, nil
}

func (l *layananManufaktur) UbahPemasok(id uint, dto *DTOPemasokSimpan) (*Pemasok, error) {
	p, err := l.repo.CariPemasokBerdasarkanID(id)
	if err != nil {
		return nil, errors.New("pemasok_tidak_ditemukan")
	}
	p.SupplierCode = dto.SupplierCode
	p.NamaEntitas = dto.NamaEntitas
	p.Kontak = dto.Kontak

	if err := l.repo.PerbaruiPemasok(p); err != nil {
		return nil, err
	}
	l.invalidateCache()
	return p, nil
}

func (l *layananManufaktur) HapusPemasokID(id uint) error {
	_, err := l.repo.CariPemasokBerdasarkanID(id)
	if err != nil {
		return errors.New("pemasok_tidak_ditemukan")
	}
	if err := l.repo.HapusPemasok(id); err != nil {
		return err
	}
	l.invalidateCache()
	return nil
}

// === LAYANAN MATERIAL ===

func (l *layananManufaktur) TambahMaterial(dto *DTOMaterialSimpan) (*Material, error) {
	// Verifikasi pemasok ada
	_, errPemasok := l.repo.CariPemasokBerdasarkanID(dto.IDPemasok)
	if errPemasok != nil {
		return nil, errors.New("pemasok_tidak_ditemukan")
	}

	m := &Material{
		KodeSKU:      dto.KodeSKU,
		NamaMaterial: dto.NamaMaterial,
		TebalCM:      dto.TebalCM,
		BeratGSM:     dto.BeratGSM,
		LebarCM:      dto.LebarCM,
		PanjangCM:    dto.PanjangCM,
		UnitSatuan:   dto.UnitSatuan,
		IDPemasok:    dto.IDPemasok,
	}

	if err := l.repo.BuatMaterial(m); err != nil {
		return nil, err
	}
	l.invalidateCache()
	return m, nil
}

func (l *layananManufaktur) CariMaterialID(id uint) (*Material, error) {
	return l.repo.CariMaterialBerdasarkanID(id)
}

func (l *layananManufaktur) AmbilSemuaMaterial() ([]Material, error) {
	cacheKey := "manufaktur_semua_material"
	if cachedVal, found := cache.GlobalCache.Get(cacheKey); found {
		return cachedVal.([]Material), nil
	}
	res, err := l.repo.DaftarMaterial()
	if err != nil {
		return nil, err
	}
	cache.GlobalCache.Set(cacheKey, res, 15*time.Minute)
	return res, nil
}

func (l *layananManufaktur) UbahMaterial(id uint, dto *DTOMaterialSimpan) (*Material, error) {
	m, err := l.repo.CariMaterialBerdasarkanID(id)
	if err != nil {
		return nil, errors.New("material_tidak_ditemukan")
	}

	_, errPemasok := l.repo.CariPemasokBerdasarkanID(dto.IDPemasok)
	if errPemasok != nil {
		return nil, errors.New("pemasok_tidak_ditemukan")
	}

	m.KodeSKU = dto.KodeSKU
	m.NamaMaterial = dto.NamaMaterial
	m.TebalCM = dto.TebalCM
	m.BeratGSM = dto.BeratGSM
	m.LebarCM = dto.LebarCM
	m.PanjangCM = dto.PanjangCM
	m.UnitSatuan = dto.UnitSatuan
	m.IDPemasok = dto.IDPemasok

	if err := l.repo.PerbaruiMaterial(m); err != nil {
		return nil, err
	}
	l.invalidateCache()
	return m, nil
}

func (l *layananManufaktur) HapusMaterialID(id uint) error {
	_, err := l.repo.CariMaterialBerdasarkanID(id)
	if err != nil {
		return errors.New("material_tidak_ditemukan")
	}
	if err := l.repo.HapusMaterial(id); err != nil {
		return err
	}
	l.invalidateCache()
	return nil
}

// === LAYANAN CUSTOMER ===

func (l *layananManufaktur) TambahCustomer(dto *DTOCustomerSimpan) (*Customer, error) {
	c := &Customer{
		CustomerCode: dto.CustomerCode,
		Nama:         dto.Nama,
		Kontak:       dto.Kontak,
	}
	if err := l.repo.BuatCustomer(c); err != nil {
		return nil, err
	}
	l.invalidateCache()
	return c, nil
}

func (l *layananManufaktur) CariCustomerID(id uint) (*Customer, error) {
	return l.repo.CariCustomerBerdasarkanID(id)
}

func (l *layananManufaktur) AmbilSemuaCustomer() ([]Customer, error) {
	cacheKey := "manufaktur_semua_customer"
	if cachedVal, found := cache.GlobalCache.Get(cacheKey); found {
		return cachedVal.([]Customer), nil
	}
	res, err := l.repo.DaftarCustomer()
	if err != nil {
		return nil, err
	}
	cache.GlobalCache.Set(cacheKey, res, 15*time.Minute)
	return res, nil
}

func (l *layananManufaktur) UbahCustomer(id uint, dto *DTOCustomerSimpan) (*Customer, error) {
	c, err := l.repo.CariCustomerBerdasarkanID(id)
	if err != nil {
		return nil, errors.New("customer_tidak_ditemukan")
	}
	c.CustomerCode = dto.CustomerCode
	c.Nama = dto.Nama
	c.Kontak = dto.Kontak

	if err := l.repo.PerbaruiCustomer(c); err != nil {
		return nil, err
	}
	l.invalidateCache()
	return c, nil
}

func (l *layananManufaktur) HapusCustomerID(id uint) error {
	_, err := l.repo.CariCustomerBerdasarkanID(id)
	if err != nil {
		return errors.New("customer_tidak_ditemukan")
	}
	if err := l.repo.HapusCustomer(id); err != nil {
		return err
	}
	l.invalidateCache()
	return nil
}

// === LAYANAN BOM ===

func (l *layananManufaktur) TambahBOM(dto *DTOBOMSimpan) (*KomposisiMaterialBOM, error) {
	// Verifikasi raw material ada
	_, errRaw := l.repo.CariMaterialBerdasarkanID(dto.IDRawMaterial)
	if errRaw != nil {
		return nil, errors.New("material_baku_tidak_ditemukan")
	}

	// Verifikasi parent material ada jika terdefinisi
	if dto.IDParentMaterial != nil {
		_, errParent := l.repo.CariMaterialBerdasarkanID(*dto.IDParentMaterial)
		if errParent != nil {
			return nil, errors.New("material_induk_tidak_ditemukan")
		}
	}

	b := &KomposisiMaterialBOM{
		IDParentMaterial:            dto.IDParentMaterial,
		IDRawMaterial:               dto.IDRawMaterial,
		ParameterKuantitasPembentuk: dto.ParameterKuantitasPembentuk,
	}

	if err := l.repo.BuatBOM(b); err != nil {
		return nil, err
	}
	l.invalidateCache()
	return b, nil
}

func (l *layananManufaktur) CariBOMID(id uint) (*KomposisiMaterialBOM, error) {
	return l.repo.CariBOMBerdasarkanID(id)
}

func (l *layananManufaktur) AmbilSemuaBOM() ([]KomposisiMaterialBOM, error) {
	cacheKey := "manufaktur_semua_bom"
	if cachedVal, found := cache.GlobalCache.Get(cacheKey); found {
		return cachedVal.([]KomposisiMaterialBOM), nil
	}
	res, err := l.repo.DaftarBOM()
	if err != nil {
		return nil, err
	}
	cache.GlobalCache.Set(cacheKey, res, 15*time.Minute)
	return res, nil
}

func (l *layananManufaktur) UbahBOM(id uint, dto *DTOBOMSimpan) (*KomposisiMaterialBOM, error) {
	b, err := l.repo.CariBOMBerdasarkanID(id)
	if err != nil {
		return nil, errors.New("bom_tidak_ditemukan")
	}

	_, errRaw := l.repo.CariMaterialBerdasarkanID(dto.IDRawMaterial)
	if errRaw != nil {
		return nil, errors.New("material_baku_tidak_ditemukan")
	}

	if dto.IDParentMaterial != nil {
		_, errParent := l.repo.CariMaterialBerdasarkanID(*dto.IDParentMaterial)
		if errParent != nil {
			return nil, errors.New("material_induk_tidak_ditemukan")
		}
	}

	b.IDParentMaterial = dto.IDParentMaterial
	b.IDRawMaterial = dto.IDRawMaterial
	b.ParameterKuantitasPembentuk = dto.ParameterKuantitasPembentuk

	if err := l.repo.PerbaruiBOM(b); err != nil {
		return nil, err
	}
	l.invalidateCache()
	return b, nil
}

func (l *layananManufaktur) HapusBOMID(id uint) error {
	_, err := l.repo.CariBOMBerdasarkanID(id)
	if err != nil {
		return errors.New("bom_tidak_ditemukan")
	}
	if err := l.repo.HapusBOM(id); err != nil {
		return err
	}
	l.invalidateCache()
	return nil
}

func (l *layananManufaktur) AmbilSnapshotMasterData() (*MasterDataSnapshotDto, *MetadataSnapshotDto, error) {
	materialList, err := l.repo.AmbilSnapshotMasterData()
	if err != nil {
		return nil, nil, err
	}

	materialDtos := make([]MaterialSnapshotDto, 0, len(materialList))
	for _, m := range materialList {
		materialDtos = append(materialDtos, MaterialSnapshotDto{
			ID:           strconv.FormatUint(uint64(m.ID), 10),
			KodeSKU:      m.KodeSKU,
			NamaMaterial: m.NamaMaterial,
			Aktif:        true,
		})
	}

	// Inisialisasi slice kosong (bukan nil) untuk menjamin output JSON [] dan menghindari null
	snapshot := &MasterDataSnapshotDto{
		VersiMasterData:  "v1.0.0",
		LineProduksi:     []LineProduksiSnapshotDto{},
		SlotWaktu:        []SlotWaktuSnapshotDto{},
		Material:         materialDtos,
		Part:             []PartSnapshotDto{},
		KategoriDefect:   []KategoriDefectSnapshotDto{},
		JenisDefect:      []JenisDefectSnapshotDto{},
		RelasiPartDefect: []RelasiPartDefectSnapshotDto{},
	}

	metadata := &MetadataSnapshotDto{
		JumlahLineProduksi:     0,
		JumlahSlotWaktu:        0,
		JumlahMaterial:         len(materialDtos),
		JumlahPart:             0,
		JumlahJenisDefect:      0,
		JumlahRelasiPartDefect: 0,
		JumlahShiftOperasional: 3,
	}

	return snapshot, metadata, nil
}
