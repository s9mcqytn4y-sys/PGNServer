package analitik

import (
	"strings"
)

// Ambang dan aturan konstan untuk DSS
const (
	AmbangRasioNGWaspada = 2.0
	AmbangRasioNGKritis  = 5.0
	AmbangDominasiPareto = 50.0
	MinimumHariTrend     = 3
	PeriodeTrendHari     = 7
)

var processDefects = []string{
	"OVERCUTTING", "DIMENSI OUT STD", "TERBALIK", "HOLE T/A", "SEWING MIRING",
	"SEWING LONCAT", "SEWING NITIK", "KUNCIAN JEBOL", "SEWING PUTUS",
	"ALUR SERAT TDK SESUAI", "SEWING LONGGAR", "LANGKAH SEWING TIDAK SESUAI",
}

// NormalisasiKodeCacat mengubah kode cacat ke huruf besar untuk standarisasi pencocokan
func NormalisasiKodeCacat(kode string) string {
	return strings.ToUpper(strings.TrimSpace(kode))
}

// ApakahDefectMaterial mengembalikan true jika cacat terkait dengan material di DefectToMaterialMapping
func ApakahDefectMaterial(kode string) bool {
	kodeNormal := NormalisasiKodeCacat(kode)
	_, ada := DefectToMaterialMapping[kodeNormal]
	return ada
}

// ApakahDefectProcess mengembalikan true jika cacat merupakan kesalahan proses produksi internal
func ApakahDefectProcess(kode string) bool {
	kodeNormal := NormalisasiKodeCacat(kode)
	for _, pd := range processDefects {
		if kodeNormal == pd {
			return true
		}
	}
	return false
}
