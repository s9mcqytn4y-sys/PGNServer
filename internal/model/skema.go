package model

import (
	"time"

	"gorm.io/datatypes"
)

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

func (User) TableName() string {
	return "users"
}

// Customer merepresentasikan customer
type Customer struct {
	ID           string    `gorm:"type:text;primaryKey;column:id" json:"id"`
	NamaCustomer string    `gorm:"column:nama_customer;not null" json:"nama_customer"`
	DibuatPada   time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

func (Customer) TableName() string {
	return "customers"
}

// Supplier merepresentasikan supplier
type Supplier struct {
	ID           string    `gorm:"type:text;primaryKey;column:id" json:"id"`
	NamaSupplier string    `gorm:"column:nama_supplier;not null" json:"nama_supplier"`
	DibuatPada   time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

func (Supplier) TableName() string {
	return "suppliers"
}

// LiniProduksi merepresentasikan lini_produksi
type LiniProduksi struct {
	ID         string    `gorm:"type:text;primaryKey;column:id" json:"id"`
	NamaLini   string    `gorm:"column:nama_lini;not null" json:"nama_lini"`
	DibuatPada time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

func (LiniProduksi) TableName() string {
	return "production_lines"
}

// Material merepresentasikan data material/bahan baku
type Material struct {
	ID            string         `gorm:"type:text;primaryKey;column:id" json:"id"`
	SupplierID    string         `gorm:"column:supplier_id" json:"supplier_id"`
	NomorUnik     string         `gorm:"column:nomor_unik" json:"nomor_unik"`
	NamaPart      string         `gorm:"column:nama_part;not null" json:"nama_part"`
	Satuan        string         `gorm:"column:satuan" json:"satuan"`
	PotensiDefect []DefectMaster `gorm:"many2many:material_potential_defects;joinForeignKey:MATERIAL_ID;joinReferences:DEFECT_ID" json:"potensi_defect"`
	DibuatPada    time.Time      `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

func (Material) TableName() string {
	return "MATERIAL"
}

// Produk merepresentasikan data produk di aplikasi
type Produk struct {
	ID         string           `gorm:"type:text;primaryKey;column:id" json:"id"`
	NomorPart  string           `gorm:"column:nomor_part;not null" json:"nomor_part"`
	NomorUnik  string           `gorm:"column:nomor_unik;uniqueIndex" json:"nomor_unik"`
	NamaPart   string           `gorm:"column:nama_part;not null" json:"nama_part"`
	Model      string           `gorm:"column:model" json:"model"`
	CustomerID string           `gorm:"column:customer_id" json:"customer_id"`
	LineID     string           `gorm:"column:line_id" json:"line_id"`
	AssyName   *string          `gorm:"column:ASSY_NAME" json:"assy_name"`
	BOM        []BillOfMaterial `gorm:"foreignKey:ProdukID" json:"bom"`
	DibuatPada time.Time        `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

func (Produk) TableName() string {
	return "products"
}

// BillOfMaterial merepresentasikan relasi produk dan material
type BillOfMaterial struct {
	ProdukID   string   `gorm:"type:text;primaryKey;column:produk_id" json:"produk_id"`
	MaterialID string   `gorm:"type:text;primaryKey;column:material_id" json:"material_id"`
	UsageQty   float64  `gorm:"column:usage_qty" json:"usage_qty"`
	Material   Material `gorm:"foreignKey:MaterialID;references:ID" json:"material"`
}

func (BillOfMaterial) TableName() string {
	return "bill_of_materials"
}

// DefectMaster merepresentasikan master_defect
type DefectMaster struct {
	ID         string    `gorm:"type:text;primaryKey;column:id" json:"id"`
	Kategori   string    `gorm:"column:kategori;not null" json:"kategori"` // MATERIAL or PROCESS
	NamaNG     string    `gorm:"column:nama_ng;not null" json:"nama_ng"`
	DibuatPada time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

func (DefectMaster) TableName() string {
	return "defect_master"
}

// InspeksiHarian merepresentasikan daily_inspections
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
}

func (InspeksiHarian) TableName() string {
	return "daily_inspections"
}

// LogInspeksi merepresentasikan inspection_logs
type LogInspeksi struct {
	ID            uint      `gorm:"primaryKey;column:id" json:"id"`
	InspeksiID    string    `gorm:"column:inspeksi_id" json:"inspeksi_id"`
	DefectID      string    `gorm:"column:defect_id" json:"defect_id"`
	JendelaWaktu  string    `gorm:"column:jendela_waktu;not null" json:"jendela_waktu"`
	WaktuKejadian time.Time `gorm:"column:waktu_kejadian;not null" json:"waktu_kejadian"`
	JumlahNG      int       `gorm:"column:jumlah_ng;not null" json:"jumlah_ng"`
	DibuatPada    time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

func (LogInspeksi) TableName() string {
	return "inspection_logs"
}

// BukuBesarDefectMaterial merepresentasikan material_defect_ledger
type BukuBesarDefectMaterial struct {
	ID            uint      `gorm:"primaryKey;column:id" json:"id"`
	LogInspeksiID uint      `gorm:"column:log_inspeksi_id" json:"log_inspeksi_id"`
	MaterialID    string    `gorm:"column:material_id" json:"material_id"`
	DefectID      string    `gorm:"column:defect_id" json:"defect_id"`
	JumlahNG      int       `gorm:"column:jumlah_ng;not null" json:"jumlah_ng"`
	DicatatPada   time.Time `gorm:"column:dicatat_pada;autoCreateTime" json:"dicatat_pada"`
}

func (BukuBesarDefectMaterial) TableName() string {
	return "material_defect_ledger"
}

// --- Tambahan model untuk checksheet ---

type ChecksheetHeader struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID     uint      `gorm:"column:user_id" json:"user_id"`
	LineID     string    `gorm:"column:line_id" json:"line_id"`
	ProdukID   string    `gorm:"column:produk_id" json:"produk_id"`
	Shift      int       `gorm:"column:shift" json:"shift"`
	Status     string    `gorm:"column:status;default:'Open'" json:"status"`
	DibuatPada time.Time `gorm:"column:created_at;autoCreateTime" json:"dibuat_pada"`
}

type ChecksheetDetail struct {
	ID             uint           `gorm:"primaryKey;column:id" json:"id"`
	HeaderID       uint           `gorm:"column:header_id" json:"header_id"`
	WaktuCek       time.Time      `gorm:"column:waktu_cek" json:"waktu_cek"`
	Proses         string         `gorm:"column:proses" json:"proses"`
	ParameterHasil datatypes.JSON `gorm:"column:parameter_hasil" json:"parameter_hasil"`
	Keterangan     string         `gorm:"column:keterangan" json:"keterangan"`
}
