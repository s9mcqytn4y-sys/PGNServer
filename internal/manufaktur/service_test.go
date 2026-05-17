package manufaktur

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepositoriManufaktur adalah tiruan dari RepositoriManufaktur menggunakan testify/mock.
type MockRepositoriManufaktur struct {
	mock.Mock
}

func (m *MockRepositoriManufaktur) BuatPemasok(p *Pemasok) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) CariPemasokBerdasarkanID(id uint) (*Pemasok, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*Pemasok), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriManufaktur) DaftarPemasok() ([]Pemasok, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]Pemasok), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriManufaktur) PerbaruiPemasok(p *Pemasok) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) HapusPemasok(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) BuatMaterial(mat *Material) error {
	args := m.Called(mat)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) CariMaterialBerdasarkanID(id uint) (*Material, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*Material), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriManufaktur) CariMaterialBerdasarkanSKU(sku string) (*Material, error) {
	args := m.Called(sku)
	if args.Get(0) != nil {
		return args.Get(0).(*Material), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriManufaktur) DaftarMaterial() ([]Material, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]Material), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriManufaktur) PerbaruiMaterial(mat *Material) error {
	args := m.Called(mat)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) HapusMaterial(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) BuatCustomer(c *Customer) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) CariCustomerBerdasarkanID(id uint) (*Customer, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriManufaktur) DaftarCustomer() ([]Customer, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriManufaktur) PerbaruiCustomer(c *Customer) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) HapusCustomer(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) BuatBOM(b *KomposisiMaterialBOM) error {
	args := m.Called(b)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) CariBOMBerdasarkanID(id uint) (*KomposisiMaterialBOM, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*KomposisiMaterialBOM), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriManufaktur) DaftarBOM() ([]KomposisiMaterialBOM, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]KomposisiMaterialBOM), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriManufaktur) PerbaruiBOM(b *KomposisiMaterialBOM) error {
	args := m.Called(b)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) HapusBOM(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepositoriManufaktur) AmbilSnapshotMasterData() ([]Material, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]Material), args.Error(1)
	}
	return nil, args.Error(1)
}


func TestTambahPemasok_Sukses(t *testing.T) {
	mockRepo := new(MockRepositoriManufaktur)
	layanan := KonstruksiLayananBaru(mockRepo)

	dto := &DTOPemasokSimpan{
		SupplierCode: "KRS",
		NamaEntitas:  "Krakatau Steel",
		Kontak:       "sales@krakatausteel.co.id",
	}

	mockRepo.On("BuatPemasok", mock.AnythingOfType("*manufaktur.Pemasok")).Return(nil).Once()

	pemasok, err := layanan.TambahPemasok(dto)

	assert.NoError(t, err)
	assert.NotNil(t, pemasok)
	assert.Equal(t, dto.SupplierCode, pemasok.SupplierCode)
	assert.Equal(t, dto.NamaEntitas, pemasok.NamaEntitas)
	mockRepo.AssertExpectations(t)
}

func TestTambahMaterial_PemasokTidakDitemukan(t *testing.T) {
	mockRepo := new(MockRepositoriManufaktur)
	layanan := KonstruksiLayananBaru(mockRepo)

	dto := &DTOMaterialSimpan{
		KodeSKU:      "MAT-FLT-01",
		NamaMaterial: "Filter Element",
		UnitSatuan:   "pcs",
		IDPemasok:    99,
	}

	// Pemasok ID 99 tidak ada
	mockRepo.On("CariPemasokBerdasarkanID", uint(99)).Return((*Pemasok)(nil), errors.New("pemasok_tidak_ditemukan")).Once()

	material, err := layanan.TambahMaterial(dto)

	assert.Error(t, err)
	assert.Nil(t, material)
	assert.Equal(t, "pemasok_tidak_ditemukan", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTambahMaterial_Sukses(t *testing.T) {
	mockRepo := new(MockRepositoriManufaktur)
	layanan := KonstruksiLayananBaru(mockRepo)

	dto := &DTOMaterialSimpan{
		KodeSKU:      "MAT-FLT-02",
		NamaMaterial: "Filter Element A",
		UnitSatuan:   "pcs",
		IDPemasok:    1,
	}

	pemasokMock := &Pemasok{
		SupplierCode: "KRS",
		NamaEntitas:  "Krakatau Steel",
	}
	pemasokMock.ID = 1

	mockRepo.On("CariPemasokBerdasarkanID", uint(1)).Return(pemasokMock, nil).Once()
	mockRepo.On("BuatMaterial", mock.AnythingOfType("*manufaktur.Material")).Return(nil).Once()

	material, err := layanan.TambahMaterial(dto)

	assert.NoError(t, err)
	assert.NotNil(t, material)
	assert.Equal(t, dto.KodeSKU, material.KodeSKU)
	assert.Equal(t, dto.IDPemasok, material.IDPemasok)
	mockRepo.AssertExpectations(t)
}

func TestAmbilSnapshotMasterData_Sukses(t *testing.T) {
	mockRepo := new(MockRepositoriManufaktur)
	layanan := KonstruksiLayananBaru(mockRepo)

	materialsMock := []Material{
		{
			KodeSKU:      "MAT-FLT-01",
			NamaMaterial: "Filter Element",
		},
	}
	materialsMock[0].ID = 123

	mockRepo.On("AmbilSnapshotMasterData").Return(materialsMock, nil).Once()

	snapshot, metadata, err := layanan.AmbilSnapshotMasterData()

	assert.NoError(t, err)
	assert.NotNil(t, snapshot)
	assert.NotNil(t, metadata)
	assert.Equal(t, "v1.0.0", snapshot.VersiMasterData)
	assert.Len(t, snapshot.Material, 1)
	assert.Equal(t, "123", snapshot.Material[0].ID)
	assert.Equal(t, "MAT-FLT-01", snapshot.Material[0].KodeSKU)
	assert.Equal(t, "Filter Element", snapshot.Material[0].NamaMaterial)
	assert.Equal(t, 1, metadata.JumlahMaterial)
	assert.Equal(t, 3, metadata.JumlahShiftOperasional)
	mockRepo.AssertExpectations(t)
}

