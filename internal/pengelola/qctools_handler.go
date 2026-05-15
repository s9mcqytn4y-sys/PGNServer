package pengelola

import (
	"net/http"
	"pgn-server/internal/konfigurasi"
	"pgn-server/pkg/utilitas"
	"sync"

	"github.com/gin-gonic/gin"
)

type ParetoData struct {
	NamaNG   string `json:"nama_ng"`
	JumlahNG int    `json:"jumlah_ng"`
}

type ControlChartData struct {
	Label string  `json:"label"`
	Nilai float64 `json:"nilai"`
}

// GetParetoData mengembalikan data agregasi cacat tertinggi
// @Summary Ambil Data Pareto
// @Description Mengambil data 80/20 cacat tertinggi menggunakan Goroutine
// @Tags QC-Tools
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utilitas.ResponsAPI
// @Router /qc-tools/pareto [get]
func GetParetoData(c *gin.Context) {
	var results []ParetoData
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// Simulasi query berat atau agregasi dari view
		konfigurasi.DB.Raw("SELECT nama_ng, SUM(jumlah_ng) as jumlah_ng FROM analytics_pareto_data GROUP BY nama_ng ORDER BY jumlah_ng DESC").Scan(&results)
	}()

	wg.Wait()
	c.JSON(http.StatusOK, utilitas.FormatRespons(true, "Berhasil mengambil data Pareto", results))
}

// GetControlChartData mengembalikan data P-Chart
// @Summary Ambil Data Control Chart
// @Description Kalkulasi paralel UCL, LCL, dan Rata-rata menggunakan Goroutine
// @Tags QC-Tools
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utilitas.ResponsAPI
// @Router /qc-tools/control-chart [get]
func GetControlChartData(c *gin.Context) {
	var ucl, lcl, avg float64
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		konfigurasi.DB.Raw("SELECT COALESCE(MAX(ucl), 0) FROM analytics_control_chart").Scan(&ucl)
	}()

	go func() {
		defer wg.Done()
		konfigurasi.DB.Raw("SELECT COALESCE(MIN(lcl), 0) FROM analytics_control_chart").Scan(&lcl)
	}()

	go func() {
		defer wg.Done()
		konfigurasi.DB.Raw("SELECT COALESCE(AVG(p_bar), 0) FROM analytics_control_chart").Scan(&avg)
	}()

	wg.Wait()

	c.JSON(http.StatusOK, utilitas.FormatRespons(true, "Berhasil kalkulasi Control Chart", gin.H{
		"ucl": ucl,
		"lcl": lcl,
		"avg": avg,
	}))
}
