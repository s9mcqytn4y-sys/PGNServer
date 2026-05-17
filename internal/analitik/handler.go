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

// TanganiParetoBulanan memberikan data histogram Pareto bulanan.
// @Summary Dapatkan Metrik Pareto Bulanan
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
			respon.Galat_Validasi(k, "Parameter bulan harus berupa angka yang valid", nil)
			return
		}
	}

	if tahunStr != "" {
		if t, err := strconv.Atoi(tahunStr); err == nil && t > 2000 {
			tahun = t
		} else {
			respon.Galat_Validasi(k, "Parameter tahun harus berupa angka yang valid", nil)
			return
		}
	}

	hasil, err := h.layanan.DapatkanParetoBulanan(k, bulan, tahun)
	if err != nil {
		respon.Galat_Server(k, "Gagal mengkalkulasi agregat Pareto.", err)
		return
	}

	respon.Sukses(k, "Kalkulasi Pareto berhasil ditarik.", hasil)
}

// TanganiPareto memberikan data histogram Pareto secara dinamis.
// @Summary Dapatkan Pareto Dinamis
// @Description Mengembalikan kalkulasi Pareto 80/20 berdasarkan tanggal mulai, tanggal selesai, dan lini.
// @Tags Analitik
// @Accept json
// @Produce json
// @Param start_date query string false "Tanggal Mulai (YYYY-MM-DD)"
// @Param end_date query string false "Tanggal Selesai (YYYY-MM-DD)"
// @Param line query string false "Nama Lini / Zona"
// @Success 200 {object} respon.ResponStandar
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/analitik/pareto [get]
func (h *PenangananAnalitik) TanganiPareto(k *gin.Context) {
	tanggalMulai := k.Query("start_date")
	tanggalSelesai := k.Query("end_date")
	zonaLini := k.Query("line")

	hasil, err := h.layanan.DapatkanPareto(k, tanggalMulai, tanggalSelesai, zonaLini)
	if err != nil {
		respon.Galat_Server(k, "Gagal mengkalkulasi agregat Pareto.", err)
		return
	}

	respon.Sukses(k, "Kalkulasi Pareto dinamis berhasil ditarik.", hasil)
}

// TanganiLacakAkarMasalah melacak akar masalah material defect dari finished good.
// @Summary Lacak Akar Masalah Defect (BOM Tracing)
// @Description Menelusuri defects produk hingga ke level bahan baku pembentuk & supplier
// @Tags Analitik
// @Accept json
// @Produce json
// @Param kode_cacat query string true "Kode Defect Cacat"
// @Param parent_sku query string false "SKU Produk Jadi (Finished Good)"
// @Success 200 {object} respon.ResponStandar
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/analitik/lacak [get]
func (h *PenangananAnalitik) TanganiLacakAkarMasalah(k *gin.Context) {
	kodeCacat := k.Query("kode_cacat")
	parentSKU := k.Query("parent_sku")

	if kodeCacat == "" {
		respon.Galat_Validasi(k, "Parameter kode_cacat wajib disertakan", nil)
		return
	}

	hasil, err := h.layanan.LacakAkarMasalah(k, kodeCacat, parentSKU)
	if err != nil {
		respon.Galat_Server(k, "Gagal melacak akar masalah BOM.", err)
		return
	}

	respon.Sukses(k, "BOM Tracing akar masalah cacat berhasil dilakukan.", hasil)
}

// TanganiRingkasanNG memberikan ringkasan metrik NG.
// @Summary Dapatkan Ringkasan NG
// @Description Mengembalikan ringkasan produksi dan defect (Total Produksi, OK, Defect, Rasio NG).
// @Tags Analitik
// @Accept json
// @Produce json
// @Param start_date query string false "Tanggal Mulai (YYYY-MM-DD)"
// @Param end_date query string false "Tanggal Selesai (YYYY-MM-DD)"
// @Param line query string false "Nama Lini / Zona"
// @Success 200 {object} respon.ResponStandar{data=DTORingkasanNG}
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/analitik/ringkasan_ng [get]
func (h *PenangananAnalitik) TanganiRingkasanNG(k *gin.Context) {
	tanggalMulai := k.Query("start_date")
	tanggalSelesai := k.Query("end_date")
	zonaLini := k.Query("line")

	hasil, err := h.layanan.DapatkanRingkasanNG(k, tanggalMulai, tanggalSelesai, zonaLini)
	if err != nil {
		respon.Galat_Server(k, "Gagal mendapatkan ringkasan NG.", err)
		return
	}

	respon.Sukses(k, "Ringkasan NG berhasil ditarik.", hasil)
}

