// Package analitik menangani agregasi dan kalkulasi data pelaporan 7 QC Tools.
package analitik

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"pgn-server/internal/manufaktur"
	"pgn-server/pkg/cache"
	"pgn-server/pkg/pencatatan_log"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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
	DapatkanParetoBulanan(k *gin.Context, bulan, tahun int) ([]DTOParetoMetrik, error)
	DapatkanPareto(k *gin.Context, tanggalMulai, tanggalSelesai string, zonaLini string) ([]DTOParetoMetrik, error)
	LacakAkarMasalah(k *gin.Context, kodeCacat string, parentSKU string) ([]DTOLacakAkarMasalah, error)
}

type layananAnalitik struct {
	repo RepositoriAnalitik
	db   *gorm.DB
}

// KonstruksiLayananBaru membuat objek LayananAnalitik.
func KonstruksiLayananBaru(repo RepositoriAnalitik, db *gorm.DB) LayananAnalitik {
	return &layananAnalitik{repo: repo, db: db}
}

// DapatkanParetoBulanan mengorkestrasi penarikan agregasi Pareto dengan caching dan logging performa.
func (l *layananAnalitik) DapatkanParetoBulanan(k *gin.Context, bulan, tahun int) ([]DTOParetoMetrik, error) {
	waktuMulai := time.Now()
	cacheKey := fmt.Sprintf("pareto_bulanan_%d_%d", bulan, tahun)

	if cachedVal, found := cache.GlobalCache.Get(cacheKey); found {
		pencatatan_log.Info(k, "[Telemetry] DapatkanParetoBulanan Cache HIT untuk kunci: %s, Latency: %v", cacheKey, time.Since(waktuMulai))
		return cachedVal.([]DTOParetoMetrik), nil
	}

	hasil, err := l.repo.KalkulasiParetoBulanan(bulan, tahun)
	if err != nil {
		pencatatan_log.Galat(k, "Kesalahan KalkulasiParetoBulanan: %v", err)
		return nil, err
	}

	cache.GlobalCache.Set(cacheKey, hasil, 5*time.Minute)
	pencatatan_log.Info(k, "[Telemetry] DapatkanParetoBulanan Cache MISS untuk kunci: %s, Latency: %v", cacheKey, time.Since(waktuMulai))
	return hasil, nil
}

// DapatkanPareto mengorkestrasi penarikan agregasi Pareto secara dinamis dengan caching dan logging performa.
func (l *layananAnalitik) DapatkanPareto(k *gin.Context, tanggalMulai, tanggalSelesai string, zonaLini string) ([]DTOParetoMetrik, error) {
	waktuMulai := time.Now()
	cacheKey := fmt.Sprintf("pareto_%s_%s_%s", tanggalMulai, tanggalSelesai, zonaLini)

	if cachedVal, found := cache.GlobalCache.Get(cacheKey); found {
		pencatatan_log.Info(k, "[Telemetry] DapatkanPareto Cache HIT untuk kunci: %s, Latency: %v", cacheKey, time.Since(waktuMulai))
		return cachedVal.([]DTOParetoMetrik), nil
	}

	hasil, err := l.repo.KalkulasiPareto(tanggalMulai, tanggalSelesai, zonaLini)
	if err != nil {
		pencatatan_log.Galat(k, "Kesalahan KalkulasiPareto: %v", err)
		return nil, err
	}

	cache.GlobalCache.Set(cacheKey, hasil, 5*time.Minute)
	pencatatan_log.Info(k, "[Telemetry] DapatkanPareto Cache MISS untuk kunci: %s, Latency: %v", cacheKey, time.Since(waktuMulai))
	return hasil, nil
}

