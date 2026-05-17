// Package analitik menangani agregasi dan kalkulasi data pelaporan 7 QC Tools.
package analitik

// DTOParetoMetrik mewakili struktur agregasi cacat kumulatif.
type DTOParetoMetrik struct {
	KodeCacat           string  `json:"kode_cacat"`
	JumlahCacat         float64 `json:"jumlah_cacat"`
	Persentase          float64 `json:"persentase"`
	PersentaseKumulatif float64 `json:"persentase_kumulatif"`
}

// DTOLacakAkarMasalah mewakili visualisasi pelacakan dari cacat ke raw material pembentuk dan pemasoknya.
type DTOLacakAkarMasalah struct {
	KodeCacat              string  `json:"kode_cacat"`
	ParentSKU              string  `json:"parent_sku"`
	ParentNama             string  `json:"parent_nama"`
	RawMaterialSKU         string  `json:"raw_material_sku"`
	RawMaterialNama        string  `json:"raw_material_nama"`
	PemasokNama            string  `json:"pemasok_nama"`
	PemasokKontak          string  `json:"pemasok_kontak"`
	KuantitasBOM           float64 `json:"kuantitas_bom"`
	CircularDependencyDetected bool `json:"circular_dependency_detected"`
}
