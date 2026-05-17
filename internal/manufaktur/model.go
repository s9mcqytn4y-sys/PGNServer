package manufaktur

import "time"

// Customer merepresentasikan entitas pelanggan (Customer).
type Customer struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	CustomerCode string    `gorm:"column:customer_code;type:varchar(100);uniqueIndex;not null" json:"customerCode"`
	Nama         string    `gorm:"column:nama;type:varchar(255);not null" json:"nama"`
	Kontak       string    `gorm:"column:kontak;type:varchar(255)" json:"kontak"`
	DibuatPada   time.Time `gorm:"autoCreateTime" json:"dibuatPada"`
	DiubahPada   time.Time `gorm:"autoUpdateTime" json:"diubahPada"`
}

func (Customer) TableName() string {
	return "customers"
}

// Pemasok merepresentasikan entitas pemasok bahan baku.
type Pemasok struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	SupplierCode string    `gorm:"column:supplier_code;type:varchar(100);not null;default:'SUP-000'" json:"supplierCode"`
	NamaEntitas  string    `gorm:"column:nama_entitas;type:varchar(255);not null" json:"namaEntitas"`
	Kontak       string    `gorm:"column:kontak;type:varchar(255)" json:"kontak"`
	DibuatPada   time.Time `gorm:"autoCreateTime" json:"dibuatPada"`
	DiubahPada   time.Time `gorm:"autoUpdateTime" json:"diubahPada"`
}

func (Pemasok) TableName() string {
	return "suppliers"
}

// Material merepresentasikan bahan baku yang disuplai oleh Pemasok.
type Material struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	KodeSKU      string    `gorm:"column:kode_sku;type:varchar(100);uniqueIndex;not null" json:"kodeSKU"`
	NamaMaterial string    `gorm:"column:nama_material;type:varchar(255);not null" json:"namaMaterial"`
	TebalCM      float64   `gorm:"column:tebal_cm" json:"tebalCM"`
	BeratGSM     int       `gorm:"column:berat_gsm;type:int" json:"beratGSM"`
	LebarCM      float64   `gorm:"column:lebar_cm" json:"lebarCM"`
	PanjangCM    float64   `gorm:"column:panjang_cm" json:"panjangCM"`
	UnitSatuan   string    `gorm:"column:unit_satuan;type:varchar(50)" json:"unitSatuan"`
	IDPemasok    uint      `gorm:"column:id_pemasok;index;not null" json:"idPemasok"`
	Pemasok      Pemasok   `gorm:"foreignKey:IDPemasok;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"pemasok,omitempty"`
	DibuatPada   time.Time `gorm:"autoCreateTime" json:"dibuatPada"`
	DiubahPada   time.Time `gorm:"autoUpdateTime" json:"diubahPada"`
}

func (Material) TableName() string {
	return "MATERIAL"
}

// KomposisiMaterialBOM merepresentasikan relasi kuantitatif hirarkis dalam Bill of Materials.
type KomposisiMaterialBOM struct {
	ID                          uint      `gorm:"primaryKey;column:id" json:"id"`
	IDParentMaterial            *uint     `gorm:"column:id_parent_material;index:idx_bom_parent_raw" json:"idParentMaterial"`
	ParentMaterial              *Material `gorm:"foreignKey:IDParentMaterial;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"parentMaterial,omitempty"`
	IDRawMaterial               uint      `gorm:"column:id_raw_material;index:idx_bom_parent_raw;not null" json:"idRawMaterial"`
	MaterialBaku                Material  `gorm:"foreignKey:IDRawMaterial;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"materialBaku,omitempty"`
	ParameterKuantitasPembentuk float64   `gorm:"column:parameter_kuantitas_pembentuk;not null" json:"parameterKuantitasPembentuk"`
	DibuatPada                  time.Time `gorm:"autoCreateTime" json:"dibuatPada"`
}

func (KomposisiMaterialBOM) TableName() string {
	return "bill_of_materials"
}

// LineProduksiSnapshotDto merepresentasikan DTO lini produksi untuk snapshot.
type LineProduksiSnapshotDto struct {
	ID           string `json:"id"`
	KodeLine     string `json:"kodeLine"`
	NamaLine     string `json:"namaLine"`
	Aktif        bool   `json:"aktif"`
	UrutanTampil int    `json:"urutanTampil"`
}

