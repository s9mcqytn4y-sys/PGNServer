package media

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pgn-server/internal/manufaktur"

	"gorm.io/gorm"
)

// Konstanta sistem untuk kebijakan media
const (
	MaksimalUkuranBerkas = 5 * 1024 * 1024 // 5 MB
	DirektoriSimpanUtama = "./penyimpanan/materials/"
)

// LayananMedia mengatur logika bisnis untuk verifikasi, validasi relasi, dan persisten file.
type LayananMedia interface {
	UnggahBerkasLokal(idMaterial uint, berkas *multipart.FileHeader) (*AsetDigital, error)
	HapusBerkas(id uint) error
	CariAsetMedia(id uint) (*AsetDigital, error)
}

type layananMedia struct {
	repoMedia RepositoriMedia
	db        *gorm.DB // Untuk pengecekan part reference validation
}

// KonstruksiLayananBaru membuat instance LayananMedia.
func KonstruksiLayananBaru(repo RepositoriMedia, db *gorm.DB) LayananMedia {
	// Pastikan direktori unggahan eksis
	_ = os.MkdirAll(DirektoriSimpanUtama, os.ModePerm)
	_ = os.MkdirAll("./penyimpanan/profiles/", os.ModePerm)

	// Inisialisasi dummy PNG untuk fallback
	dummyPNG := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89\x00\x00\x00\rIDATx\x9cc`\x00\x01\x00\x00\x05\x00\x01\r\n-\xb4\x00\x00\x00\x00IEND\xaeB`\x82")
	_ = os.WriteFile("./penyimpanan/materials/part.png", dummyPNG, 0644)
	_ = os.WriteFile("./penyimpanan/profiles/avatar.png", dummyPNG, 0644)

	return &layananMedia{repoMedia: repo, db: db}
}

func (l *layananMedia) validasiTipeMIME(tipe string) bool {
	validMIME := []string{"image/jpeg", "image/png"}
	for _, v := range validMIME {
		if tipe == v {
			return true
		}
	}
	return false
}

func (l *layananMedia) validasiEkstensi(namaBerkas string) bool {
	ekstensi := strings.ToLower(filepath.Ext(namaBerkas))
	validEkstensi := []string{".jpg", ".jpeg", ".png"}
	for _, v := range validEkstensi {
		if ekstensi == v {
			return true
		}
	}
	return false
}

func (l *layananMedia) UnggahBerkasLokal(idMaterial uint, berkas *multipart.FileHeader) (*AsetDigital, error) {
	// 1. Part Reference Validation: Validasi eksistensi IDMaterial
	var material manufaktur.Material
	if errCari := l.db.First(&material, idMaterial).Error; errCari != nil {
		return nil, errors.New("referensi_material_tidak_ditemukan")
	}

	// 2. File Validation: Maksimal 5MB
	if berkas.Size > MaksimalUkuranBerkas {
		return nil, errors.New("ukuran_berkas_terlalu_besar")
	}

	// 3. File Validation: Ekstensi
	if !l.validasiEkstensi(berkas.Filename) {
		return nil, errors.New("ekstensi_berkas_tidak_diizinkan")
	}

	berkasDibuka, errBuka := berkas.Open()
	if errBuka != nil {
		return nil, errBuka
	}
	defer berkasDibuka.Close()

	// 4. File Validation: MIME Type
	buffer := make([]byte, 512)
	_, errBaca := berkasDibuka.Read(buffer)
	if errBaca != nil && errBaca != io.EOF {
		return nil, errBaca
	}

	tipeMimeDeteksi := http.DetectContentType(buffer)
	if !l.validasiTipeMIME(tipeMimeDeteksi) {
		return nil, errors.New("tipe_mime_berkas_tidak_valid")
	}

	// Untuk keamanan nyata, kita akan deteksi ulang
	// Reset pointer file
	if _, errSeek := berkasDibuka.Seek(0, io.SeekStart); errSeek != nil {
		return nil, errSeek
	}

	tipeMIMEAsli := berkas.Header.Get("Content-Type")
	if !l.validasiTipeMIME(tipeMIMEAsli) {
		return nil, errors.New("tipe_mime_berkas_tidak_valid")
	}

	// Persiapkan penyimpanan
	namaBerkasBaru := fmt.Sprintf("mat_%d_%d%s", idMaterial, time.Now().UnixNano(), filepath.Ext(berkas.Filename))
	jalurTujuan := filepath.Join(DirektoriSimpanUtama, namaBerkasBaru)

	fileTujuan, errBuat := os.Create(jalurTujuan)
	if errBuat != nil {
		return nil, errBuat
	}
	defer fileTujuan.Close()

	if _, errSalin := io.Copy(fileTujuan, berkasDibuka); errSalin != nil {
		return nil, errSalin
	}

	// 5. Rekam ke Database
	asetBaru := &AsetDigital{
		IDMaterial:      idMaterial,
		TipeMIME:        tipeMIMEAsli,
		UkuranBerkas:    berkas.Size,
		Ekstensi:        strings.ToLower(filepath.Ext(berkas.Filename)),
		DirektoriLokal:  jalurTujuan,
		TipePenyimpanan: "LOKAL",
	}

	if errSimpan := l.repoMedia.Simpan(asetBaru); errSimpan != nil {
		_ = os.Remove(jalurTujuan) // rollback file fisik jika DB gagal
		return nil, errSimpan
	}

	return asetBaru, nil
}

func (l *layananMedia) HapusBerkas(id uint) error {
	aset, errCari := l.repoMedia.CariBerdasarkanID(id)
	if errCari != nil {
		return errCari
	}

	if aset.TipePenyimpanan == "LOKAL" {
		// Proteksi file default agar tidak terhapus
		if !strings.Contains(aset.DirektoriLokal, "avatar.png") && !strings.Contains(aset.DirektoriLokal, "part.png") {
			_ = os.Remove(aset.DirektoriLokal)
		}
	}

	return l.repoMedia.HapusBerdasarkanID(id)
}

func (l *layananMedia) CariAsetMedia(id uint) (*AsetDigital, error) {
	aset, err := l.repoMedia.CariBerdasarkanID(id)
	if err != nil {
		// Fallback default jika tidak ada di DB
		return &AsetDigital{
			TipePenyimpanan: "LOKAL",
			DirektoriLokal:  "./penyimpanan/materials/part.png",
			TipeMIME:        "image/png",
		}, nil
	}

	if aset.TipePenyimpanan == "LOKAL" {
		if _, errStat := os.Stat(aset.DirektoriLokal); os.IsNotExist(errStat) {
			if aset.IDMaterial != 0 {
				aset.DirektoriLokal = "./penyimpanan/materials/part.png"
			} else {
				aset.DirektoriLokal = "./penyimpanan/profiles/avatar.png"
			}
		}
	}

	return aset, nil
}
