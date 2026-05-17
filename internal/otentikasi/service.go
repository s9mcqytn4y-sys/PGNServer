package otentikasi

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LayananOtentikasi interface {
	Registrasi(surel string, sandi string, peran string) (*Pengguna, error)
	Login(surel string, sandi string) (string, error)
	LupaSandi(surel string, sandiBaru string) error
	AmbilProfilBerdasarkanID(id uint) (*Pengguna, error)
}

type layananOtentikasi struct {
	repositori RepositoriOtentikasi
}

func KonstruksiLayananBaru(repo RepositoriOtentikasi) LayananOtentikasi {
	return &layananOtentikasi{repositori: repo}
}

func (l *layananOtentikasi) Registrasi(surel string, sandi string, peran string) (*Pengguna, error) {
	// Mitigasi Happy Path terhadap registrasi ganda
	penggunaEksisting, _ := l.repositori.CariBerdasarkanSurel(surel)
	if penggunaEksisting != nil {
		return nil, errors.New("surel_telah_terdaftar")
	}

	sandiEnkripsi, errEnkripsi := bcrypt.GenerateFromPassword([]byte(sandi), bcrypt.DefaultCost)
	if errEnkripsi != nil {
		return nil, errEnkripsi
	}

	penggunaBaru := &Pengguna{
		SurelKredensial:       surel,
		KataSandiTerenskripsi: string(sandiEnkripsi),
		PeranOtorisasi:        peran,
	}

	errSimpan := l.repositori.Daftar(penggunaBaru)
	if errSimpan != nil {
		return nil, errSimpan
	}

	return penggunaBaru, nil
}

func (l *layananOtentikasi) Login(surel string, sandi string) (string, error) {
	pengguna, errCari := l.repositori.CariBerdasarkanSurel(surel)
	if errCari != nil {
		return "", errors.New("kredensial_tidak_valid")
	}

	errValidasi := bcrypt.CompareHashAndPassword([]byte(pengguna.KataSandiTerenskripsi), []byte(sandi))
	if errValidasi != nil {
		return "", errors.New("kredensial_tidak_valid")
	}

	// Pembangkitan JWT v5
	kunciRahasia := os.Getenv("JWT_SECRET")
	if kunciRahasia == "" {
		kunciRahasia = "rahasia-default"
	}

	klaim := jwt.MapClaims{
		"id":    pengguna.ID,
		"surel": pengguna.SurelKredensial,
		"peran": pengguna.PeranOtorisasi,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, klaim)
	tokenBerbasisString, errToken := token.SignedString([]byte(kunciRahasia))
	if errToken != nil {
		return "", errToken
	}

	return tokenBerbasisString, nil
}

func (l *layananOtentikasi) LupaSandi(surel string, sandiBaru string) error {
	pengguna, errCari := l.repositori.CariBerdasarkanSurel(surel)
	if errCari != nil {
		return errors.New("akun_tidak_ditemukan")
	}

	sandiEnkripsi, errEnkripsi := bcrypt.GenerateFromPassword([]byte(sandiBaru), bcrypt.DefaultCost)
	if errEnkripsi != nil {
		return errEnkripsi
	}

	pengguna.KataSandiTerenskripsi = string(sandiEnkripsi)
	return l.repositori.PerbaruiSandi(pengguna)
}

func (l *layananOtentikasi) AmbilProfilBerdasarkanID(id uint) (*Pengguna, error) {
	return l.repositori.CariBerdasarkanID(id)
}
