package konfigurasi

import (
	"log/slog"
	"pgn-server/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// JalankanSeeder memasukkan data master awal ke database
func JalankanSeeder(db *gorm.DB) {
	slog.Info("Memulai proses seeding data master...")

	// 1. User
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	user := model.User{
		NIP:      "2211019",
		Nama:     "Leader QC System",
		Role:     "Leader QC",
		Password: string(hash),
	}
	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "nip"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "role", "password"}),
	}).Create(&user)

	// 2. Customers
	customers := []model.Customer{
		{ID: "CUST-001", NamaCustomer: "BONECOM TRICOM"},
		{ID: "CUST-002", NamaCustomer: "RAJAWALI MITRA PRATAMA"},
		{ID: "CUST-003", NamaCustomer: "RAVALIA INTI MANDIRI"},
	}
	for _, c := range customers {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&c)
	}

	// 3. Suppliers
	suppliers := []model.Supplier{
		{ID: "SUP-001", NamaSupplier: "PT. ARTHA LANGGENG MULYA"},
		{ID: "SUP-002", NamaSupplier: "PT. BONECOM"},
		{ID: "SUP-003", NamaSupplier: "PT. HASIL DAMAI TEXTILE"},
	}
	for _, s := range suppliers {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&s)
	}

	// 4. Production Lines
	lines := []model.LiniProduksi{
		{ID: "PRESS", NamaLini: "PRESS"},
		{ID: "SEWING", NamaLini: "SEWING"},
		{ID: "PASSTROUGH", NamaLini: "PASSTROUGH"},
		{ID: "MATERIAL", NamaLini: "MATERIAL"},
		{ID: "CUTTING", NamaLini: "CUTTING"},
	}
	for _, l := range lines {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&l)
	}

	// 5. Materials
	materials := []model.Material{
		{ID: "BC2", NamaPart: "Laminasi LDPE 200 Gsm", Satuan: "Roll", SupplierID: "SUP-001"},
		{ID: "9Y8_MAT", NamaPart: "HARDFELT 375", Satuan: "Roll", SupplierID: "SUP-001"},
		{ID: "CB9_MAT", NamaPart: "Carpet CB-III", Satuan: "Roll", SupplierID: "SUP-001"},
		{ID: "EPDM1", NamaPart: "EPDM Rubber", Satuan: "Pcs", SupplierID: "SUP-001"},
	}
	for _, m := range materials {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&m)
	}

	// 6. Products
	assy1 := "ASSY-CONSOLE"
	products := []model.Produk{
		{ID: "PROD-001", NomorPart: "58815-KK010-00", NomorUnik: "CB9", NamaPart: "CARPET CONSOLE BOX", CustomerID: "CUST-001", LineID: "PRESS", AssyName: &assy1},
		{ID: "PROD-002", NomorPart: "71695-VT070", NomorUnik: "B35", NamaPart: "PROTECTOR RR SEAT BACK", CustomerID: "CUST-001", LineID: "SEWING", AssyName: nil},
	}
	for _, p := range products {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&p)
	}

	// 7. Defect Master
	defects := []model.DefectMaster{
		{ID: "DEF-MAT-001", Kategori: "MATERIAL", NamaNG: "CARPET TIPIS"},
		{ID: "DEF-MAT-002", Kategori: "MATERIAL", NamaNG: "BELANG"},
		{ID: "DEF-MAT-003", Kategori: "MATERIAL", NamaNG: "BRUDUL"},
		{ID: "DEF-PROC-001", Kategori: "PROCESS", NamaNG: "DENT"},
		{ID: "DEF-PROC-002", Kategori: "PROCESS", NamaNG: "GALER"},
		{ID: "DEF-PROC-003", Kategori: "PROCESS", NamaNG: "OVERCUTTING"},
	}
	for _, d := range defects {
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&d)
	}

	slog.Info("Seeding data master selesai.")
}
