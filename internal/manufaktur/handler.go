package manufaktur

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"pgn-server/pkg/respon"
)

type PenangananManufaktur struct {
	layanan LayananManufaktur
}

// KonstruksiPenangananBaru membuat instance baru PenangananManufaktur.
func KonstruksiPenangananBaru(layanan LayananManufaktur) *PenangananManufaktur {
	return &PenangananManufaktur{layanan: layanan}
}

// === HANDLER PEMASOK (SUPPLIER) ===

// TanganiTambahPemasok menambahkan pemasok baru.
// @Summary Tambah Pemasok Baru
// @Description Menyimpan data pemasok (supplier) baru ke pangkalan data
// @Tags Manufaktur - Pemasok
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body DTOPemasokSimpan true "Payload Tambah Pemasok"
// @Success 200 {object} respon.ResponStandar{data=Pemasok}
// @Failure 400 {object} respon.ResponStandar
// @Failure 401 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/suppliers [post]
func (p *PenangananManufaktur) TanganiTambahPemasok(k *gin.Context) {
	var dto DTOPemasokSimpan
	if err := k.ShouldBindJSON(&dto); err != nil {
		respon.Galat_Validasi(k, "Data input tidak valid", []string{err.Error()})
		return
	}

	hasil, err := p.layanan.TambahPemasok(&dto)
	if err != nil {
		respon.Galat_Server(k, "Gagal menambahkan pemasok baru", err)
		return
	}

	respon.Sukses(k, "Pemasok berhasil ditambahkan", hasil)
}

// TanganiAmbilSemuaPemasok mengambil daftar semua pemasok.
// @Summary Ambil Semua Pemasok
// @Description Menampilkan daftar seluruh pemasok bahan baku
// @Tags Manufaktur - Pemasok
// @Produce json
// @Security BearerAuth
// @Success 200 {object} respon.ResponStandar{data=[]Pemasok}
// @Failure 401 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/suppliers [get]
func (p *PenangananManufaktur) TanganiAmbilSemuaPemasok(k *gin.Context) {
	list, err := p.layanan.AmbilSemuaPemasok()
	if err != nil {
		respon.Galat_Server(k, "Gagal mengambil daftar pemasok", err)
		return
	}
	respon.Sukses(k, "Daftar pemasok berhasil diambil", list)
}

// TanganiCariPemasokID mencari pemasok berdasarkan ID.
// @Summary Detail Pemasok
// @Description Mendapatkan data detail pemasok berdasarkan ID
// @Tags Manufaktur - Pemasok
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Pemasok"
// @Success 200 {object} respon.ResponStandar{data=Pemasok}
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/suppliers/{id} [get]
func (p *PenangananManufaktur) TanganiCariPemasokID(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	hasil, errCari := p.layanan.CariPemasokID(uint(id))
	if errCari != nil {
		respon.Galat_TidakDitemukan(k, "Pemasok tidak ditemukan")
		return
	}
	respon.Sukses(k, "Detail pemasok berhasil ditemukan", hasil)
}

// TanganiUbahPemasok mengubah data pemasok.
// @Summary Perbarui Pemasok
// @Description Memperbarui informasi pemasok berdasarkan ID
// @Tags Manufaktur - Pemasok
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Pemasok"
// @Param payload body DTOPemasokSimpan true "Payload Perbarui Pemasok"
// @Success 200 {object} respon.ResponStandar{data=Pemasok}
// @Failure 400 {object} respon.ResponStandar
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/suppliers/{id} [put]
func (p *PenangananManufaktur) TanganiUbahPemasok(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	var dto DTOPemasokSimpan
	if errBind := k.ShouldBindJSON(&dto); errBind != nil {
		respon.Galat_Validasi(k, "Data input tidak valid", []string{errBind.Error()})
		return
	}

	hasil, errUbah := p.layanan.UbahPemasok(uint(id), &dto)
	if errUbah != nil {
		if errUbah.Error() == "pemasok_tidak_ditemukan" {
			respon.Galat_TidakDitemukan(k, "Pemasok yang ingin diubah tidak ditemukan")
			return
		}
		respon.Galat_Server(k, "Gagal mengubah data pemasok", errUbah)
		return
	}

	respon.Sukses(k, "Data pemasok berhasil diperbarui", hasil)
}

