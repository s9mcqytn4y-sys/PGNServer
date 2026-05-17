package manufaktur

import (
	"log"

	"gorm.io/gorm"
)

// JalankanSeeder akan mengisi database dengan data awal manufaktur
func JalankanSeeder(db *gorm.DB) error {
	// Seeder Customer
	customers := []Customer{
		{CustomerCode: "CUST-001", Nama: "PT. Astra Honda Motor", Kontak: "pic.ahm@astra.co.id"},
		{CustomerCode: "CUST-002", Nama: "PT. Toyota Motor Manufacturing", Kontak: "pic.tmin@toyota.co.id"},
	}

	for _, c := range customers {
		var count int64
		db.Model(&Customer{}).Where("customer_code = ?", c.CustomerCode).Count(&count)
		if count == 0 {
			db.Create(&c)
		}
	}

	// Seeder Pemasok (Pemasok)
	pemasok := []Pemasok{
		{SupplierCode: "SUP-001", NamaEntitas: "PT. Krakatau Steel", Kontak: "sales@krakatausteel.com"},
		{SupplierCode: "SUP-002", NamaEntitas: "PT. YKK Zip", Kontak: "sales@ykk.co.id"},
		{SupplierCode: "SUP-BTI", NamaEntitas: "PT Artha Langgeng Mulya (BTI)", Kontak: "sales@arthalanggeng.co.id"},
	}

	for _, p := range pemasok {
		var count int64
		db.Model(&Pemasok{}).Where("supplier_code = ?", p.SupplierCode).Count(&count)
		if count == 0 {
			db.Create(&p)
		}
	}

	// Fetch Pemasok
	var sup1, supBTI Pemasok
	db.Where("supplier_code = ?", "SUP-001").First(&sup1)
	db.Where("supplier_code = ?", "SUP-BTI").First(&supBTI)

	// Seeder Material
	materials := []Material{
		{KodeSKU: "MAT-A01", NamaMaterial: "Plat Baja 2mm", TebalCM: 0.2, BeratGSM: 500, UnitSatuan: "Pcs", IDPemasok: sup1.ID},
		{KodeSKU: "MAT-B01", NamaMaterial: "Baut M8", UnitSatuan: "Pcs", IDPemasok: sup1.ID},
		{KodeSKU: "FG-001", NamaMaterial: "Bracket Engine Mount", UnitSatuan: "Pcs", IDPemasok: sup1.ID}, // Finished Good 1
		
		// TPS/BTI Materials
		{KodeSKU: "MAT-001", NamaMaterial: "PS Polyester Non Woven Spunbond 100 Gsm White", TebalCM: 0.1, BeratGSM: 100, UnitSatuan: "Roll", IDPemasok: supBTI.ID},
		{KodeSKU: "MAT-002", NamaMaterial: "Laminasi LDPE 200 Gsm", TebalCM: 0.05, BeratGSM: 200, UnitSatuan: "Roll", IDPemasok: supBTI.ID},
		{KodeSKU: "MAT-003", NamaMaterial: "Recycle Felt GWPS 2mm 375 Gsm", TebalCM: 0.2, BeratGSM: 375, UnitSatuan: "Roll", IDPemasok: supBTI.ID},
		{KodeSKU: "MAT-004", NamaMaterial: "Carpet STKD19 Black", TebalCM: 0.3, BeratGSM: 400, UnitSatuan: "Roll", IDPemasok: supBTI.ID},
		{KodeSKU: "MAT-005", NamaMaterial: "Carpet CBIII", TebalCM: 0.4, BeratGSM: 450, UnitSatuan: "Roll", IDPemasok: supBTI.ID},
		{KodeSKU: "MAT-006", NamaMaterial: "EPDM", TebalCM: 0.15, BeratGSM: 300, UnitSatuan: "Roll", IDPemasok: supBTI.ID},
		{KodeSKU: "MAT-007", NamaMaterial: "Ester Canvas SAB10-NS121 SSP", TebalCM: 0.25, BeratGSM: 350, UnitSatuan: "Roll", IDPemasok: supBTI.ID},
		{KodeSKU: "MAT-008", NamaMaterial: "Silincer T. 15mm 1000 Gsm", TebalCM: 1.5, BeratGSM: 1000, UnitSatuan: "Roll", IDPemasok: supBTI.ID},
		{KodeSKU: "MAT-009", NamaMaterial: "Silincer T. 6mm 350 Gsm", TebalCM: 0.6, BeratGSM: 350, UnitSatuan: "Roll", IDPemasok: supBTI.ID},

		// Finished Good: Protector
		{KodeSKU: "FG-002", NamaMaterial: "Protector", TebalCM: 0.0, BeratGSM: 0, UnitSatuan: "Pcs", IDPemasok: supBTI.ID},
	}

	for _, m := range materials {
		var count int64
		db.Model(&Material{}).Where("kode_sku = ?", m.KodeSKU).Count(&count)
		if count == 0 {
			db.Create(&m)
		} else {
			// Update Pemasok ID just in case
			db.Model(&Material{}).Where("kode_sku = ?", m.KodeSKU).Update("id_pemasok", m.IDPemasok)
		}
	}

	// Fetch Material
	var matBaja, matBaut, fgBracket Material
	db.Where("kode_sku = ?", "MAT-A01").First(&matBaja)
	db.Where("kode_sku = ?", "MAT-B01").First(&matBaut)
	db.Where("kode_sku = ?", "FG-001").First(&fgBracket)

	var fgProtector Material
	db.Where("kode_sku = ?", "FG-002").First(&fgProtector)

	var btiMats []Material
	db.Where("kode_sku IN ?", []string{
		"MAT-001", "MAT-002", "MAT-003", "MAT-004", "MAT-005", "MAT-006", "MAT-007", "MAT-008", "MAT-009",
	}).Find(&btiMats)

	// Seeder BOM (KomposisiMaterialBOM)
	boms := []KomposisiMaterialBOM{
		{IDParentMaterial: &fgBracket.ID, IDRawMaterial: matBaja.ID, ParameterKuantitasPembentuk: 1.5},
		{IDParentMaterial: &fgBracket.ID, IDRawMaterial: matBaut.ID, ParameterKuantitasPembentuk: 4.0},
	}

	// Tambahkan komposisi material pembentuk Protector
	for _, raw := range btiMats {
		boms = append(boms, KomposisiMaterialBOM{
			IDParentMaterial:            &fgProtector.ID,
			IDRawMaterial:               raw.ID,
			ParameterKuantitasPembentuk: 1.0, // Default 1.0 roll per Pcs protector
		})
	}

	for _, bom := range boms {
		var count int64
		db.Model(&KomposisiMaterialBOM{}).Where("id_parent_material = ? AND id_raw_material = ?", bom.IDParentMaterial, bom.IDRawMaterial).Count(&count)
		if count == 0 {
			db.Create(&bom)
		}
	}

	log.Println("Seeder Manufaktur berhasil dieksekusi")
	return nil
}
