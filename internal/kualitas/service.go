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
	AmbilOpsiLembarPeriksa() OpsiLembarPeriksaDto
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
	// Validasi Shift (harus NORMAL atau kosong yang dianggap NORMAL)
	if dto.Shift != "" && dto.Shift != "NORMAL" && dto.Shift != "normal" {
		return errors.New("validasi_gagal: hanya shift NORMAL yang didukung oleh sistem")
	}

	// Validasi Tanggal
	if dto.Tanggal == "" {
		return errors.New("validasi_gagal: tanggal tidak boleh kosong")
	}

	// Validasi Zona Lini
	if dto.ZonaLini == "" {
		return errors.New("validasi_gagal: zona lini tidak boleh kosong")
	}

	// Daftar valid time slots sesuai standar PGN
	validTimeSlots := map[string]bool{
		"08:00-12:00":   true,
		"13:00-15:30":   true,
		"16:00-17:30":   true,
		"18:30-selesai": true,
	}

	// Validasi TPS: Total Produksi == OK + NG dan cek rentang waktu
	if len(dto.Detail) == 0 {
		return errors.New("validasi_gagal: detail inspeksi tidak boleh kosong")
	}

	for _, d := range dto.Detail {
		if d.TotalProduksi < 0 || d.RasioTotalOK < 0 || d.RasioCacat < 0 {
			return errors.New("validasi_gagal: nilai produksi dan cacat tidak boleh negatif")
		}
		if d.TotalProduksi != (d.RasioTotalOK + d.RasioCacat) {
			return errors.New("validasi_tps_gagal: total produksi harus presisi sama dengan jumlah OK dan NG")
		}
		if !validTimeSlots[d.WaktuPergeseran] {
			return errors.New("validasi_gagal: waktu pergeseran (time slot) tidak valid")
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

// AmbilOpsiLembarPeriksa mengembalikan daftar konfigurasi statis untuk pengisian QC.
func (l *layananKualitas) AmbilOpsiLembarPeriksa() OpsiLembarPeriksaDto {
	return OpsiLembarPeriksaDto{
		Shifts: []string{"NORMAL"},
		ZonaLini: []string{
			"Lini 1", "Lini 2", "Lini 3", "Lini 4",
		},
		TimeSlots: []string{
			"08:00-12:00",
			"13:00-15:30",
			"16:00-17:30",
			"18:30-selesai",
		},
	}
}
