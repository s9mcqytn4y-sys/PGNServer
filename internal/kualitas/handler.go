package kualitas

import (
	"github.com/gin-gonic/gin"
	"pgn-server/pkg/respon"
)

// PenangananKualitas menjadi garda depan validasi pintu masuk.
type PenangananKualitas struct {
	layanan LayananKualitas
}

func KonstruksiPenangananBaru(layanan LayananKualitas) *PenangananKualitas {
	return &PenangananKualitas{layanan: layanan}
}

// TanganiRekamLembarPeriksa menerima permintaan pencatatan dari ujung gerbang API.
func (p *PenangananKualitas) TanganiRekamLembarPeriksa(k *gin.Context) {
	var dto DTO_LembarPeriksa_Kirim

	// Tangkap eksepsi bila permohonan antarmuka terdistorsi
	if err := k.ShouldBindJSON(&dto); err != nil {
		respon.Galat_Validasi(k, "Struktur laporan inspeksi cacat tidak lengkap: "+err.Error())
		return
	}

	errProses := p.layanan.RekamLembarPeriksa(&dto)
	if errProses != nil {
		respon.Galat_Server(k, "Gagal mencatat transmisi himpunan Lembar Periksa ke pangkalan data.")
		return
	}

	respon.Sukses(k, "Data lembar periksa harian berhasil direkam.", nil)
}
