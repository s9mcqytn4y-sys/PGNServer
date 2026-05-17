package kualitas

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MockRepositoriKualitas adalah tiruan dari RepositoriKualitas menggunakan testify/mock.
type MockRepositoriKualitas struct {
	mock.Mock
}

func (m *MockRepositoriKualitas) SimpanLembarPeriksaMassal(dto *DTOLembarPeriksaKirim, tx *gorm.DB) error {
	args := m.Called(dto, tx)
	return args.Error(0)
}

func (m *MockRepositoriKualitas) DaftarRiwayat(limit int, offset int, tanggalMulai string, tanggalSelesai string, zonaLini string) ([]LembarPeriksa, error) {
	args := m.Called(limit, offset, tanggalMulai, tanggalSelesai, zonaLini)
	if args.Get(0) != nil {
		return args.Get(0).([]LembarPeriksa), args.Error(1)
	}
	return nil, args.Error(1)
}

// dapatkanDBTestOffline membuat instance *gorm.DB tiruan (sqlmock) tanpa koneksi fisik ke PostgreSQL.
func dapatkanDBTestOffline() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return gormDB, mock
}

func TestRekamLembarPeriksa_ValidasiTPS_Gagal(t *testing.T) {
	mockRepo := new(MockRepositoriKualitas)
	dbOffline, sMock := dapatkanDBTestOffline()
	layanan := KonstruksiLayananBaru(mockRepo, dbOffline)

	// Skenario: Total Produksi (100) != OK (80) + NG (10) -> Selisih 10 (Gagal)
	dtoTpsGagal := &DTOLembarPeriksaKirim{
		Tanggal:            "2026-05-17",
		ZonaLini:           "Lini A",
		PenggunaIDTercatat: 1,
		Detail: []DTODetailInspeksi{
			{
				UnikPartID:      1,
				KodeCacat:       "LAMINATING BOLONG",
				WaktuPergeseran: "Shift-1",
				TotalProduksi:   100,
				RasioTotalOK:    80,
				RasioCacat:      10, 
			},
		},
	}

	err := layanan.RekamLembarPeriksa(dtoTpsGagal)

	// Validasi gerbang TPS wajib gagal & return error spesifik sebelum menyentuh DB/Transaksi
	assert.Error(t, err)
	assert.Equal(t, "validasi_tps_gagal: total produksi harus sama dengan jumlah OK dan NG", err.Error())
	mockRepo.AssertNotCalled(t, "SimpanLembarPeriksaMassal", mock.Anything, mock.Anything)
	assert.NoError(t, sMock.ExpectationsWereMet())
}

func TestRekamLembarPeriksa_ValidasiTPS_Sukses(t *testing.T) {
	mockRepo := new(MockRepositoriKualitas)
	dbOffline, sMock := dapatkanDBTestOffline()
	layanan := KonstruksiLayananBaru(mockRepo, dbOffline)

	// Skenario: Total Produksi (100) == OK (95) + NG (5) -> Cocok (Sukses)
	dtoTpsSukses := &DTOLembarPeriksaKirim{
		Tanggal:            "2026-05-17",
		ZonaLini:           "Lini A",
		PenggunaIDTercatat: 1,
		Detail: []DTODetailInspeksi{
			{
				UnikPartID:      1,
				KodeCacat:       "LAMINATING BOLONG",
				WaktuPergeseran: "Shift-1",
				TotalProduksi:   100,
				RasioTotalOK:    100,
				RasioCacat:      0,
			},
		},
	}

	// Persiapkan ekspektasi SQL mock untuk transaksi
	sMock.ExpectBegin()
	sMock.ExpectCommit()

	// Atur mock agar SimpanLembarPeriksaMassal sukses
	mockRepo.On("SimpanLembarPeriksaMassal", dtoTpsSukses, mock.AnythingOfType("*gorm.DB")).Return(nil).Once()

	err := layanan.RekamLembarPeriksa(dtoTpsSukses)

	// Harusnya lolos validasi TPS dan mengeksekusi penyimpanan
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	assert.NoError(t, sMock.ExpectationsWereMet())
}

func TestRekamLembarPeriksa_SimpanGagal_Rollback(t *testing.T) {
	mockRepo := new(MockRepositoriKualitas)
	dbOffline, sMock := dapatkanDBTestOffline()
	layanan := KonstruksiLayananBaru(mockRepo, dbOffline)

	dtoTpsSukses := &DTOLembarPeriksaKirim{
		Tanggal:            "2026-05-17",
		ZonaLini:           "Lini A",
		PenggunaIDTercatat: 1,
		Detail: []DTODetailInspeksi{
			{
				UnikPartID:      1,
				KodeCacat:       "LAMINATING BOLONG",
				WaktuPergeseran: "Shift-1",
				TotalProduksi:   100,
				RasioTotalOK:    95,
				RasioCacat:      5,
			},
		},
	}

	// Persiapkan ekspektasi SQL mock untuk transaksi
	sMock.ExpectBegin()
	sMock.ExpectRollback()

	// Atur mock agar SimpanLembarPeriksaMassal mengembalikan error database
	errDB := errors.New("db_connection_error_or_constraint_violation")
	mockRepo.On("SimpanLembarPeriksaMassal", dtoTpsSukses, mock.AnythingOfType("*gorm.DB")).Return(errDB).Once()

	err := layanan.RekamLembarPeriksa(dtoTpsSukses)

	// Harus mengembalikan error dan rollback dieksekusi secara aman
	assert.Error(t, err)
	assert.Equal(t, errDB, err)
	mockRepo.AssertExpectations(t)
	assert.NoError(t, sMock.ExpectationsWereMet())
}
