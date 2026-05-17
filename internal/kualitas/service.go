// Package kualitas menyediakan layanan dan repositori untuk modul kontrol kualitas.
package kualitas

import (
	"errors"
	"log"

	"gorm.io/gorm"
	"pgn-server/internal/manufaktur"
)

// LayananKualitas melayani penelusuran BOM dan pencatatan komposit.
type LayananKualitas interface {
	RekamLembarPeriksa(dto *DTOLembarPeriksaKirim) error
	DaftarRiwayat(limit int, offset int, tanggalMulai string, tanggalSelesai string, zonaLini string) ([]LembarPeriksa, error)
}

type layananKualitas struct {
	repo RepositoriKualitas
	db   *gorm.DB
}

func KonstruksiLayananBaru(repo RepositoriKualitas, db *gorm.DB) LayananKualitas {
	return &layananKualitas{repo: repo, db: db}
}

// RekamLembarPeriksa mengatur arus logika bisnis pencatatan inspeksi fisik.
func (l *layananKualitas) RekamLembarPeriksa(dto *DTOLembarPeriksaKirim) error {
	// Validasi TPS: Total Produksi == OK + NG
	for _, d := range dto.Detail {
		if d.TotalProduksi != (d.RasioTotalOK + d.RasioCacat) {
			return errors.New("validasi_tps_gagal: total produksi harus sama dengan jumlah OK dan NG")
		}
	}

	// Penjagaan skema atomik menggunakan GORM Transaction untuk auto-rollback pada panic/error.
	errTx := l.db.Transaction(func(tx *gorm.DB) error {
		// Langkah 1: Eksekusi penyisipan cepat multi-baris (menghindari N+1)
		errSimpan := l.repo.SimpanLembarPeriksaMassal(dto, tx)
		if errSimpan != nil {
			return errSimpan
		}

		// Langkah 2: Algoritma penelusuran susunan graf hirarki BOM.
		// Jika ditemukan kode cacat spesifik, sistem secara independen mendeklarasikan
		// catatan buku besar penyusutan.
		for _, detail := range dto.Detail {
			if detail.RasioCacat > 0 {
				var bom manufaktur.KomposisiMaterialBOM

				// Kueri melacak ID material komponen berdasarkan hirarki BOM produk akhir (UnikPartID).
				errLacak := tx.Joins("JOIN materials ON materials.id = komposisi_material_boms.id_raw_material").
					Where("komposisi_material_boms.id_produk_final = ?", detail.UnikPartID).
					First(&bom).Error

				if errLacak == nil {
					// Deklarasikan pencatatan rasio penyusutan hanya terhadap entri Buku Besar Cacat
					entriLedger := BukuBesarCacat{
						IDMaterial:      bom.IDRawMaterial,
						TotalPenyusutan: detail.RasioCacat,
					}
					if errLedger := tx.Create(&entriLedger).Error; errLedger != nil {
						return errLedger
					}
				} else if !errors.Is(errLacak, gorm.ErrRecordNotFound) {
					return errLacak
				} else {
					// Bila BOM tidak diatur, log info namun tak batalkan transaksi.
					log.Printf("Informasi: Relasi material BOM tidak terdeteksi untuk part ID %d. Mengabaikan buku besar.", detail.UnikPartID)
				}
			}
		}
		return nil
	})

	return errTx
}

// DaftarRiwayat mengembalikan daftar riwayat lembar periksa harian.
func (l *layananKualitas) DaftarRiwayat(limit int, offset int, tanggalMulai string, tanggalSelesai string, zonaLini string) ([]LembarPeriksa, error) {
	return l.repo.DaftarRiwayat(limit, offset, tanggalMulai, tanggalSelesai, zonaLini)
}
