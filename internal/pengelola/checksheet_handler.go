package pengelola

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pgn-server/internal/konfigurasi"
	"pgn-server/internal/model"
	"pgn-server/pkg/utilitas"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type ChecksheetInput struct {
	Header struct {
		UserID   uint   `json:"user_id" binding:"required"`
		LineID   string `json:"line_id" binding:"required"`
		ProdukID string `json:"produk_id" binding:"required"`
		Shift    int    `json:"shift" binding:"required"`
	} `json:"header" binding:"required"`
	Details []struct {
		Proses         string                 `json:"proses" binding:"required"` // PRESS, SEWING, CUTTING
		ParameterHasil map[string]interface{} `json:"parameter_hasil" binding:"required"`
		Defects        []struct {
			DefectID string `json:"defect_id"`
			Jumlah   int    `json:"jumlah"`
		} `json:"defects"`
		Keterangan string `json:"keterangan"`
	} `json:"details" binding:"required"`
}

// TambahChecksheet menangani penyimpanan checksheet dan logging NG
// @Summary Simpan Checksheet QC & Log NG
// @Description Menyimpan data checksheet (Header & Detail) serta mencatat log defect secara otomatis
// @Tags QC-Checksheet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param checksheet body ChecksheetInput true "Data Checksheet"
// @Success 201 {object} utilitas.ResponsAPI
// @Failure 400 {object} utilitas.ResponsAPI
// @Router /qc/checksheet [post]
func TambahChecksheet(c *gin.Context) {
	var input ChecksheetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utilitas.FormatRespons(false, "Input tidak valid", nil))
		return
	}

	tx := konfigurasi.DB.Begin()

	header := model.ChecksheetHeader{
		UserID:   input.Header.UserID,
		LineID:   input.Header.LineID,
		ProdukID: input.Header.ProdukID,
		Shift:    input.Header.Shift,
		Status:   "Closed",
	}

	if err := tx.Create(&header).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utilitas.FormatRespons(false, "Gagal menyimpan header", nil))
		return
	}

	for _, d := range input.Details {
		paramJSON, _ := json.Marshal(d.ParameterHasil)
		detail := model.ChecksheetDetail{
			HeaderID:       header.ID,
			WaktuCek:       time.Now(),
			Proses:         d.Proses,
			ParameterHasil: datatypes.JSON(paramJSON),
			Keterangan:     d.Keterangan,
		}

		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, utilitas.FormatRespons(false, "Gagal menyimpan detail", nil))
			return
		}

		// Jika ada defect, catat ke log_inspeksi
		for _, def := range d.Defects {
			if def.Jumlah > 0 {
				logNG := model.LogInspeksi{
					InspeksiID:    fmt.Sprintf("INSP-%d", time.Now().Unix()),
					DefectID:      def.DefectID,
					JendelaWaktu:  time.Now().Format("15:04"),
					WaktuKejadian: time.Now(),
					JumlahNG:      def.Jumlah,
				}
				tx.Create(&logNG)
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, utilitas.FormatRespons(true, "Checksheet dan Log NG berhasil disimpan", header))
}

// GetFormChecksheet mengambil data defect potensial untuk form checksheet
// @Summary Ambil Form Checksheet Dinamis
// @Description Mengambil potensi defect material (via BOM) dan defect proses secara paralel
// @Tags QC-Checksheet
// @Produce json
// @Security BearerAuth
// @Param produk_id path string true "ID Produk"
// @Success 200 {object} utilitas.ResponsAPI
// @Router /qc/form-checksheet/{produk_id} [get]
func GetFormChecksheet(c *gin.Context) {
	produkID := c.Param("produk_id")
	db := konfigurasi.DB

	var produk model.Produk
	var processDefects []model.DefectMaster
	var wg sync.WaitGroup
	var errA, errB error

	wg.Add(2)
	go func() {
		defer wg.Done()
		errA = db.Preload("BOM.Material.PotensiDefect").First(&produk, "id = ?", produkID).Error
	}()
	go func() {
		defer wg.Done()
		errB = db.Where("kategori = ?", "PROCESS").Find(&processDefects).Error
	}()
	wg.Wait()

	if errA != nil {
		c.JSON(http.StatusNotFound, utilitas.FormatRespons(false, "Produk tidak ditemukan", nil))
		return
	}

	if errB != nil {
		c.JSON(http.StatusInternalServerError, utilitas.FormatRespons(false, "Gagal mengambil data defect proses", nil))
		return
	}

	data := gin.H{
		"produk":        produk,
		"defect_proses": processDefects,
		"server_time":   time.Now(),
	}

	c.JSON(http.StatusOK, utilitas.FormatRespons(true, "Data form checksheet berhasil diambil", data))
}
