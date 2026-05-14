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
func AmbilSemuaProduk(c *gin.Context) {
	var daftarProduk []model.Produk
	if err := konfigurasi.DB.Find(&daftarProduk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utilitas.FormatRespons(false, fmt.Sprintf("Gagal mengambil data produk: %v", err), nil))
		return
	}

	c.JSON(http.StatusOK, utilitas.FormatRespons(true, "Berhasil mengambil semua produk", daftarProduk))
}

// TambahProduk menambahkan data produk baru ke database
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
