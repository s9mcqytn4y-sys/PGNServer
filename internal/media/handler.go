package media

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pgn-server/pkg/respon"
)

type PenangananMedia struct {
	layanan LayananMedia
}

func KonstruksiPenangananBaru(layanan LayananMedia) *PenangananMedia {
	return &PenangananMedia{layanan: layanan}
}

// TanganiUnggahMedia menangani penerimaan berkas media untuk spesifik material.
// @Summary Unggah Media Material
// @Description Menyimpan foto atau bukti cacat terkait suatu material
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Material"
// @Param berkas formData file true "Berkas gambar (JPG/PNG, Maks 5MB)"
// @Success 200 {object} respon.ResponStandar
// @Failure 400 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/materials/{id}/media [post]
func (p *PenangananMedia) TanganiUnggahMedia(k *gin.Context) {
	idParam := k.Param("id")
	idMaterial, errParse := strconv.ParseUint(idParam, 10, 32)
	if errParse != nil {
		respon.Galat_Validasi(k, "ID Material tidak sah", nil)
		return
	}

	berkas, errUnggah := k.FormFile("berkas")
	if errUnggah != nil {
		respon.Galat_Validasi(k, "Berkas gagal dilampirkan atau tidak ditemukan", nil)
		return
	}

	aset, errProses := p.layanan.UnggahBerkasLokal(uint(idMaterial), berkas)
	if errProses != nil {
		pesanGalat := "Gagal mengunggah media"
		if errProses.Error() == "referensi_material_tidak_ditemukan" {
			respon.Galat_Validasi(k, "Material referensi tidak ditemukan di pangkalan data kami", nil)
			return
		} else if errProses.Error() == "ukuran_berkas_terlalu_besar" {
			respon.Galat_Validasi(k, "Ukuran berkas melebihi ambang batas (Maksimal 5MB)", nil)
			return
		} else if errProses.Error() == "ekstensi_berkas_tidak_diizinkan" || errProses.Error() == "tipe_mime_berkas_tidak_valid" {
			respon.Galat_Validasi(k, "Format berkas tidak diizinkan. Hanya mendukung PNG dan JPG/JPEG", nil)
			return
		}

		respon.Galat_Server(k, pesanGalat, errProses)
		return
	}

	respon.Sukses(k, "Media berhasil direkam dan dihubungkan ke material terkait.", aset)
}

// TanganiPratinjauMedia menyajikan preview gambar aset media.
// @Summary Pratinjau Media
// @Description Merender berkas media ke peramban
// @Tags Media
// @Produce image/jpeg
// @Param id path int true "ID Media"
// @Success 200 {file} file
// @Router /api/v1/media/{id}/pratinjau [get]
func (p *PenangananMedia) TanganiPratinjauMedia(k *gin.Context) {
	idParam := k.Param("id")
	idMedia, errParse := strconv.ParseUint(idParam, 10, 32)
	if errParse != nil {
		respon.Galat_Validasi(k, "Format ID Media tidak sah", nil)
		return
	}

	aset, errCari := p.layanan.CariAsetMedia(uint(idMedia))
	if errCari != nil {
		respon.Galat_TidakDitemukan(k, "Media tidak ditemukan")
		return
	}

	if aset.TipePenyimpanan == "EKSTERNAL" {
		k.Redirect(http.StatusTemporaryRedirect, aset.TautanEksternal)
		return
	}

	k.File(aset.DirektoriLokal)
}
