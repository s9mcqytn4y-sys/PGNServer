package analitik

import (
	"fmt"
	"os"
	"testing"

	"pgn-server/internal/infrastruktur"
	"pgn-server/internal/kualitas"
	"pgn-server/internal/manufaktur"
	"pgn-server/pkg/cache"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dapatkanKoneksiDBTest() (*gorm.DB, error) {
	inangDB := os.Getenv("DB_HOST")
	if inangDB == "" {
		inangDB = "localhost"
	}
	penggunaDB := os.Getenv("DB_USER")
	if penggunaDB == "" {
		penggunaDB = "admin"
	}
	sandiDB := os.Getenv("DB_PASSWORD")
	if sandiDB == "" {
		sandiDB = "admin"
	}
	namaDB := os.Getenv("DB_NAME")
	if namaDB == "" {
		namaDB = "pgn_db"
	}
	pelabuhanDB := os.Getenv("DB_PORT")
	if pelabuhanDB == "" {
		pelabuhanDB = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		inangDB, penggunaDB, sandiDB, namaDB, pelabuhanDB)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func TestLayananAnalitik_ParetoDanLacak(t *testing.T) {
	db, err := dapatkanKoneksiDBTest()
	if err != nil {
		t.Skip("Pangkalan data pengujian tidak dapat dihubungi, lewati tes integrasi")
		return
	}

	// Drop tables first to guarantee clean state for tests
	_ = db.Migrator().DropTable(
		"lembar_periksas",
		"detail_inspeksis",
		"buku_besar_cacats",
		"bill_of_materials",
		"MATERIAL",
		"suppliers",
		"customers",
		"aset_digitals",
		"penggunas",
	)

	// Jalankan migrasi
	_ = infrastruktur.PelaksanaanAutoMigrasi(db)
	_ = manufaktur.JalankanSeeder(db)

	repo := KonstruksiRepositoriBaru(db)
	layanan := KonstruksiLayananBaru(repo, db)

	t.Run("Kalkulasi Pareto Dinamis", func(t *testing.T) {
		// Simpan lembar periksa dummy untuk testing
		repoKualitas := kualitas.KonstruksiRepositoriBaru(db)
		dto := &kualitas.DTOLembarPeriksaKirim{
			Tanggal:            "2026-05-17",
			ZonaLini:           "Lini-Test-1",
			PenggunaIDTercatat: 1,
			Detail: []kualitas.DTODetailInspeksi{
				{UnikPartID: 1, KodeCacat: "LAMINATING BOLONG", WaktuPergeseran: "Shift-1", TotalProduksi: 100, RasioTotalOK: 95, RasioCacat: 5},
				{UnikPartID: 1, KodeCacat: "LAMINATING TIDAK MATANG", WaktuPergeseran: "Shift-1", TotalProduksi: 100, RasioTotalOK: 98, RasioCacat: 2},
			},
		}

		tx := db.Begin()
		errSimpan := repoKualitas.SimpanLembarPeriksaMassal(dto, tx)
		assert.NoError(t, errSimpan)
		tx.Commit()

		// Test Kalkulasi Pareto
		pareto, errPareto := layanan.DapatkanPareto(nil, "2026-05-17", "2026-05-17", "Lini-Test-1")
		assert.NoError(t, errPareto)
		assert.NotEmpty(t, pareto)

		// Verifikasi order by DESC & persentase kumulatif
		if len(pareto) >= 2 {
			assert.True(t, pareto[0].JumlahCacat >= pareto[1].JumlahCacat)
			assert.Equal(t, float64(100), pareto[len(pareto)-1].PersentaseKumulatif)
		}
	})

	t.Run("BOM Tracing Akar Masalah Cacat Material", func(t *testing.T) {
		// Lacak cacat LAMINATING BOLONG pada Protector
		lacak, errLacak := layanan.LacakAkarMasalah(nil, "LAMINATING BOLONG", "FG-002")
		assert.NoError(t, errLacak)
		assert.NotEmpty(t, lacak)

		// Verifikasi raw material pembentuk dan pemasok
		found := false
		for _, item := range lacak {
			if item.RawMaterialNama == "Laminasi LDPE 200 Gsm" {
				assert.Equal(t, "PT Artha Langgeng Mulya (BTI)", item.PemasokNama)
				found = true
			}
		}
		assert.True(t, found, "Harus menemukan Laminasi LDPE 200 Gsm disuplai oleh BTI")
	})

	t.Run("Circular Dependency Detection", func(t *testing.T) {
		// Buat circular dependency buatan di BOM
		var mat1, mat2 manufaktur.Material
		db.Where("kode_sku = ?", "MAT-001").First(&mat1)
		db.Where("kode_sku = ?", "MAT-002").First(&mat2)

		// Hapus komposisi lama jika ada
		db.Where("id_parent_material = ? AND id_raw_material = ?", mat1.ID, mat2.ID).Delete(&manufaktur.KomposisiMaterialBOM{})
		db.Where("id_parent_material = ? AND id_raw_material = ?", mat2.ID, mat1.ID).Delete(&manufaktur.KomposisiMaterialBOM{})

		// C -> A & A -> C
		bom1 := manufaktur.KomposisiMaterialBOM{IDParentMaterial: &mat1.ID, IDRawMaterial: mat2.ID, ParameterKuantitasPembentuk: 1.0}
		bom2 := manufaktur.KomposisiMaterialBOM{IDParentMaterial: &mat2.ID, IDRawMaterial: mat1.ID, ParameterKuantitasPembentuk: 1.0}

		db.Create(&bom1)
		db.Create(&bom2)

		defer func() {
			db.Delete(&bom1)
			db.Delete(&bom2)
			cache.GlobalCache.Clear()
		}()

		// Bersihkan cache sebelum pelacakan
		cache.GlobalCache.Clear()

		lacak, errLacak := layanan.LacakAkarMasalah(nil, "SPUNBOUND TIDAK MEREKAT", mat1.KodeSKU)
		assert.NoError(t, errLacak)
		
		circularDetected := false
		for _, item := range lacak {
			if item.CircularDependencyDetected {
				circularDetected = true
			}
		}
		assert.True(t, circularDetected, "Harus mendeteksi dependensi melingkar (circular dependency)")
	})
}

func TestLayananAnalitik_ProcessDefect(t *testing.T) {
	db, err := dapatkanKoneksiDBTest()
	if err != nil {
		t.Skip("Pangkalan data pengujian tidak dapat dihubungi, lewati tes integrasi")
		return
	}

	repo := KonstruksiRepositoriBaru(db)
	layanan := KonstruksiLayananBaru(repo, db)

	t.Run("Internal Process Defect", func(t *testing.T) {
		lacak, errLacak := layanan.LacakAkarMasalah(nil, "SEWING MIRING", "FG-002")
		assert.NoError(t, errLacak)
		assert.NotEmpty(t, lacak)
		assert.Equal(t, "INTERNAL PROCESS DEFECT", lacak[0].RawMaterialNama)
		assert.Equal(t, "INTERNAL PRODUCTION LINE", lacak[0].PemasokNama)
	})
}
