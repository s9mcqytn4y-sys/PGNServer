package analitik

import (
	"testing"
	"time"

	"pgn-server/internal/infrastruktur"
	"pgn-server/internal/kualitas"
	"pgn-server/internal/manufaktur"

	"github.com/stretchr/testify/assert"
)

func TestLayananAnalitik_DapatkanRekomendasiTindakan(t *testing.T) {
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
	repoKualitas := kualitas.KonstruksiRepositoriBaru(db)

	hariIni := time.Now().Format("2006-01-02")
	kemarin := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// Helper func to clear and seed inspections
	resetAndSeedInspeksi := func(details []kualitas.DTODetailInspeksi, tgl string) {
		db.Exec("TRUNCATE TABLE detail_inspeksis CASCADE")
		db.Exec("TRUNCATE TABLE lembar_periksas CASCADE")
		db.Exec("TRUNCATE TABLE buku_besar_cacats CASCADE")

		if len(details) > 0 {
			dto := &kualitas.DTOLembarPeriksaKirim{
				Tanggal:            tgl,
				ZonaLini:           "Lini-DSS-Test",
				PenggunaIDTercatat: 1,
				Detail:             details,
			}
			tx := db.Begin()
			err := repoKualitas.SimpanLembarPeriksaMassal(dto, tx)
			assert.NoError(t, err)
			tx.Commit()
		}
	}

	t.Run("1. No Data (TotalProduksi == 0)", func(t *testing.T) {
		resetAndSeedInspeksi(nil, hariIni) // No data

		resp, err := layanan.DapatkanRekomendasiTindakan(nil, hariIni, hariIni, "Lini-DSS-Test", "")
		assert.NoError(t, err)
		assert.Equal(t, "STABIL", resp.Status)
		assert.Contains(t, resp.Ringkasan, "Total Produksi = 0")
		assert.Empty(t, resp.Rekomendasi)
	})

	t.Run("2. Critical Ratio (RasioNG >= AmbangRasioNGKritis)", func(t *testing.T) {
		// Target > 5.0%
		details := []kualitas.DTODetailInspeksi{
			{UnikPartID: 1, KodeCacat: "LAIN-LAIN", WaktuPergeseran: "Shift-1", TotalProduksi: 100, RasioTotalOK: 90, RasioCacat: 10}, // 10% NG
		}
		resetAndSeedInspeksi(details, hariIni)

		resp, err := layanan.DapatkanRekomendasiTindakan(nil, hariIni, hariIni, "Lini-DSS-Test", "")
		assert.NoError(t, err)
		assert.Equal(t, "KRITIS", resp.Status)
		assert.True(t, resp.Indikator.RasioNG >= float64(AmbangRasioNGKritis))

		// Should have MANAGEMENT and QA recommendations
		var hasManagement, hasQA bool
		for _, rec := range resp.Rekomendasi {
			if rec.Target == "MANAGEMENT" {
				hasManagement = true
			}
			if rec.Target == "QA" {
				hasQA = true
			}
		}
		assert.True(t, hasManagement)
		assert.True(t, hasQA)
	})

	t.Run("3. Warning Ratio (RasioNG >= AmbangRasioNGWaspada)", func(t *testing.T) {
		// Target between 2.0% and 5.0%
		details := []kualitas.DTODetailInspeksi{
			{UnikPartID: 1, KodeCacat: "LAIN-LAIN", WaktuPergeseran: "Shift-1", TotalProduksi: 100, RasioTotalOK: 97, RasioCacat: 3}, // 3% NG
		}
		resetAndSeedInspeksi(details, hariIni)

		resp, err := layanan.DapatkanRekomendasiTindakan(nil, hariIni, hariIni, "Lini-DSS-Test", "")
		assert.NoError(t, err)
		assert.Equal(t, "WASPADA", resp.Status)
		assert.True(t, resp.Indikator.RasioNG >= float64(AmbangRasioNGWaspada))
		assert.True(t, resp.Indikator.RasioNG < float64(AmbangRasioNGKritis))

		var hasQA, hasQC bool
		for _, rec := range resp.Rekomendasi {
			if rec.Target == "QA" {
				hasQA = true
			}
			if rec.Target == "QC" {
				hasQC = true
			}
		}
		assert.True(t, hasQA)
		assert.True(t, hasQC)
	})

	t.Run("4. Pareto Dominance Material (Pareto >= 50% & Material)", func(t *testing.T) {
		// LAMINATING BOLONG is material defect based on DefectToMaterialMapping
		details := []kualitas.DTODetailInspeksi{
			{UnikPartID: 1, KodeCacat: "LAMINATING BOLONG", WaktuPergeseran: "Shift-1", TotalProduksi: 100, RasioTotalOK: 99, RasioCacat: 1}, // 100% of defect
		}
		resetAndSeedInspeksi(details, hariIni)

		resp, err := layanan.DapatkanRekomendasiTindakan(nil, hariIni, hariIni, "Lini-DSS-Test", "")
		assert.NoError(t, err)

		// Ratio is 1%, so overall status is STABIL
		assert.Equal(t, "STABIL", resp.Status)

		var hasSupplier, hasPCD bool
		for _, rec := range resp.Rekomendasi {
			if rec.Target == "SUPPLIER" {
				hasSupplier = true
			}
			if rec.Target == "PCD" {
				hasPCD = true
			}
		}
		assert.True(t, hasSupplier)
		assert.True(t, hasPCD)
	})

	t.Run("5. Pareto Dominance Process (Pareto >= 50% & Process)", func(t *testing.T) {
		// SEWING MIRING is a process defect based on ApakahDefectProcess helper
		details := []kualitas.DTODetailInspeksi{
			{UnikPartID: 1, KodeCacat: "SEWING MIRING", WaktuPergeseran: "Shift-1", TotalProduksi: 100, RasioTotalOK: 99, RasioCacat: 1},
		}
		resetAndSeedInspeksi(details, hariIni)

		resp, err := layanan.DapatkanRekomendasiTindakan(nil, hariIni, hariIni, "Lini-DSS-Test", "")
		assert.NoError(t, err)

		var hasQC bool
		for _, rec := range resp.Rekomendasi {
			if rec.Target == "QC" {
				hasQC = true
			}
		}
		assert.True(t, hasQC)
	})

	t.Run("6. Material Issue PCD Trigger (RasioNG >= 2.0 & Material)", func(t *testing.T) {
		details := []kualitas.DTODetailInspeksi{
			{UnikPartID: 1, KodeCacat: "LAMINATING BOLONG", WaktuPergeseran: "Shift-1", TotalProduksi: 100, RasioTotalOK: 97, RasioCacat: 3}, // 3% NG
		}
		resetAndSeedInspeksi(details, hariIni)

		resp, err := layanan.DapatkanRekomendasiTindakan(nil, hariIni, hariIni, "Lini-DSS-Test", "")
		assert.NoError(t, err)
		assert.Equal(t, "WASPADA", resp.Status)

		pcdHighPriority := false
		for _, rec := range resp.Rekomendasi {
			if rec.Target == "PCD" && rec.Prioritas == "TINGGI" {
				pcdHighPriority = true
			}
		}
		assert.True(t, pcdHighPriority, "Expected high priority recommendation for PCD due to high material defect ratio")
	})

	t.Run("7. Upward Trend (Trend == NAIK)", func(t *testing.T) {
		db.Exec("TRUNCATE TABLE detail_inspeksis CASCADE")
		db.Exec("TRUNCATE TABLE lembar_periksas CASCADE")
		db.Exec("TRUNCATE TABLE buku_besar_cacats CASCADE")

		duaHariLalu := time.Now().AddDate(0, 0, -2).Format("2006-01-02")

		// Day 1 (2 hari lalu) - 1 defect
		tx0 := db.Begin()
		err0 := repoKualitas.SimpanLembarPeriksaMassal(&kualitas.DTOLembarPeriksaKirim{
			Tanggal:            duaHariLalu,
			ZonaLini:           "Lini-DSS-Test",
			PenggunaIDTercatat: 1,
			Detail: []kualitas.DTODetailInspeksi{
				{UnikPartID: 1, KodeCacat: "LAIN-LAIN", WaktuPergeseran: "Shift-1", TotalProduksi: 100, RasioTotalOK: 99, RasioCacat: 1},
			},
		}, tx0)
		assert.NoError(t, err0)
		tx0.Commit()

		// Day 2 (kemarin) - 2 defects
		tx1 := db.Begin()
		err1 := repoKualitas.SimpanLembarPeriksaMassal(&kualitas.DTOLembarPeriksaKirim{
			Tanggal:            kemarin,
			ZonaLini:           "Lini-DSS-Test",
			PenggunaIDTercatat: 1,
			Detail: []kualitas.DTODetailInspeksi{
				{UnikPartID: 1, KodeCacat: "LAIN-LAIN", WaktuPergeseran: "Shift-1", TotalProduksi: 100, RasioTotalOK: 98, RasioCacat: 2},
			},
		}, tx1)
		assert.NoError(t, err1)
		tx1.Commit()

		// Day 3 (hari ini) - 3 defects
		tx2 := db.Begin()
		err2 := repoKualitas.SimpanLembarPeriksaMassal(&kualitas.DTOLembarPeriksaKirim{
			Tanggal:            hariIni,
			ZonaLini:           "Lini-DSS-Test",
			PenggunaIDTercatat: 1,
			Detail: []kualitas.DTODetailInspeksi{
				{UnikPartID: 1, KodeCacat: "LAIN-LAIN", WaktuPergeseran: "Shift-1", TotalProduksi: 100, RasioTotalOK: 97, RasioCacat: 3},
			},
		}, tx2)
		assert.NoError(t, err2)
		tx2.Commit()

		resp, err := layanan.DapatkanRekomendasiTindakan(nil, hariIni, hariIni, "Lini-DSS-Test", "")
		assert.NoError(t, err)
		assert.Equal(t, "NAIK", resp.Indikator.Trend7Hari)
		// Since total NG today is 3/100 (3%), normal status is WASPADA, NAIK keeps it WASPADA
		assert.Equal(t, "WASPADA", resp.Status)
	})
}
