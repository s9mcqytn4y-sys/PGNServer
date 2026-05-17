package media

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/textproto"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MockRepositoriMedia is a test double for RepositoriMedia
type MockRepositoriMedia struct {
	mock.Mock
}

func (m *MockRepositoriMedia) Simpan(aset *AsetDigital) error {
	args := m.Called(aset)
	return args.Error(0)
}

func (m *MockRepositoriMedia) HapusBerdasarkanID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepositoriMedia) CariBerdasarkanID(id uint) (*AsetDigital, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*AsetDigital), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepositoriMedia) DapatkanBerdasarkanIDMaterial(idMaterial uint) ([]*AsetDigital, error) {
	args := m.Called(idMaterial)
	if args.Get(0) != nil {
		return args.Get(0).([]*AsetDigital), args.Error(1)
	}
	return nil, args.Error(1)
}

// dapatkanDBTest membuat instance gorm DB mock
func dapatkanDBTest() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return gormDB, mock
}

// Helper to create a fake multipart.FileHeader
func buatFileHeaderTiruan(nama string, tipeMime string, data []byte) *multipart.FileHeader {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("berkas", nama)
	if err != nil {
		panic(err)
	}
	_, _ = part.Write(data)
	_ = writer.Close()

	// Parse it back to multipart.Reader to get a FileHeader
	boundary := writer.Boundary()
	reader := multipart.NewReader(bytes.NewReader(buf.Bytes()), boundary)
	form, err := reader.ReadForm(10 * 1024 * 1024)
	if err != nil {
		panic(err)
	}

	files := form.File["berkas"]
	if len(files) == 0 {
		// Fallback header manually if parser fails
		return &multipart.FileHeader{
			Filename: nama,
			Size:     int64(len(data)),
			Header: textproto.MIMEHeader{
				"Content-Type": []string{tipeMime},
			},
		}
	}

	// Make sure the header has correct content type
	files[0].Header.Set("Content-Type", tipeMime)
	return files[0]
}