// LacakAkarMasalah menelusuri defects produk hingga ke level bahan baku pembentuk & supplier
func (l *layananAnalitik) LacakAkarMasalah(k *gin.Context, kodeCacat string, parentSKU string) ([]DTOLacakAkarMasalah, error) {
	waktuMulai := time.Now()
	var hasil []DTOLacakAkarMasalah

	// 1. Dapatkan Finished Good (Parent) Material
	var parents []manufaktur.Material
	var cacheHit = 0
	var cacheMiss = 0

	if parentSKU != "" {
		cacheKey := "material_sku_" + parentSKU
		if cachedVal, found := cache.GlobalCache.Get(cacheKey); found {
			parents = append(parents, cachedVal.(manufaktur.Material))
			cacheHit++
		} else {
			var p manufaktur.Material
			if err := l.db.Where("kode_sku = ?", parentSKU).First(&p).Error; err != nil {
				return nil, errors.New("parent_sku_tidak_ditemukan")
			}
			parents = append(parents, p)
			cache.GlobalCache.Set(cacheKey, p, 5*time.Minute)
			cacheMiss++
		}
	} else {
		cacheKey := "bom_parent_materials_all"
		if cachedVal, found := cache.GlobalCache.Get(cacheKey); found {
			parents = cachedVal.([]manufaktur.Material)
			cacheHit++
		} else {
			var parentIDs []uint
			l.db.Model(&manufaktur.KomposisiMaterialBOM{}).Distinct("id_parent_material").Pluck("id_parent_material", &parentIDs)
			if len(parentIDs) > 0 {
				l.db.Where("id IN ?", parentIDs).Find(&parents)
			}
			cache.GlobalCache.Set(cacheKey, parents, 5*time.Minute)
			cacheMiss++
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
		pencatatan_log.Info(k, "[Telemetry] LacakAkarMasalah (Internal Process Defect): %s, Latency: %v, Parents: %d", kodeCacat, time.Since(waktuMulai), len(parents))
		return hasil, nil
	}

	// 2. Lakukan rekursi BOM Tracing untuk setiap finished good secara konkuren via errgroup
	var mu sync.Mutex
	var localCacheHit int32
	var localCacheMiss int32

	ctx := context.Background()
	if k != nil && k.Request != nil {
		ctx = k.Request.Context()
	}
	g, gCtx := errgroup.WithContext(ctx)

	for _, parentVal := range parents {
		p := parentVal // Pinned loop variable
		g.Go(func() error {
			visited := make(map[uint]bool)
			var boms []manufaktur.KomposisiMaterialBOM
			circularDetected := false

			// Dapatkan BOM direct children (caching)
			directCacheKey := fmt.Sprintf("bom_by_parent_%d", p.ID)
			if cachedBoms, found := cache.GlobalCache.Get(directCacheKey); found {
				boms = cachedBoms.([]manufaktur.KomposisiMaterialBOM)
				localCacheHit++
			} else {
				var dbBoms []manufaktur.KomposisiMaterialBOM
				err := l.db.WithContext(gCtx).Preload("MaterialBaku.Pemasok").Where("id_parent_material = ?", p.ID).Find(&dbBoms).Error
				if err == nil {
					boms = dbBoms
					cache.GlobalCache.Set(directCacheKey, dbBoms, 5*time.Minute)
				}
				localCacheMiss++
			}

			var traverse func(parentID uint)
			traverse = func(parentID uint) {
				if visited[parentID] {
					circularDetected = true
					return
				}
				visited[parentID] = true
				defer func() { visited[parentID] = false }()

				cacheKey := fmt.Sprintf("bom_by_parent_%d", parentID)
				var subBoms []manufaktur.KomposisiMaterialBOM
				if cachedSubBoms, found := cache.GlobalCache.Get(cacheKey); found {
					subBoms = cachedSubBoms.([]manufaktur.KomposisiMaterialBOM)
					localCacheHit++
				} else {
					l.db.WithContext(gCtx).Preload("MaterialBaku.Pemasok").Where("id_parent_material = ?", parentID).Find(&subBoms)
					cache.GlobalCache.Set(cacheKey, subBoms, 5*time.Minute)
					localCacheMiss++
				}

				for _, sb := range subBoms {
					boms = append(boms, sb)
					traverse(sb.IDRawMaterial)
				}
			}

			traverse(p.ID)

			// Cari material pembentuk yang sesuai dengan defect
			targetMaterials, isMapped := DefectToMaterialMapping[strings.ToUpper(kodeCacat)]
			var localHasil []DTOLacakAkarMasalah

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
					matched = true
				}

				if matched {
					localHasil = append(localHasil, DTOLacakAkarMasalah{
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

			mu.Lock()
			hasil = append(hasil, localHasil...)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	totalCacheHit := cacheHit + int(localCacheHit)
	totalCacheMiss := cacheMiss + int(localCacheMiss)
	totalAccesses := totalCacheHit + totalCacheMiss
	hitRatio := 0.0
	if totalAccesses > 0 {
		hitRatio = float64(totalCacheHit) / float64(totalAccesses) * 100
	}

	pencatatan_log.Info(k, "[Telemetry] LacakAkarMasalah: %s, Latency: %v, Worker Pool: %d concurrent tasks, Cache Hit Ratio: %.2f%% (Hit: %d, Miss: %d)",
		kodeCacat, time.Since(waktuMulai), len(parents), hitRatio, totalCacheHit, totalCacheMiss)

	return hasil, nil
}
