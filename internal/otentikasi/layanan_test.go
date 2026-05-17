package otentikasi

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// RepositoriMock merupakan tiruan dari RepositoriOtentikasi untuk keperluan unit testing.
type RepositoriMock struct {
	mock.Mock
}

func (m *RepositoriMock) Daftar(pengguna *Pengguna) error {
	args := m.Called(pengguna)
	return args.Error(0)
}

func (m *RepositoriMock) CariBerdasarkanSurel(surel string) (*Pengguna, error) {
	args := m.Called(surel)
	if args.Get(0) != nil {
		return args.Get(0).(*Pengguna), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *RepositoriMock) CariBerdasarkanID(id uint) (*Pengguna, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*Pengguna), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *RepositoriMock) PerbaruiSandi(pengguna *Pengguna) error {
	args := m.Called(pengguna)
	return args.Error(0)
}

func (m *RepositoriMock) PerbaruiTenggatSesi(id uint, tenggatSesi *gorm.DeletedAt) error {
	args := m.Called(id, tenggatSesi)
	return args.Error(0)
}

func TestRegistrasi(t *testing.T) {
	repo := new(RepositoriMock)
	layanan := KonstruksiLayananBaru(repo)

	t.Run("Sukses", func(t *testing.T) {
		repo.On("CariBerdasarkanSurel", "tester@pgn.com").Return(nil, errors.New("tidak_ditemukan")).Once()
		repo.On("Daftar", mock.AnythingOfType("*otentikasi.Pengguna")).Return(nil).Once()

		pengguna, err := layanan.Registrasi("tester@pgn.com", "sandi12345", "Leader")

		assert.NoError(t, err)
		assert.NotNil(t, pengguna)
		assert.Equal(t, "tester@pgn.com", pengguna.SurelKredensial)
		assert.Equal(t, "Leader", pengguna.PeranOtorisasi)
		repo.AssertExpectations(t)
	})

	t.Run("Gagal_SurelSudahAda", func(t *testing.T) {
		repo.On("CariBerdasarkanSurel", "tester@pgn.com").Return(&Pengguna{SurelKredensial: "tester@pgn.com"}, nil).Once()

		pengguna, err := layanan.Registrasi("tester@pgn.com", "sandi12345", "Leader")

		assert.Error(t, err)
		assert.Equal(t, "surel_telah_terdaftar", err.Error())
		assert.Nil(t, pengguna)
		repo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	repo := new(RepositoriMock)
	layanan := KonstruksiLayananBaru(repo)

	os.Setenv("JWT_SECRET", "rahasia-test")

	t.Run("Gagal_KredensialSalah", func(t *testing.T) {
		repo.On("CariBerdasarkanSurel", "tester@pgn.com").Return(nil, errors.New("tidak_ditemukan")).Once()

		token, err := layanan.Login("tester@pgn.com", "sandi12345")

		assert.Error(t, err)
		assert.Equal(t, "kredensial_tidak_valid", err.Error())
		assert.Empty(t, token)
		repo.AssertExpectations(t)
	})
}

func TestAmbilProfilBerdasarkanID(t *testing.T) {
	repo := new(RepositoriMock)
	layanan := KonstruksiLayananBaru(repo)

	t.Run("Sukses", func(t *testing.T) {
		dummyUser := &Pengguna{
			ID:              12,
			SurelKredensial: "tester@pgn.com",
			PeranOtorisasi:  "LEADER",
		}
		repo.On("CariBerdasarkanID", uint(12)).Return(dummyUser, nil).Once()

		pengguna, err := layanan.AmbilProfilBerdasarkanID(12)

		assert.NoError(t, err)
		assert.NotNil(t, pengguna)
		assert.Equal(t, uint(12), pengguna.ID)
		assert.Equal(t, "tester@pgn.com", pengguna.SurelKredensial)
		assert.Equal(t, "LEADER", pengguna.PeranOtorisasi)
		repo.AssertExpectations(t)
	})

	t.Run("Gagal_TidakDitemukan", func(t *testing.T) {
		repo.On("CariBerdasarkanID", uint(99)).Return(nil, errors.New("record not found")).Once()

		pengguna, err := layanan.AmbilProfilBerdasarkanID(99)

		assert.Error(t, err)
		assert.Nil(t, pengguna)
		repo.AssertExpectations(t)
	})
}
