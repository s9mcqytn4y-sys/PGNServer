package pengelola

import (
	"fmt"
	"net/http"
	"pgn-server/internal/konfigurasi"
	"pgn-server/internal/model"
	"pgn-server/pkg/utilitas"

	"github.com/gin-gonic/gin"
)

// AmbilSemuaProduk mengembalikan daftar seluruh produk
// @Summary Ambil Semua Produk
// @Description Mengambil semua data produk dari database
// @Tags Produk
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utilitas.ResponsAPI
// @Failure 500 {object} utilitas.ResponsAPI
// @Router /produk [get]
func AmbilSemuaProduk(c *gin.Context) {
	var daftarProduk []model.Produk
	if err := konfigurasi.DB.Find(&daftarProduk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utilitas.FormatRespons(false, fmt.Sprintf("Gagal mengambil data produk: %v", err), nil))
		return
	}

	c.JSON(http.StatusOK, utilitas.FormatRespons(true, "Berhasil mengambil semua produk", daftarProduk))
}

// TambahProduk menambahkan data produk baru ke database
// @Summary Tambah Produk Baru
// @Description Menyimpan data produk baru ke database
// @Tags Produk
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param produk body model.Produk true "Data Produk"
// @Success 201 {object} utilitas.ResponsAPI
// @Failure 400 {object} utilitas.ResponsAPI
// @Failure 500 {object} utilitas.ResponsAPI
// @Router /produk [post]
func TambahProduk(c *gin.Context) {
	var input model.Produk
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utilitas.FormatRespons(false, "Input tidak valid", nil))
		return
	}

	if err := konfigurasi.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utilitas.FormatRespons(false, fmt.Sprintf("Gagal menyimpan produk: %v", err), nil))
		return
	}

	c.JSON(http.StatusCreated, utilitas.FormatRespons(true, "Produk berhasil ditambahkan", input))
}
