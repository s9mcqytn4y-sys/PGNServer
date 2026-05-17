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

// DTORingkasanNG mewakili ringkasan metrik kualitas (OK vs NG).
type DTORingkasanNG struct {
	TotalProduksi float64 `json:"total_produksi" example:"1000"` // Total barang diproduksi
	TotalOK       float64 `json:"total_ok" example:"950"`        // Total barang OK (Lolos QC)
	TotalDefect   float64 `json:"total_defect" example:"50"`     // Total barang NG (Reject)
	RasioNG       float64 `json:"rasio_ng" example:"5.0"`        // Persentase NG (%)
	JumlahHari    int     `json:"jumlah_hari" example:"5"`       // Jumlah hari aktif produksi
}

// DTOHistogramDefect mewakili frekuensi defect berdasarkan grup tertentu.
type DTOHistogramDefect struct {
	Kategori string  `json:"kategori" example:"08:00-12:00"` // Kategori grouping (contoh: slot waktu, tanggal, kode)
	Jumlah   float64 `json:"jumlah" example:"15"`            // Frekuensi/jumlah defect pada kategori tersebut
}

// DTOTrendDefect mewakili data deret waktu untuk run chart.
type DTOTrendDefect struct {
	Periode      string  `json:"periode" example:"2026-05-01"` // Waktu periode (harian, mingguan, bulanan)
	JumlahDefect float64 `json:"jumlah_defect" example:"12"`   // Frekuensi defect
}

// DTOStratifikasiDefect mewakili pemilahan data defect ke dalam sub-kategori.
type DTOStratifikasiDefect struct {
	Dimensi  string  `json:"dimensi" example:"kode_cacat"` // Dasar stratifikasi
	Kategori string  `json:"kategori" example:"A"`         // Nilai kategori stratifikasi
	Jumlah   float64 `json:"jumlah" example:"35"`          // Jumlah defect
}

// DTOSinyalKualitas mewakili status kendali (Control Signal).
type DTOSinyalKualitas struct {
	Status    string  `json:"status" example:"STABIL"`               // KRITIS, WASPADA, atau STABIL
	Alasan    string  `json:"alasan" example:"Rasio NG di bawah 2%"` // Penjelasan status
	Indikator float64 `json:"indikator" example:"1.5"`               // Metrik utama yang diukur (contoh: rasio NG)
}

// DTORekomendasiTindakan mewakili response dari Decision Support System (DSS)
type DTORekomendasiTindakan struct {
	Status      string                          `json:"status" example:"WASPADA"`
	Ringkasan   string                          `json:"ringkasan" example:"Rasio NG berada pada level waspada."`
	Indikator   DTOIndikatorRekomendasiTindakan `json:"indikator"`
	Rekomendasi []DTODetailRekomendasiTindakan  `json:"rekomendasi"`
}

// DTOIndikatorRekomendasiTindakan mewakili metrik dasar untuk rekomendasi DSS
type DTOIndikatorRekomendasiTindakan struct {
	TotalProduksi float64 `json:"total_produksi" example:"1000"`
	TotalOK       float64 `json:"total_ok" example:"950"`
	TotalDefect   float64 `json:"total_defect" example:"50"`
	RasioNG       float64 `json:"rasio_ng" example:"5.0"`
	DefectDominan string  `json:"defect_dominan" example:"LAMINATING BOLONG"`
	Trend7Hari    string  `json:"trend_7_hari" example:"NAIK"`
}

// DTODetailRekomendasiTindakan mewakili satu item tindakan saran dari DSS
type DTODetailRekomendasiTindakan struct {
	Target    string `json:"target" example:"QA"`
	Prioritas string `json:"prioritas" example:"TINGGI"`
	Tindakan  string `json:"tindakan" example:"Lakukan investigasi defect dominan pada line terkait."`
	Alasan    string `json:"alasan" example:"Rasio NG melebihi ambang kritis."`
}
