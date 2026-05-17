// Package kualitas menangani modul pencatatan inspeksi kontrol kualitas.
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
// @Summary Rekam Lembar Periksa
// @Description Menyimpan entri lembar periksa beserta detail inspeksinya
// @Tags Kualitas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body DTOLembarPeriksaKirim true "Payload Lembar Periksa"
// @Success 200 {object} respon.ResponStandar
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/operasi/rekam_lembar_periksa [post]
func (p *PenangananKualitas) TanganiRekamLembarPeriksa(k *gin.Context) {
	var dto DTOLembarPeriksaKirim

	// Tangkap eksepsi bila permohonan antarmuka terdistorsi
	if err := k.ShouldBindJSON(&dto); err != nil {
		respon.Galat_Validasi(k, "Struktur laporan inspeksi cacat tidak lengkap: "+err.Error(), nil)
		return
	}

	errProses := p.layanan.RekamLembarPeriksa(&dto)
	if errProses != nil {
		respon.Galat_Server(k, "Gagal mencatat transmisi himpunan Lembar Periksa ke pangkalan data.", errProses)
		return
	}

	respon.Sukses(k, "Data lembar periksa harian berhasil direkam.", nil)
}
