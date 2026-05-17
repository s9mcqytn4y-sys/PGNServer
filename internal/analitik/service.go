// Package analitik menangani agregasi dan kalkulasi data pelaporan 7 QC Tools.
package analitik

import (
	"errors"
	"strings"

	"pgn-server/internal/manufaktur"

	"gorm.io/gorm"
)

// DefectToMaterialMapping memetakan kode cacat ke nama material pembentuknya berdasarkan PARTLIST.xlsx
var DefectToMaterialMapping = map[string][]string{
	"SPUNBOUND TIDAK MEREKAT": {"PS Polyester Non Woven Spunbond 100 Gsm White"},
	"SPUNBOUND KOTOR":         {"PS Polyester Non Woven Spunbond 100 Gsm White"},
	"SPUNBOUND TERLIPAT":      {"PS Polyester Non Woven Spunbond 100 Gsm White"},
	"SPUNBOUND HARDEN":        {"PS Polyester Non Woven Spunbond 100 Gsm White"},
	"LAMINATING BOLONG":       {"Laminasi LDPE 200 Gsm"},
	"LAMINATING TIDAK MATANG": {"Laminasi LDPE 200 Gsm"},
	"SOBEK":                   {"Recycle Felt GWPS 2mm 375 Gsm", "Carpet STKD19 Black", "Carpet CBIII", "Ester Canvas SAB10-NS121 SSP", "Silincer T. 15mm 1000 Gsm", "Silincer T. 6mm 350 Gsm"},
	"BRUDUL":                  {"Recycle Felt GWPS 2mm 375 Gsm", "Carpet STKD19 Black", "Carpet CBIII", "Ester Canvas SAB10-NS121 SSP", "Silincer T. 15mm 1000 Gsm", "Silincer T. 6mm 350 Gsm"},
	"TIPIS":                   {"Recycle Felt GWPS 2mm 375 Gsm", "Carpet STKD19 Black", "Carpet CBIII", "EPDM"},
	"BERJAMUR":                {"Carpet STKD19 Black", "Carpet CBIII"},
	"GALER":                   {"Carpet STKD19 Black", "Carpet CBIII"},
	"DENT":                    {"Carpet STKD19 Black", "Carpet CBIII"},
	"TERLIPAT":                {"Carpet STKD19 Black", "Carpet CBIII"},
	"BELANG":                  {"Carpet STKD19 Black", "Carpet CBIII"},
	"BERLUBANG":               {"EPDM"},
	"KOTOR":                   {"Ester Canvas SAB10-NS121 SSP", "Silincer T. 15mm 1000 Gsm", "Silincer T. 6mm 350 Gsm"},
	"MENGEMBANG":              {"Silincer T. 15mm 1000 Gsm", "Silincer T. 6mm 350 Gsm"},
}

// LayananAnalitik menyediakan logika pelaporan.
type LayananAnalitik interface {
	DapatkanParetoBulanan(bulan, tahun int) ([]DTOParetoMetrik, error)
	DapatkanPareto(tanggalMulai, tanggalSelesai string, zonaLini string) ([]DTOParetoMetrik, error)
	LacakAkarMasalah(kodeCacat string, parentSKU string) ([]DTOLacakAkarMasalah, error)
}

type layananAnalitik struct {
	repo RepositoriAnalitik
	db   *gorm.DB
}

// KonstruksiLayananBaru membuat objek LayananAnalitik.
func KonstruksiLayananBaru(repo RepositoriAnalitik, db *gorm.DB) LayananAnalitik {
	return &layananAnalitik{repo: repo, db: db}
}

// DapatkanParetoBulanan mengorkestrasi penarikan agregasi Pareto.
func (l *layananAnalitik) DapatkanParetoBulanan(bulan, tahun int) ([]DTOParetoMetrik, error) {
	return l.repo.KalkulasiParetoBulanan(bulan, tahun)
}

// DapatkanPareto mengorkestrasi penarikan agregasi Pareto secara dinamis.
func (l *layananAnalitik) DapatkanPareto(tanggalMulai, tanggalSelesai string, zonaLini string) ([]DTOParetoMetrik, error) {
	return l.repo.KalkulasiPareto(tanggalMulai, tanggalSelesai, zonaLini)
}