// TanganiHapusPemasok menghapus data pemasok.
// @Summary Hapus Pemasok
// @Description Menghapus data pemasok dari sistem berdasarkan ID
// @Tags Manufaktur - Pemasok
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Pemasok"
// @Success 200 {object} respon.ResponStandar
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/suppliers/{id} [delete]
func (p *PenangananManufaktur) TanganiHapusPemasok(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	if errHapus := p.layanan.HapusPemasokID(uint(id)); errHapus != nil {
		if errHapus.Error() == "pemasok_tidak_ditemukan" {
			respon.Galat_TidakDitemukan(k, "Pemasok tidak ditemukan")
			return
		}
		respon.Galat_Server(k, "Gagal menghapus pemasok", errHapus)
		return
	}

	respon.Sukses(k, "Pemasok berhasil dihapus dari sistem", nil)
}

// === HANDLER MATERIAL ===

// TanganiTambahMaterial menambahkan material baru.
// @Summary Tambah Material Baru
// @Description Menyimpan data material baru ke pangkalan data
// @Tags Manufaktur - Material
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body DTOMaterialSimpan true "Payload Tambah Material"
// @Success 200 {object} respon.ResponStandar{data=Material}
// @Failure 400 {object} respon.ResponStandar
// @Failure 401 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/materials [post]
func (p *PenangananManufaktur) TanganiTambahMaterial(k *gin.Context) {
	var dto DTOMaterialSimpan
	if err := k.ShouldBindJSON(&dto); err != nil {
		respon.Galat_Validasi(k, "Data input tidak valid", []string{err.Error()})
		return
	}

	hasil, err := p.layanan.TambahMaterial(&dto)
	if err != nil {
		if err.Error() == "pemasok_tidak_ditemukan" {
			respon.Galat_Validasi(k, "Pemasok yang dipilih tidak ditemukan", nil)
			return
		}
		respon.Galat_Server(k, "Gagal menambahkan material baru", err)
		return
	}

	respon.Sukses(k, "Material berhasil ditambahkan", hasil)
}

// TanganiAmbilSemuaMaterial mengambil daftar semua material.
// @Summary Ambil Semua Material
// @Description Menampilkan daftar seluruh material beserta info pemasoknya
// @Tags Manufaktur - Material
// @Produce json
// @Security BearerAuth
// @Success 200 {object} respon.ResponStandar{data=[]Material}
// @Failure 401 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/materials [get]
func (p *PenangananManufaktur) TanganiAmbilSemuaMaterial(k *gin.Context) {
	list, err := p.layanan.AmbilSemuaMaterial()
	if err != nil {
		respon.Galat_Server(k, "Gagal mengambil daftar material", err)
		return
	}
	respon.Sukses(k, "Daftar material berhasil diambil", list)
}

// TanganiCariMaterialID mencari material berdasarkan ID.
// @Summary Detail Material
// @Description Mendapatkan data detail material berdasarkan ID
// @Tags Manufaktur - Material
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Material"
// @Success 200 {object} respon.ResponStandar{data=Material}
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/materials/{id} [get]
func (p *PenangananManufaktur) TanganiCariMaterialID(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	hasil, errCari := p.layanan.CariMaterialID(uint(id))
	if errCari != nil {
		respon.Galat_TidakDitemukan(k, "Material tidak ditemukan")
		return
	}
	respon.Sukses(k, "Detail material berhasil ditemukan", hasil)
}

