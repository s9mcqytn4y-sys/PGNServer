// Package pengelola berisi handler untuk manajemen bisnis logic dan API PGNServer.
package pengelola

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"pgn-server/pkg/utilitas"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadHandler menangani upload file media (gambar)
// @Summary Upload Media
// @Description Mengunggah file gambar ke folder penyimpanan
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File gambar"
// @Success 200 {object} utilitas.ResponsAPI
// @Failure 400 {object} utilitas.ResponsAPI
// @Failure 500 {object} utilitas.ResponsAPI
// @Router /media/upload [post]
func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utilitas.FormatRespons(false, "File tidak ditemukan", nil))
		return
	}

	// Validasi Ekstensi
	ext := filepath.Ext(file.Filename)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, utilitas.FormatRespons(false, "Format file tidak didukung (hanya jpg, jpeg, png)", nil))
		return
	}

	// Buat folder jika belum ada
	uploadDir := "./penyimpanan"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// Nama file unik
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, utilitas.FormatRespons(false, "Gagal menyimpan file", nil))
		return
	}

	fileURL := fmt.Sprintf("/penyimpanan/%s", filename)
	c.JSON(http.StatusOK, utilitas.FormatRespons(true, "File berhasil diunggah", gin.H{
		"url":      fileURL,
		"filename": filename,
	}))
}
