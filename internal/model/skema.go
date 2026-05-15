package model

import (
	"time"

	"gorm.io/datatypes"
)

// Produk merepresentasikan data produk di aplikasi
type Produk struct {
	ID         string    `gorm:"primaryKey;column:id" json:"id"`
	NomorPart  string    `gorm:"column:nomor_part;not null" json:"nomor_part"`
	NomorUnik  string    `gorm:"column:nomor_unik;uniqueIndex" json:"nomor_unik"`
	NamaPart   string    `gorm:"column:nama_part;not null" json:"nama_part"`
	Model      string    `gorm:"column:model" json:"model"`
	CustomerID string    `gorm:"column:customer_id" json:"customer_id"`
	LineID     string    `gorm:"column:line_id" json:"line_id"`
	NamaAssy   string    `gorm:"column:nama_assy" json:"nama_assy"`
	LokasiFoto string    `gorm:"column:lokasi_foto" json:"lokasi_foto"`
	DibuatPada time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

// Material merepresentasikan data material/bahan baku
type Material struct {
	ID         string    `gorm:"primaryKey;column:id" json:"id"`
	SupplierID string    `gorm:"column:supplier_id" json:"supplier_id"`
	NomorUnik  string    `gorm:"column:nomor_unik" json:"nomor_unik"`
	NamaPart   string    `gorm:"column:nama_part;not null" json:"nama_part"`
	TebalMM    float64   `gorm:"column:tebal_mm" json:"tebal_mm"`
	LebarCM    float64   `gorm:"column:lebar_cm" json:"lebar_cm"`
	PanjangCM  float64   `gorm:"column:panjang_cm" json:"panjang_cm"`
	BeratGSM   float64   `gorm:"column:berat_gsm" json:"berat_gsm"`
	MassaKG    float64   `gorm:"column:massa_kg" json:"massa_kg"`
	Satuan     string    `gorm:"column:satuan" json:"satuan"`
	LokasiFoto string    `gorm:"column:lokasi_foto" json:"lokasi_foto"`
	DibuatPada time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

// BillOfMaterial merepresentasikan relasi produk dan material
type BillOfMaterial struct {
	ProdukID        string   `gorm:"primaryKey;column:produk_id" json:"produk_id"`
	MaterialID      string   `gorm:"primaryKey;column:material_id" json:"material_id"`
	JumlahPemakaian float64  `gorm:"column:jumlah_pemakaian" json:"jumlah_pemakaian"`
	Produk          Produk   `gorm:"foreignKey:ProdukID" json:"produk"`
	Material        Material `gorm:"foreignKey:MaterialID" json:"material"`
}

// KategoriDefect merepresentasikan kategori_defect
type KategoriDefect struct {
	ID           string    `gorm:"primaryKey;column:id" json:"id"`
	NamaKategori string    `gorm:"column:nama_kategori;not null" json:"nama_kategori"`
	DibuatPada   time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

// MasterDefect merepresentasikan master_defect
type MasterDefect struct {
	ID         string         `gorm:"primaryKey;column:id" json:"id"`
	KategoriID string         `gorm:"column:kategori_id" json:"kategori_id"`
	NamaNG     string         `gorm:"column:nama_ng;not null" json:"nama_ng"`
	DibuatPada time.Time      `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
	Kategori   KategoriDefect `gorm:"foreignKey:KategoriID" json:"kategori"`
}

// MasterMesin merepresentasikan master mesin untuk traceability
type MasterMesin struct {
	ID         string    `gorm:"primaryKey;column:id" json:"id"`
	NamaMesin  string    `gorm:"column:nama_mesin;not null" json:"nama_mesin"`
	TipeMesin  string    `gorm:"column:tipe_mesin" json:"tipe_mesin"` // e.g., Press, Cutting, Sewing
	DibuatPada time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

// InspeksiHarian merepresentasikan inspeksi_harian
type InspeksiHarian struct {
	ID              string    `gorm:"primaryKey;column:id" json:"id"`
	TanggalInspeksi time.Time `gorm:"column:tanggal_inspeksi;not null" json:"tanggal_inspeksi"`
	LeaderID        uint      `gorm:"column:leader_id" json:"leader_id"`
	LineID          string    `gorm:"column:line_id" json:"line_id"`
	ProdukID        string    `gorm:"column:produk_id" json:"produk_id"`
	TotalProduksi   int       `gorm:"column:total_produksi" json:"total_produksi"`
	TotalOK         int       `gorm:"column:total_ok" json:"total_ok"`
	TotalNG         int       `gorm:"column:total_ng" json:"total_ng"`
	StatTahun       int       `gorm:"column:stat_tahun" json:"stat_tahun"`
	StatBulan       int       `gorm:"column:stat_bulan" json:"stat_bulan"`
	StatMinggu      int       `gorm:"column:stat_minggu" json:"stat_minggu"`
	DibuatPada      time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
	Produk          Produk    `gorm:"foreignKey:ProdukID" json:"produk"`
}

// ChecksheetHeader merepresentasikan header dari checksheet transaksi
type ChecksheetHeader struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID     uint      `gorm:"column:user_id" json:"user_id"`
	LineID     string    `gorm:"column:line_id" json:"line_id"`
	ProdukID   string    `gorm:"column:produk_id" json:"produk_id"`
	MesinID    string    `gorm:"column:mesin_id" json:"mesin_id"`
	Shift      int       `gorm:"column:shift" json:"shift"`
	Status     string    `gorm:"column:status;default:'Open'" json:"status"`
	DibuatPada time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	Mesin      MasterMesin `gorm:"foreignKey:MesinID" json:"mesin"`
}

// ChecksheetDetail merepresentasikan detail parameter checksheet menggunakan JSONB
type ChecksheetDetail struct {
	ID             uint           `gorm:"primaryKey;column:id" json:"id"`
	HeaderID       uint           `gorm:"column:header_id" json:"header_id"`
	WaktuCek       time.Time      `gorm:"column:waktu_cek" json:"waktu_cek"`
	Proses         string         `gorm:"column:proses" json:"proses"` // Material Cutting, Part Press, etc.
	ParameterHasil datatypes.JSON `gorm:"column:parameter_hasil" json:"parameter_hasil"`
	Keterangan     string         `gorm:"column:keterangan" json:"keterangan"`
	Header         ChecksheetHeader `gorm:"foreignKey:HeaderID" json:"header"`
}

// LogInspeksi merepresentasikan log_inspeksi
type LogInspeksi struct {
	ID            uint           `gorm:"primaryKey;column:id" json:"id"`
	InspeksiID    string         `gorm:"column:inspeksi_id" json:"inspeksi_id"`
	DefectID      string         `gorm:"column:defect_id" json:"defect_id"`
	JendelaWaktu  string         `gorm:"column:jendela_waktu;not null" json:"jendela_waktu"`
	WaktuKejadian time.Time      `gorm:"column:waktu_kejadian;not null" json:"waktu_kejadian"`
	JumlahNG      int            `gorm:"column:jumlah_ng;not null" json:"jumlah_ng"`
	DibuatPada    time.Time      `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
	Inspeksi      InspeksiHarian `gorm:"foreignKey:InspeksiID" json:"inspeksi"`
	Defect        MasterDefect   `gorm:"foreignKey:DefectID" json:"defect"`
}

// BukuBesarDefectMaterial merepresentasikan buku_besar_defect_material
type BukuBesarDefectMaterial struct {
	ID            uint        `gorm:"primaryKey;column:id" json:"id"`
	LogInspeksiID uint        `gorm:"column:log_inspeksi_id" json:"log_inspeksi_id"`
	MaterialID    string      `gorm:"column:material_id" json:"material_id"`
	DefectID      string      `gorm:"column:defect_id" json:"defect_id"`
	JumlahNG      int         `gorm:"column:jumlah_ng;not null" json:"jumlah_ng"`
	DicatatPada   time.Time   `gorm:"column:dicatat_pada;autoCreateTime" json:"dicatat_pada"`
	LogInspeksi   LogInspeksi `gorm:"foreignKey:LogInspeksiID" json:"log_inspeksi"`
	Material      Material    `gorm:"foreignKey:MaterialID" json:"material"`
}

// LiniProduksi merepresentasikan lini_produksi
type LiniProduksi struct {
	ID         string    `gorm:"primaryKey;column:id" json:"id"`
	NamaLini   string    `gorm:"column:nama_lini;not null" json:"nama_lini"`
	DibuatPada time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

// Supplier merepresentasikan supplier
type Supplier struct {
	ID           string    `gorm:"primaryKey;column:id" json:"id"`
	NamaSupplier string    `gorm:"column:nama_supplier;not null" json:"nama_supplier"`
	DibuatPada   time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

// Customer merepresentasikan customer
type Customer struct {
	ID           string    `gorm:"primaryKey;column:id" json:"id"`
	NamaCustomer string    `gorm:"column:nama_customer;not null" json:"nama_customer"`
	DibuatPada   time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

// User merepresentasikan users
type User struct {
	ID            uint      `gorm:"primaryKey;column:id" json:"id"`
	NIP           string    `gorm:"column:nip;unique;not null" json:"nip"`
	Password      string    `gorm:"column:password;not null" json:"-"`
	Nama          string    `gorm:"column:name;not null" json:"nama"`
	Role          string    `gorm:"column:role;default:'Leader QC'" json:"role"`
	RememberToken string    `gorm:"column:remember_token" json:"-"`
	DibuatPada    time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
	DiubahPada    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"diubah_pada"`
}