// TanganiUbahMaterial mengubah data material.
// @Summary Perbarui Material
// @Description Memperbarui informasi material berdasarkan ID
// @Tags Manufaktur - Material
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Material"
// @Param payload body DTOMaterialSimpan true "Payload Perbarui Material"
// @Success 200 {object} respon.ResponStandar{data=Material}
// @Failure 400 {object} respon.ResponStandar
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/materials/{id} [put]
func (p *PenangananManufaktur) TanganiUbahMaterial(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	var dto DTOMaterialSimpan
	if errBind := k.ShouldBindJSON(&dto); errBind != nil {
		respon.Galat_Validasi(k, "Data input tidak valid", []string{errBind.Error()})
		return
	}

	hasil, errUbah := p.layanan.UbahMaterial(uint(id), &dto)
	if errUbah != nil {
		if errUbah.Error() == "material_tidak_ditemukan" {
			respon.Galat_TidakDitemukan(k, "Material tidak ditemukan")
			return
		}
		if errUbah.Error() == "pemasok_tidak_ditemukan" {
			respon.Galat_Validasi(k, "Pemasok tidak ditemukan", nil)
			return
		}
		respon.Galat_Server(k, "Gagal mengubah data material", errUbah)
		return
	}

	respon.Sukses(k, "Data material berhasil diperbarui", hasil)
}

// TanganiHapusMaterial menghapus data material.
// @Summary Hapus Material
// @Description Menghapus data material dari sistem berdasarkan ID
// @Tags Manufaktur - Material
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Material"
// @Success 200 {object} respon.ResponStandar
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/materials/{id} [delete]
func (p *PenangananManufaktur) TanganiHapusMaterial(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	if errHapus := p.layanan.HapusMaterialID(uint(id)); errHapus != nil {
		if errHapus.Error() == "material_tidak_ditemukan" {
			respon.Galat_TidakDitemukan(k, "Material tidak ditemukan")
			return
		}
		respon.Galat_Server(k, "Gagal menghapus material", errHapus)
		return
	}

	respon.Sukses(k, "Material berhasil dihapus dari sistem", nil)
}

// === HANDLER CUSTOMER ===

// TanganiTambahCustomer menambahkan customer baru.
// @Summary Tambah Customer Baru
// @Description Menyimpan data customer baru ke pangkalan data
// @Tags Manufaktur - Customer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body DTOCustomerSimpan true "Payload Tambah Customer"
// @Success 200 {object} respon.ResponStandar{data=Customer}
// @Failure 400 {object} respon.ResponStandar
// @Failure 401 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/customers [post]
func (p *PenangananManufaktur) TanganiTambahCustomer(k *gin.Context) {
	var dto DTOCustomerSimpan
	if err := k.ShouldBindJSON(&dto); err != nil {
		respon.Galat_Validasi(k, "Data input tidak valid", []string{err.Error()})
		return
	}

	hasil, err := p.layanan.TambahCustomer(&dto)
	if err != nil {
		respon.Galat_Server(k, "Gagal menambahkan customer baru", err)
		return
	}

	respon.Sukses(k, "Customer berhasil ditambahkan", hasil)
}

// TanganiAmbilSemuaCustomer mengambil daftar semua customer.
// @Summary Ambil Semua Customer
// @Description Menampilkan daftar seluruh customer
// @Tags Manufaktur - Customer
// @Produce json
// @Security BearerAuth
// @Success 200 {object} respon.ResponStandar{data=[]Customer}
// @Failure 401 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/customers [get]
func (p *PenangananManufaktur) TanganiAmbilSemuaCustomer(k *gin.Context) {
	list, err := p.layanan.AmbilSemuaCustomer()
	if err != nil {
		respon.Galat_Server(k, "Gagal mengambil daftar customer", err)
		return
	}
	respon.Sukses(k, "Daftar customer berhasil diambil", list)
}