// TanganiHistogramDefect memberikan data histogram defect.
// @Summary Dapatkan Histogram Defect
// @Description Mengembalikan frekuensi defect yang dikelompokkan berdasarkan parameter tertentu.
// @Tags Analitik
// @Accept json
// @Produce json
// @Param start_date query string false "Tanggal Mulai (YYYY-MM-DD)"
// @Param end_date query string false "Tanggal Selesai (YYYY-MM-DD)"
// @Param line query string false "Nama Lini / Zona"
// @Param group_by query string false "Grup (waktu | tanggal | kode_cacat)"
// @Success 200 {object} respon.ResponStandar{data=[]DTOHistogramDefect}
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/analitik/histogram_defect [get]
func (h *PenangananAnalitik) TanganiHistogramDefect(k *gin.Context) {
	tanggalMulai := k.Query("start_date")
	tanggalSelesai := k.Query("end_date")
	zonaLini := k.Query("line")
	groupBy := k.Query("group_by")

	hasil, err := h.layanan.DapatkanHistogramDefect(k, tanggalMulai, tanggalSelesai, zonaLini, groupBy)
	if err != nil {
		respon.Galat_Server(k, "Gagal mendapatkan histogram defect.", err)
		return
	}

	respon.Sukses(k, "Histogram defect berhasil ditarik.", hasil)
}

// TanganiTrendDefect memberikan data trend/run chart defect.
// @Summary Dapatkan Trend Defect
// @Description Mengembalikan tren defect berdasarkan periode.
// @Tags Analitik
// @Accept json
// @Produce json
// @Param start_date query string false "Tanggal Mulai (YYYY-MM-DD)"
// @Param end_date query string false "Tanggal Selesai (YYYY-MM-DD)"
// @Param line query string false "Nama Lini / Zona"
// @Param periode query string false "Periode (harian | mingguan | bulanan)"
// @Success 200 {object} respon.ResponStandar{data=[]DTOTrendDefect}
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/analitik/trend_defect [get]
func (h *PenangananAnalitik) TanganiTrendDefect(k *gin.Context) {
	tanggalMulai := k.Query("start_date")
	tanggalSelesai := k.Query("end_date")
	zonaLini := k.Query("line")
	periode := k.Query("periode")

	hasil, err := h.layanan.DapatkanTrendDefect(k, tanggalMulai, tanggalSelesai, zonaLini, periode)
	if err != nil {
		respon.Galat_Server(k, "Gagal mendapatkan trend defect.", err)
		return
	}

	respon.Sukses(k, "Trend defect berhasil ditarik.", hasil)
}

// TanganiStratifikasiDefect memberikan data stratifikasi defect.
// @Summary Dapatkan Stratifikasi Defect
// @Description Mengembalikan stratifikasi defect (contoh: berdasarkan kode_cacat).
// @Tags Analitik
// @Accept json
// @Produce json
// @Param start_date query string false "Tanggal Mulai (YYYY-MM-DD)"
// @Param end_date query string false "Tanggal Selesai (YYYY-MM-DD)"
// @Param line query string false "Nama Lini / Zona"
// @Param kode_cacat query string false "Kode Cacat Spesifik"
// @Success 200 {object} respon.ResponStandar{data=[]DTOStratifikasiDefect}
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/analitik/stratifikasi_defect [get]
func (h *PenangananAnalitik) TanganiStratifikasiDefect(k *gin.Context) {
	tanggalMulai := k.Query("start_date")
	tanggalSelesai := k.Query("end_date")
	zonaLini := k.Query("line")
	kodeCacat := k.Query("kode_cacat")

	hasil, err := h.layanan.DapatkanStratifikasiDefect(k, tanggalMulai, tanggalSelesai, zonaLini, kodeCacat)
	if err != nil {
		respon.Galat_Server(k, "Gagal mendapatkan stratifikasi defect.", err)
		return
	}

	respon.Sukses(k, "Stratifikasi defect berhasil ditarik.", hasil)
}

// TanganiSinyalKualitas memberikan sinyal atau alert kualitas (Control Signal).
// @Summary Dapatkan Sinyal Kualitas
// @Description Mengembalikan sinyal peringatan kualitas (KRITIS, WASPADA, STABIL) berdasarkan batas toleransi.
// @Tags Analitik
// @Accept json
// @Produce json
// @Param start_date query string false "Tanggal Mulai (YYYY-MM-DD)"
// @Param end_date query string false "Tanggal Selesai (YYYY-MM-DD)"
// @Param line query string false "Nama Lini / Zona"
// @Success 200 {object} respon.ResponStandar{data=DTOSinyalKualitas}
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/analitik/sinyal_kualitas [get]
func (h *PenangananAnalitik) TanganiSinyalKualitas(k *gin.Context) {
	tanggalMulai := k.Query("start_date")
	tanggalSelesai := k.Query("end_date")
	zonaLini := k.Query("line")

	hasil, err := h.layanan.DapatkanSinyalKualitas(k, tanggalMulai, tanggalSelesai, zonaLini)
	if err != nil {
		respon.Galat_Server(k, "Gagal mendapatkan sinyal kualitas.", err)
		return
	}

	respon.Sukses(k, "Sinyal kualitas berhasil ditarik.", hasil)
}
