// Package analitik menangani agregasi dan kalkulasi data pelaporan 7 QC Tools.
package analitik

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"pgn-server/pkg/respon"
)

// PenangananAnalitik mengurus rute dan kontrol lalu lintas data analitik.
type PenangananAnalitik struct {
	layanan LayananAnalitik
}

// KonstruksiPenangananBaru membuat router pengarah untuk analitik.
func KonstruksiPenangananBaru(layanan LayananAnalitik) *PenangananAnalitik {
	return &PenangananAnalitik{layanan: layanan}
}

// TanganiParetoBulanan memberikan data histogram Pareto.
// @Summary Dapatkan Metrik Pareto
// @Description Mengembalikan kalkulasi Pareto 80/20 per bulan menggunakan Window Function SQL
// @Tags Analitik
// @Accept json
// @Produce json
// @Param bulan query int false "Bulan (1-12)"
// @Param tahun query int false "Tahun (contoh: 2026)"
// @Success 200 {object} respon.ResponStandar
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/analitik/metrik_pareto_bulanan [get]
func (h *PenangananAnalitik) TanganiParetoBulanan(k *gin.Context) {
	bulanStr := k.Query("bulan")
	tahunStr := k.Query("tahun")

	sekarang := time.Now()
	bulan := int(sekarang.Month())
	tahun := sekarang.Year()

	if bulanStr != "" {
		if b, err := strconv.Atoi(bulanStr); err == nil && b >= 1 && b <= 12 {
			bulan = b
		} else {
			respon.Galat_Validasi(k, "Parameter bulan harus berupa angka yang valid")
			return
		}
	}

	if tahunStr != "" {
		if t, err := strconv.Atoi(tahunStr); err == nil && t > 2000 {
			tahun = t
		} else {
			respon.Galat_Validasi(k, "Parameter tahun harus berupa angka yang valid")
			return
		}
	}

	hasil, err := h.layanan.DapatkanParetoBulanan(bulan, tahun)
	if err != nil {
		respon.Galat_Server(k, "Gagal mengkalkulasi agregat Pareto.")
		return
	}

	respon.Sukses(k, "Kalkulasi Pareto berhasil ditarik.", hasil)
}