// TanganiCariCustomerID mencari customer berdasarkan ID.
// @Summary Detail Customer
// @Description Mendapatkan data detail customer berdasarkan ID
// @Tags Manufaktur - Customer
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Customer"
// @Success 200 {object} respon.ResponStandar{data=Customer}
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/customers/{id} [get]
func (p *PenangananManufaktur) TanganiCariCustomerID(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	hasil, errCari := p.layanan.CariCustomerID(uint(id))
	if errCari != nil {
		respon.Galat_TidakDitemukan(k, "Customer tidak ditemukan")
		return
	}
	respon.Sukses(k, "Detail customer berhasil ditemukan", hasil)
}

// TanganiUbahCustomer mengubah data customer.
// @Summary Perbarui Customer
// @Description Memperbarui informasi customer berdasarkan ID
// @Tags Manufaktur - Customer
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Customer"
// @Param payload body DTOCustomerSimpan true "Payload Perbarui Customer"
// @Success 200 {object} respon.ResponStandar{data=Customer}
// @Failure 400 {object} respon.ResponStandar
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/customers/{id} [put]
func (p *PenangananManufaktur) TanganiUbahCustomer(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	var dto DTOCustomerSimpan
	if errBind := k.ShouldBindJSON(&dto); errBind != nil {
		respon.Galat_Validasi(k, "Data input tidak valid", []string{errBind.Error()})
		return
	}

	hasil, errUbah := p.layanan.UbahCustomer(uint(id), &dto)
	if errUbah != nil {
		if errUbah.Error() == "customer_tidak_ditemukan" {
			respon.Galat_TidakDitemukan(k, "Customer tidak ditemukan")
			return
		}
		respon.Galat_Server(k, "Gagal mengubah data customer", errUbah)
		return
	}

	respon.Sukses(k, "Data customer berhasil diperbarui", hasil)
}

// TanganiHapusCustomer menghapus data customer.
// @Summary Hapus Customer
// @Description Menghapus data customer dari sistem berdasarkan ID
// @Tags Manufaktur - Customer
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Customer"
// @Success 200 {object} respon.ResponStandar
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/customers/{id} [delete]
func (p *PenangananManufaktur) TanganiHapusCustomer(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	if errHapus := p.layanan.HapusCustomerID(uint(id)); errHapus != nil {
		if errHapus.Error() == "customer_tidak_ditemukan" {
			respon.Galat_TidakDitemukan(k, "Customer tidak ditemukan")
			return
		}
		respon.Galat_Server(k, "Gagal menghapus customer", errHapus)
		return
	}

	respon.Sukses(k, "Customer berhasil dihapus dari sistem", nil)
}

// === HANDLER BOM (BILL OF MATERIALS) ===

// TanganiTambahBOM menambahkan BOM baru.
// @Summary Tambah BOM Baru
// @Description Menyimpan relasi hirarkis BOM baru ke pangkalan data
// @Tags Manufaktur - BOM
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body DTOBOMSimpan true "Payload Tambah BOM"
// @Success 200 {object} respon.ResponStandar{data=KomposisiMaterialBOM}
// @Failure 400 {object} respon.ResponStandar
// @Failure 401 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/boms [post]
func (p *PenangananManufaktur) TanganiTambahBOM(k *gin.Context) {
	var dto DTOBOMSimpan
	if err := k.ShouldBindJSON(&dto); err != nil {
		respon.Galat_Validasi(k, "Data input tidak valid", []string{err.Error()})
		return
	}

	hasil, err := p.layanan.TambahBOM(&dto)
	if err != nil {
		if err.Error() == "material_baku_tidak_ditemukan" || err.Error() == "material_induk_tidak_ditemukan" {
			respon.Galat_Validasi(k, err.Error(), nil)
			return
		}
		respon.Galat_Server(k, "Gagal menambahkan BOM baru", err)
		return
	}

	respon.Sukses(k, "BOM berhasil ditambahkan", hasil)
}

