// Package kualitas menangani modul pencatatan inspeksi kontrol kualitas.
package kualitas

import (
	"strconv"

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

// TanganiDaftarRiwayat menerima permintaan daftar historis.
// @Summary Riwayat Lembar Periksa
// @Description Menampilkan daftar historis lembar periksa dengan filter kalender
// @Tags Kualitas
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Batas jumlah data (default 10)"
// @Param offset query int false "Offset paginasi (default 0)"
// @Param tanggal_mulai query string false "Filter tanggal mulai (YYYY-MM-DD)"
// @Param tanggal_selesai query string false "Filter tanggal selesai (YYYY-MM-DD)"
// @Param zona_lini query string false "Filter zona lini"
// @Success 200 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/operasi/riwayat_lembar_periksa [get]
func (p *PenangananKualitas) TanganiDaftarRiwayat(k *gin.Context) {
	limitStr := k.DefaultQuery("limit", "10")
	offsetStr := k.DefaultQuery("offset", "0")
	tanggalMulai := k.Query("tanggal_mulai")
	tanggalSelesai := k.Query("tanggal_selesai")
	zonaLini := k.Query("zona_lini")

	limit, errL := strconv.Atoi(limitStr)
	if errL != nil {
		limit = 10
	}
	offset, errO := strconv.Atoi(offsetStr)
	if errO != nil {
		offset = 0
	}

	riwayat, err := p.layanan.DaftarRiwayat(limit, offset, tanggalMulai, tanggalSelesai, zonaLini)
	if err != nil {
		respon.Galat_Server(k, "Gagal memuat riwayat", err)
		return
	}

	respon.Sukses(k, "Berhasil memuat riwayat lembar periksa.", riwayat)
}
