package pengelola

import (
	"encoding/json"
	"net/http"
	"pgn-server/internal/konfigurasi"
	"pgn-server/internal/model"
	"pgn-server/pkg/utilitas"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type ChecksheetInput struct {
	Header struct {
		UserID   uint   `json:"user_id" binding:"required"`
		LineID   string `json:"line_id" binding:"required"`
		ProdukID string `json:"produk_id" binding:"required"`
		MesinID  string `json:"mesin_id" binding:"required"`
		Shift    int    `json:"shift" binding:"required"`
	} `json:"header" binding:"required"`
	Details []struct {
		Proses         string                 `json:"proses" binding:"required"`
		ParameterHasil map[string]interface{} `json:"parameter_hasil" binding:"required"`
		Keterangan     string                 `json:"keterangan"`
	} `json:"details" binding:"required"`
}

// TambahChecksheet menangani penyimpanan checksheet transaksi
// @Summary Simpan Checksheet QC
// @Description Menyimpan data checksheet (Header & Detail) dalam satu transaksi
// @Tags QC
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

	// Memulai Transaksi
	tx := konfigurasi.DB.Begin()

	header := model.ChecksheetHeader{
		UserID:   input.Header.UserID,
		LineID:   input.Header.LineID,
		ProdukID: input.Header.ProdukID,
		MesinID:  input.Header.MesinID,
		Shift:    input.Header.Shift,
		Status:   "Closed", // Langsung closed setelah input
	}

	if err := tx.Create(&header).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utilitas.FormatRespons(false, "Gagal menyimpan header checksheet", nil))
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
			c.JSON(http.StatusInternalServerError, utilitas.FormatRespons(false, "Gagal menyimpan detail checksheet", nil))
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, utilitas.FormatRespons(true, "Checksheet berhasil disimpan", header))
}