// SlotWaktuSnapshotDto merepresentasikan DTO slot waktu untuk snapshot.
type SlotWaktuSnapshotDto struct {
	ID           string  `json:"id"`
	KodeSlot     string  `json:"kodeSlot"`
	LabelSlot    string  `json:"labelSlot"`
	JamMulai     *string `json:"jamMulai"`
	JamSelesai   *string `json:"jamSelesai"`
	Aktif        bool    `json:"aktif"`
	UrutanTampil int     `json:"urutanTampil"`
}

// MaterialSnapshotDto merepresentasikan DTO material untuk snapshot.
type MaterialSnapshotDto struct {
	ID           string `json:"id"`
	KodeSKU      string `json:"kodeSKU"`
	NamaMaterial string `json:"namaMaterial"`
	Aktif        bool   `json:"aktif"`
}

// PartSnapshotDto merepresentasikan DTO part untuk snapshot.
type PartSnapshotDto struct {
	ID                  string  `json:"id"`
	KodeUnikPart        string  `json:"kodeUnikPart"`
	NamaPart            string  `json:"namaPart"`
	NomorPart           *string `json:"nomorPart"`
	MaterialID          *string `json:"materialId"`
	KodeMaterial        *string `json:"kodeMaterial"`
	NamaMaterial        *string `json:"namaMaterial"`
	KodeProyek          *string `json:"kodeProyek"`
	JumlahItemPerKanban *int    `json:"jumlahItemPerKanban"`
	LineDefaultID       *string `json:"lineDefaultId"`
	KodeLineDefault     *string `json:"kodeLineDefault"`
	NamaLineDefault     *string `json:"namaLineDefault"`
	Aktif               bool    `json:"aktif"`
	SumberData          *string `json:"sumberData"`
}

// KategoriDefectSnapshotDto merepresentasikan DTO kategori defect untuk snapshot.
type KategoriDefectSnapshotDto struct {
	ID           string `json:"id"`
	KodeKategori string `json:"kodeKategori"`
	NamaKategori string `json:"namaKategori"`
	Aktif        bool   `json:"aktif"`
	UrutanTampil int    `json:"urutanTampil"`
}

// JenisDefectSnapshotDto merepresentasikan DTO jenis defect untuk snapshot.
type JenisDefectSnapshotDto struct {
	ID               string  `json:"id"`
	KodeDefect       string  `json:"kodeDefect"`
	NamaDefect       string  `json:"namaDefect"`
	KategoriDefectID *string `json:"kategoriDefectId"`
	KodeKategori     *string `json:"kodeKategori"`
	NamaKategori     *string `json:"namaKategori"`
	Aktif            bool    `json:"aktif"`
}

// RelasiPartDefectSnapshotDto merepresentasikan DTO relasi part dan defect untuk snapshot.
type RelasiPartDefectSnapshotDto struct {
	ID                 string  `json:"id"`
	PartID             string  `json:"partId"`
	KodeUnikPart       *string `json:"kodeUnikPart"`
	JenisDefectID      string  `json:"jenisDefectId"`
	KodeDefect         *string `json:"kodeDefect"`
	KodeTampilanDefect *string `json:"kodeTampilanDefect"`
	UrutanTampil       int     `json:"urutanTampil"`
	Aktif              bool    `json:"aktif"`
}

// MasterDataSnapshotDto mengelompokkan seluruh entitas master data untuk sinkronisasi offline-first.
type MasterDataSnapshotDto struct {
	VersiMasterData  string                        `json:"versiMasterData"`
	LineProduksi     []LineProduksiSnapshotDto     `json:"lineProduksi"`
	SlotWaktu        []SlotWaktuSnapshotDto        `json:"slotWaktu"`
	Material         []MaterialSnapshotDto         `json:"material"`
	Part             []PartSnapshotDto             `json:"part"`
	KategoriDefect   []KategoriDefectSnapshotDto   `json:"kategoriDefect"`
	JenisDefect      []JenisDefectSnapshotDto      `json:"jenisDefect"`
	RelasiPartDefect []RelasiPartDefectSnapshotDto `json:"relasiPartDefect"`
}

// MetadataSnapshotDto menyajikan statistik master data untuk QControl.
type MetadataSnapshotDto struct {
	JumlahLineProduksi     int `json:"jumlahLineProduksi"`
	JumlahSlotWaktu        int `json:"jumlahSlotWaktu"`
	JumlahMaterial         int `json:"jumlahMaterial"`
	JumlahPart             int `json:"jumlahPart"`
	JumlahJenisDefect      int `json:"jumlahJenisDefect"`
	JumlahRelasiPartDefect int `json:"jumlahRelasiPartDefect"`
	JumlahShiftOperasional int `json:"jumlahShiftOperasional"`
}