func TestUnggahBerkasLokal_MaterialTidakDitemukan(t *testing.T) {
	mockRepo := new(MockRepositoriMedia)
	db, sMock := dapatkanDBTest()
	layanan := KonstruksiLayananBaru(mockRepo, db)

	// Material ID = 99
	sMock.ExpectQuery(`SELECT \* FROM "MATERIAL" WHERE "MATERIAL"\."id" = \$1`).
		WithArgs(99, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	berkas := buatFileHeaderTiruan("test.png", "image/png", []byte("fake_image_bytes"))

	aset, err := layanan.UnggahBerkasLokal(99, berkas)

	assert.Error(t, err)
	assert.Nil(t, aset)
	assert.Equal(t, "referensi_material_tidak_ditemukan", err.Error())
	assert.NoError(t, sMock.ExpectationsWereMet())
}

func TestUnggahBerkasLokal_UkuranTerlaluBesar(t *testing.T) {
	mockRepo := new(MockRepositoriMedia)
	db, sMock := dapatkanDBTest()
	layanan := KonstruksiLayananBaru(mockRepo, db)

	// mock database first
	sMock.ExpectQuery(`SELECT \* FROM "MATERIAL" WHERE "MATERIAL"\."id" = \$1`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode_sku"}).AddRow(1, "MAT-SKU-1"))

	// Create a large file (> 5MB)
	besarData := 6 * 1024 * 1024 // 6 MB
	data := make([]byte, besarData)
	berkas := buatFileHeaderTiruan("test.png", "image/png", data)

	aset, err := layanan.UnggahBerkasLokal(1, berkas)

	assert.Error(t, err)
	assert.Nil(t, aset)
	assert.Equal(t, "ukuran_berkas_terlalu_besar", err.Error())
	assert.NoError(t, sMock.ExpectationsWereMet())
}

func TestUnggahBerkasLokal_FormatTidakDiizinkan(t *testing.T) {
	mockRepo := new(MockRepositoriMedia)
	db, sMock := dapatkanDBTest()
	layanan := KonstruksiLayananBaru(mockRepo, db)

	sMock.ExpectQuery(`SELECT \* FROM "MATERIAL" WHERE "MATERIAL"\."id" = \$1`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode_sku"}).AddRow(1, "MAT-SKU-1"))

	// Format pdf is not allowed
	berkas := buatFileHeaderTiruan("document.pdf", "application/pdf", []byte("pdf_bytes"))

	aset, err := layanan.UnggahBerkasLokal(1, berkas)

	assert.Error(t, err)
	assert.Nil(t, aset)
	assert.Equal(t, "ekstensi_berkas_tidak_diizinkan", err.Error())
	assert.NoError(t, sMock.ExpectationsWereMet())
}

func TestUnggahBerkasLokal_Sukses(t *testing.T) {
	mockRepo := new(MockRepositoriMedia)
	db, sMock := dapatkanDBTest()
	layanan := KonstruksiLayananBaru(mockRepo, db)

	sMock.ExpectQuery(`SELECT \* FROM "MATERIAL" WHERE "MATERIAL"\."id" = \$1`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode_sku"}).AddRow(1, "MAT-SKU-1"))

	// Use real signature of PNG to bypass DetectContentType
	pngSignature := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01")
	berkas := buatFileHeaderTiruan("part_ok.png", "image/png", pngSignature)

	mockRepo.On("Simpan", mock.AnythingOfType("*media.AsetDigital")).Return(nil).Once()

	aset, err := layanan.UnggahBerkasLokal(1, berkas)

	assert.NoError(t, err)
	assert.NotNil(t, aset)
	assert.Equal(t, uint(1), aset.IDMaterial)
	assert.Equal(t, "image/png", aset.TipeMIME)
	assert.Equal(t, ".png", aset.Ekstensi)
	assert.Equal(t, "LOKAL", aset.TipePenyimpanan)

	// Clean up actual created file
	if aset != nil && aset.DirektoriLokal != "" {
		_ = os.Remove(aset.DirektoriLokal)
	}

	mockRepo.AssertExpectations(t)
	assert.NoError(t, sMock.ExpectationsWereMet())
}

func TestHapusBerkas_ProteksiDefault(t *testing.T) {
	mockRepo := new(MockRepositoriMedia)
	db, _ := dapatkanDBTest()
	layanan := KonstruksiLayananBaru(mockRepo, db)

	// Mock existing default asset
	defaultAsset := &AsetDigital{
		ID:              10,
		IDMaterial:      0,
		TipePenyimpanan: "LOKAL",
		DirektoriLokal:  "./penyimpanan/profiles/avatar.png",
	}

	mockRepo.On("CariBerdasarkanID", uint(10)).Return(defaultAsset, nil).Once()
	mockRepo.On("HapusBerdasarkanID", uint(10)).Return(nil).Once()

	// Try removing. It should delete from DB but NOT from filesystem since it's avatar.png
	err := layanan.HapusBerkas(10)

	assert.NoError(t, err)
	// Assert filesystem still has avatar.png
	_, statErr := os.Stat("./penyimpanan/profiles/avatar.png")
	assert.NoError(t, statErr)

	mockRepo.AssertExpectations(t)
}

func TestCariAsetMedia_FallbackDefault(t *testing.T) {
	mockRepo := new(MockRepositoriMedia)
	db, _ := dapatkanDBTest()
	layanan := KonstruksiLayananBaru(mockRepo, db)

	// Mock not found in DB
	mockRepo.On("CariBerdasarkanID", uint(999)).Return((*AsetDigital)(nil), errors.New("not_found")).Once()

	aset, err := layanan.CariAsetMedia(999)

	assert.NoError(t, err)
	assert.NotNil(t, aset)
	assert.Contains(t, aset.DirektoriLokal, "part.png")
	assert.Equal(t, "image/png", aset.TipeMIME)

	mockRepo.AssertExpectations(t)
}
