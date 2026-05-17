package otentikasi

import (
	"pgn-server/pkg/respon"

	"github.com/gin-gonic/gin"
)

type PenangananOtentikasi struct {
	layanan LayananOtentikasi
}

func KonstruksiPenangananBaru(layanan LayananOtentikasi) *PenangananOtentikasi {
	return &PenangananOtentikasi{layanan: layanan}
}

// DataTransferMasuk (DTO) untuk permintaan registrasi
type DataPermintaanRegistrasi struct {
	Surel string `json:"surel" binding:"required,email"`
	Sandi string `json:"sandi" binding:"required,min=8"`
	Peran string `json:"peran"`
}

// TanganiRegistrasi mengatur pendaftaran pengguna baru.
// @Summary Pendaftaran Pengguna
// @Description Mendaftarkan kredensial staf QC/Manajemen baru
// @Tags Otentikasi
// @Accept json
// @Produce json
// @Param body body DataPermintaanRegistrasi true "Data pendaftaran"
// @Success 200 {object} respon.ResponStandar
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/otentikasi/daftar [post]
func (p *PenangananOtentikasi) TanganiRegistrasi(k *gin.Context) {
	var masukan DataPermintaanRegistrasi
	if err := k.ShouldBindJSON(&masukan); err != nil {
		respon.Galat_Validasi(k, "Format surel tidak sah atau sandi kurang dari 8 karakter", nil)
		return
	}

	penggunaBaru, errDaftar := p.layanan.Registrasi(masukan.Surel, masukan.Sandi, masukan.Peran)
	if errDaftar != nil {
		if errDaftar.Error() == "surel_telah_terdaftar" {
			respon.Galat_Validasi(k, "Akun dengan surel tersebut telah terdaftar di sistem", nil)
			return
		}
		respon.Galat_Server(k, "Gagal memproses registrasi akun", errDaftar)
		return
	}

	respon.Sukses(k, "Registrasi akun berhasil", penggunaBaru)
}

// DataTransferMasuk (DTO) untuk permintaan login
type DataPermintaanLogin struct {
	Surel string `json:"surel" binding:"required"`
	Sandi string `json:"sandi" binding:"required"`
}

// TanganiLogin mengatur otentikasi dan mengembalikan JWT.
// @Summary Akses Masuk Pengguna
// @Description Memberikan Token JWT untuk akses sesi internal
// @Tags Otentikasi
// @Accept json
// @Produce json
// @Param body body DataPermintaanLogin true "Kredensial"
// @Success 200 {object} respon.ResponStandar
// @Failure 401 {object} respon.ResponStandar
// @Router /api/v1/otentikasi/masuk [post]
func (p *PenangananOtentikasi) TanganiLogin(k *gin.Context) {
	var masukan DataPermintaanLogin
	if err := k.ShouldBindJSON(&masukan); err != nil {
		respon.Galat_Validasi(k, "Harap berikan surel dan kata sandi", nil)
		return
	}

	token, errLogin := p.layanan.Login(masukan.Surel, masukan.Sandi)
	if errLogin != nil {
		respon.Galat_Validasi(k, "Surel atau kata sandi tidak cocok", nil)
		return
	}

	respon.Sukses(k, "Autentikasi berhasil, token diterbitkan", map[string]string{
		"token": token,
	})
}

// DataTransferMasuk (DTO) untuk permintaan lupa sandi
type DataPermintaanLupaSandi struct {
	Surel     string `json:"surel" binding:"required,email"`
	SandiBaru string `json:"sandiBaru" binding:"required,min=8"`
}

// TanganiLupaSandi mengatur penyetelan ulang sandi.
// @Summary Lupa Sandi
// @Description Pemulihan akun QC
// @Tags Otentikasi
// @Accept json
// @Produce json
// @Param body body DataPermintaanLupaSandi true "Permohonan Reset"
// @Success 200 {object} respon.ResponStandar
// @Router /api/v1/otentikasi/lupa-sandi [post]
func (p *PenangananOtentikasi) TanganiLupaSandi(k *gin.Context) {
	var masukan DataPermintaanLupaSandi
	if err := k.ShouldBindJSON(&masukan); err != nil {
		respon.Galat_Validasi(k, "Pastikan format surel dan kata sandi baru (min 8 karakter) sesuai", nil)
		return
	}

	errLupa := p.layanan.LupaSandi(masukan.Surel, masukan.SandiBaru)
	if errLupa != nil {
		respon.Galat_Validasi(k, "Permintaan gagal diproses. Pastikan akun terdaftar di platform kami.", nil)
		return
	}

	respon.Sukses(k, "Kata sandi berhasil diperbarui", nil)
}

// TanganiLogout mengatur proses keluar pengguna.
// @Summary Akses Keluar
// @Description Mengakhiri sesi pengguna (Client-side token drop)
// @Tags Otentikasi
// @Produce json
// @Success 200 {object} respon.ResponStandar
// @Router /api/v1/otentikasi/keluar [post]
func (p *PenangananOtentikasi) TanganiLogout(k *gin.Context) {
	// Karena menggunakan JWT (stateless), logout dilakukan di klien dengan menghapus token.
	// Kami memberikan konfirmasi sukses.
	respon.Sukses(k, "Berhasil keluar. Silakan hapus token di sisi klien.", nil)
}