// TanganiAmbilSemuaBOM mengambil daftar semua BOM.
// @Summary Ambil Semua BOM
// @Description Menampilkan daftar seluruh relasi BOM beserta detail material terkait
// @Tags Manufaktur - BOM
// @Produce json
// @Security BearerAuth
// @Success 200 {object} respon.ResponStandar{data=[]KomposisiMaterialBOM}
// @Failure 401 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/boms [get]
func (p *PenangananManufaktur) TanganiAmbilSemuaBOM(k *gin.Context) {
	list, err := p.layanan.AmbilSemuaBOM()
	if err != nil {
		respon.Galat_Server(k, "Gagal mengambil daftar BOM", err)
		return
	}
	respon.Sukses(k, "Daftar BOM berhasil diambil", list)
}

// TanganiCariBOMID mencari BOM berdasarkan ID.
// @Summary Detail BOM
// @Description Mendapatkan data detail relasi BOM berdasarkan ID
// @Tags Manufaktur - BOM
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID BOM"
// @Success 200 {object} respon.ResponStandar{data=KomposisiMaterialBOM}
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/boms/{id} [get]
func (p *PenangananManufaktur) TanganiCariBOMID(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	hasil, errCari := p.layanan.CariBOMID(uint(id))
	if errCari != nil {
		respon.Galat_TidakDitemukan(k, "BOM tidak ditemukan")
		return
	}
	respon.Sukses(k, "Detail BOM berhasil ditemukan", hasil)
}

// TanganiUbahBOM mengubah data BOM.
// @Summary Perbarui BOM
// @Description Memperbarui informasi relasi BOM berdasarkan ID
// @Tags Manufaktur - BOM
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID BOM"
// @Param payload body DTOBOMSimpan true "Payload Perbarui BOM"
// @Success 200 {object} respon.ResponStandar{data=KomposisiMaterialBOM}
// @Failure 400 {object} respon.ResponStandar
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/boms/{id} [put]
func (p *PenangananManufaktur) TanganiUbahBOM(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	var dto DTOBOMSimpan
	if errBind := k.ShouldBindJSON(&dto); errBind != nil {
		respon.Galat_Validasi(k, "Data input tidak valid", []string{errBind.Error()})
		return
	}

	hasil, errUbah := p.layanan.UbahBOM(uint(id), &dto)
	if errUbah != nil {
		if errUbah.Error() == "bom_tidak_ditemukan" {
			respon.Galat_TidakDitemukan(k, "BOM tidak ditemukan")
			return
		}
		if errUbah.Error() == "material_baku_tidak_ditemukan" || errUbah.Error() == "material_induk_tidak_ditemukan" {
			respon.Galat_Validasi(k, errUbah.Error(), nil)
			return
		}
		respon.Galat_Server(k, "Gagal mengubah data BOM", errUbah)
		return
	}

	respon.Sukses(k, "Data BOM berhasil diperbarui", hasil)
}

// TanganiHapusBOM menghapus data BOM.
// @Summary Hapus BOM
// @Description Menghapus data relasi BOM dari sistem berdasarkan ID
// @Tags Manufaktur - BOM
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID BOM"
// @Success 200 {object} respon.ResponStandar
// @Failure 404 {object} respon.ResponStandar
// @Failure 500 {object} respon.ResponStandar
// @Router /api/v1/boms/{id} [delete]
func (p *PenangananManufaktur) TanganiHapusBOM(k *gin.Context) {
	idParam := k.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respon.Galat_Validasi(k, "Format ID tidak valid", nil)
		return
	}

	if errHapus := p.layanan.HapusBOMID(uint(id)); errHapus != nil {
		if errHapus.Error() == "bom_tidak_ditemukan" {
			respon.Galat_TidakDitemukan(k, "BOM tidak ditemukan")
			return
		}
		respon.Galat_Server(k, "Gagal menghapus BOM", errHapus)
		return
	}

	respon.Sukses(k, "BOM berhasil dihapus dari sistem", nil)
}