// LacakAkarMasalah menelusuri defects produk hingga ke level bahan baku pembentuk & supplier
func (l *layananAnalitik) LacakAkarMasalah(kodeCacat string, parentSKU string) ([]DTOLacakAkarMasalah, error) {
	var hasil []DTOLacakAkarMasalah

	// 1. Dapatkan Finished Good (Parent) Material
	var parents []manufaktur.Material
	if parentSKU != "" {
		var p manufaktur.Material
		if err := l.db.Where("kode_sku = ?", parentSKU).First(&p).Error; err != nil {
			return nil, errors.New("parent_sku_tidak_ditemukan")
		}
		parents = append(parents, p)
	} else {
		// Dapatkan semua material yang bertindak sebagai Finished Good (memiliki BOM anak)
		var parentIDs []uint
		l.db.Model(&manufaktur.KomposisiMaterialBOM{}).Distinct("id_parent_material").Pluck("id_parent_material", &parentIDs)
		if len(parentIDs) > 0 {
			l.db.Where("id IN ?", parentIDs).Find(&parents)
		}
	}

	// Cek apakah defect merupakan defect proses internal (tidak berakar ke material)
	isProcessDefect := false
	processDefects := []string{
		"OVERCUTTING", "DIMENSI OUT STD", "TERBALIK", "HOLE T/A", "SEWING MIRING",
		"SEWING LONCAT", "SEWING NITIK", "KUNCIAN JEBOL", "SEWING PUTUS",
		"ALUR SERAT TDK SESUAI", "SEWING LONGGAR", "LANGKAH SEWING TIDAK SESUAI",
	}
	for _, pd := range processDefects {
		if strings.ToUpper(kodeCacat) == pd {
			isProcessDefect = true
			break
		}
	}

	if isProcessDefect {
		for _, p := range parents {
			hasil = append(hasil, DTOLacakAkarMasalah{
				KodeCacat:       kodeCacat,
				ParentSKU:       p.KodeSKU,
				ParentNama:      p.NamaMaterial,
				RawMaterialSKU:  "N/A",
				RawMaterialNama: "INTERNAL PROCESS DEFECT",
				PemasokNama:     "INTERNAL PRODUCTION LINE",
				PemasokKontak:   "internal@pgn-quality.co.id",
				KuantitasBOM:    0.0,
			})
		}
		return hasil, nil
	}

	// 2. Lakukan rekursi BOM Tracing untuk setiap finished good
	for _, p := range parents {
		visited := make(map[uint]bool)
		var boms []manufaktur.KomposisiMaterialBOM
		circularDetected := false

		// Dapatkan BOM direct children
		err := l.db.Preload("MaterialBaku.Pemasok").Where("id_parent_material = ?", p.ID).Find(&boms).Error
		if err != nil {
			continue
		}

		var traverse func(parentID uint)
		traverse = func(parentID uint) {
			if visited[parentID] {
				circularDetected = true
				return
			}
			visited[parentID] = true
			defer func() { visited[parentID] = false }()

			var subBoms []manufaktur.KomposisiMaterialBOM
			l.db.Preload("MaterialBaku.Pemasok").Where("id_parent_material = ?", parentID).Find(&subBoms)
			for _, sb := range subBoms {
				boms = append(boms, sb)
				traverse(sb.IDRawMaterial)
			}
		}

		traverse(p.ID)

		// 3. Cari material pembentuk yang sesuai dengan defect
		targetMaterials, isMapped := DefectToMaterialMapping[strings.ToUpper(kodeCacat)]
		for _, bomItem := range boms {
			matched := false
			if isMapped {
				for _, tm := range targetMaterials {
					if strings.Contains(strings.ToLower(bomItem.MaterialBaku.NamaMaterial), strings.ToLower(tm)) {
						matched = true
						break
					}
				}
			} else {
				// Jika tidak ada di map, cocokkan secara fallback
				matched = true
			}

			if matched {
				hasil = append(hasil, DTOLacakAkarMasalah{
					KodeCacat:                  kodeCacat,
					ParentSKU:                  p.KodeSKU,
					ParentNama:                 p.NamaMaterial,
					RawMaterialSKU:             bomItem.MaterialBaku.KodeSKU,
					RawMaterialNama:            bomItem.MaterialBaku.NamaMaterial,
					PemasokNama:                bomItem.MaterialBaku.Pemasok.NamaEntitas,
					PemasokKontak:              bomItem.MaterialBaku.Pemasok.Kontak,
					KuantitasBOM:               bomItem.ParameterKuantitasPembentuk,
					CircularDependencyDetected: circularDetected,
				})
			}
		}
	}

	return hasil, nil
}
