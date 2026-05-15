package konfigurasi

import (
	"log/slog"
	"pgn-server/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JalankanSeeder memasukkan data master awal ke database
func JalankanSeeder(db *gorm.DB) {
	slog.Info("Memulai proses seeding data master...")

	// 1. Seed Users
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	user := model.User{
		NIP:      "2211019",
		Nama:     "Leader QC",
		Role:     "Leader_QC",
		Password: string(hash),
	}
	db.Where(model.User{NIP: "2211019"}).FirstOrCreate(&user)

	// 2. Seed Lini Produksi
	lini := []model.LiniProduksi{
		{ID: "L01", NamaLini: "Lini Press"},
		{ID: "L02", NamaLini: "Lini Cutting"},
		{ID: "L03", NamaLini: "Lini Sewing"},
	}
	for _, l := range lini {
		db.Where(model.LiniProduksi{ID: l.ID}).FirstOrCreate(&l)
	}

	// 3. Seed Master Mesin
	mesin := []model.MasterMesin{
		{ID: "M-PR-01", NamaMesin: "AIDA 200T", TipeMesin: "Press"},
		{ID: "M-PR-02", NamaMesin: "KOMATSU 110T", TipeMesin: "Press"},
		{ID: "M-CT-01", NamaMesin: "Band Saw A", TipeMesin: "Cutting"},
		{ID: "M-CT-02", NamaMesin: "Laser Cut X1", TipeMesin: "Cutting"},
		{ID: "M-SW-01", NamaMesin: "JUKI Lockstitch", TipeMesin: "Sewing"},
		{ID: "M-SW-02", NamaMesin: "Brother Overlock", TipeMesin: "Sewing"},
	}
	for _, m := range mesin {
		db.Where(model.MasterMesin{ID: m.ID}).FirstOrCreate(&m)
	}

	// 4. Seed Kategori Defect
	kategori := []model.KategoriDefect{
		{ID: "K01", NamaKategori: "Material"},
		{ID: "K02", NamaKategori: "Sewing"},
		{ID: "K03", NamaKategori: "Press"},
	}
	for _, k := range kategori {
		db.Where(model.KategoriDefect{ID: k.ID}).FirstOrCreate(&k)
	}

	// 5. Seed Master Defect
	defects := []model.MasterDefect{
		{ID: "D01", KategoriID: "K01", NamaNG: "Crack"},
		{ID: "D02", KategoriID: "K01", NamaNG: "Rust"},
		{ID: "D03", KategoriID: "K01", NamaNG: "Scratch"},
		{ID: "D04", KategoriID: "K02", NamaNG: "Puckering"},
		{ID: "D05", KategoriID: "K02", NamaNG: "Broken Stitch"},
		{ID: "D06", KategoriID: "K02", NamaNG: "Oil Stain"},
		{ID: "D07", KategoriID: "K03", NamaNG: "Burry"},
		{ID: "D08", KategoriID: "K03", NamaNG: "Dent"},
	}
	for _, d := range defects {
		db.Where(model.MasterDefect{ID: d.ID}).FirstOrCreate(&d)
	}

	slog.Info("Seeding data master selesai.")
}
