// Package analitik menangani agregasi dan kalkulasi data pelaporan 7 QC Tools.
package analitik

// DTOParetoMetrik mewakili struktur agregasi cacat kumulatif.
type DTOParetoMetrik struct {
	KodeCacat           string  `json:"kode_cacat" example:"A"`              // Kode cacat defect (misal: "A", "B", dll)
	JumlahCacat         float64 `json:"jumlah_cacat" example:"45"`           // Jumlah cacat kumulatif yang terdeteksi
	Persentase          float64 `json:"persentase" example:"60.0"`           // Persentase cacat spesifik terhadap total defect (%)
	PersentaseKumulatif float64 `json:"persentase_kumulatif" example:"60.0"` // Persentase kumulatif untuk analisis pareto 80/20 (%)
}

// DTOLacakAkarMasalah mewakili visualisasi pelacakan dari cacat ke raw material pembentuk dan pemasoknya.
type DTOLacakAkarMasalah struct {
	KodeCacat                  string  `json:"kode_cacat" example:"LAMINATING BOLONG"`            // Kode cacat defect yang ditelusuri
	ParentSKU                  string  `json:"parent_sku" example:"FG-002"`                       // SKU Produk Jadi (Finished Good)
	ParentNama                 string  `json:"parent_nama" example:"Tas Spunbound Premium"`       // Nama Produk Jadi (Finished Good)
	RawMaterialSKU             string  `json:"raw_material_sku" example:"RAW-001"`                // SKU Bahan Baku Pembentuk
	RawMaterialNama            string  `json:"raw_material_nama" example:"Bahan Spunbound Merah"` // Nama Bahan Baku Pembentuk
	PemasokNama                string  `json:"pemasok_nama" example:"CV Sumber Tekstil"`          // Nama Supplier / Pemasok Bahan Baku
	PemasokKontak              string  `json:"pemasok_kontak" example:"+62-812-3456-7890"`        // Informasi kontak Supplier / Pemasok
	KuantitasBOM               float64 `json:"kuantitas_bom" example:"0.12"`                      // Proporsi pemakaian bahan baku dalam BOM
	CircularDependencyDetected bool    `json:"circular_dependency_detected" example:"false"`      // Status deteksi dependensi melingkar dalam BOM
}
