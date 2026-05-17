// Package analitik menangani agregasi dan kalkulasi data pelaporan 7 QC Tools.
package analitik

// DTOParetoMetrik mewakili struktur agregasi cacat kumulatif bulanan.
type DTOParetoMetrik struct {
	KodeCacat           string  `json:"kode_cacat"`
	JumlahCacat         int64   `json:"jumlah_cacat"`
	Persentase          float64 `json:"persentase"`
	PersentaseKumulatif float64 `json:"persentase_kumulatif"`
}
