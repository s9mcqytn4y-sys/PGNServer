package manufaktur

import "time"

// Pemasok merepresentasikan entitas pemasok bahan baku.
type Pemasok struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	NamaEntitas string    `gorm:"column:nama_entitas;type:varchar(255);not null" json:"namaEntitas"`
	Kontak      string    `gorm:"column:kontak;type:varchar(255)" json:"kontak"`
	DibuatPada  time.Time `gorm:"autoCreateTime" json:"dibuatPada"`
	DiubahPada  time.Time `gorm:"autoUpdateTime" json:"diubahPada"`
}

// Material merepresentasikan bahan baku yang disuplai oleh Pemasok.
type Material struct {
	ID           uint      `gorm:"primaryKey;column:id" json:"id"`
	KodeSKU      string    `gorm:"column:kode_sku;type:varchar(100);uniqueIndex;not null" json:"kodeSKU"`
	NamaMaterial string    `gorm:"column:nama_material;type:varchar(255);not null" json:"namaMaterial"`
	TebalCM      float64   `gorm:"column:tebal_cm" json:"tebalCM"`
	BeratGSM     float64   `gorm:"column:berat_gsm" json:"beratGSM"`
	LebarCM      float64   `gorm:"column:lebar_cm" json:"lebarCM"`
	PanjangCM    float64   `gorm:"column:panjang_cm" json:"panjangCM"`
	UnitSatuan   string    `gorm:"column:unit_satuan;type:varchar(50)" json:"unitSatuan"`
	IDPemasok    uint      `gorm:"column:id_pemasok;not null" json:"idPemasok"`
	Pemasok      Pemasok   `gorm:"foreignKey:IDPemasok;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"pemasok,omitempty"`
	DibuatPada   time.Time `gorm:"autoCreateTime" json:"dibuatPada"`
	DiubahPada   time.Time `gorm:"autoUpdateTime" json:"diubahPada"`
}

// KomposisiMaterialBOM merepresentasikan relasi kuantitatif dalam Bill of Materials.
type KomposisiMaterialBOM struct {
	ID                          uint      `gorm:"primaryKey;column:id" json:"id"`
	IDProdukFinal               uint      `gorm:"column:id_produk_final;not null" json:"idProdukFinal"`
	IDRawMaterial               uint      `gorm:"column:id_raw_material;not null" json:"idRawMaterial"`
	ParameterKuantitasPembentuk float64   `gorm:"column:parameter_kuantitas_pembentuk;not null" json:"parameterKuantitasPembentuk"`
	MaterialBaku                Material  `gorm:"foreignKey:IDRawMaterial;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"materialBaku,omitempty"`
	DibuatPada                  time.Time `gorm:"autoCreateTime" json:"dibuatPada"`
}
