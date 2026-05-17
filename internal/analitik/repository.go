// Package analitik menangani agregasi dan kalkulasi data pelaporan 7 QC Tools.
package analitik

import (
	"gorm.io/gorm"
)

// RepositoriAnalitik mendefinisikan antarmuka pengumpulan data analitik.
type RepositoriAnalitik interface {
	KalkulasiParetoBulanan(bulan, tahun int) ([]DTOParetoMetrik, error)
}

type repositoriAnalitik struct {
	db *gorm.DB
}

// KonstruksiRepositoriBaru membuat objek RepositoriAnalitik baru.
func KonstruksiRepositoriBaru(db *gorm.DB) RepositoriAnalitik {
	return &repositoriAnalitik{db: db}
}

// KalkulasiParetoBulanan menggunakan Window Function PostgreSQL (SUM OVER)
// untuk menghindari rekursi pada lapisan Go.
func (r *repositoriAnalitik) KalkulasiParetoBulanan(bulan, tahun int) ([]DTOParetoMetrik, error) {
	var hasil []DTOParetoMetrik

	// Kueri Window Function murni PostgreSQL 17.x
	kueri := `
	WITH AgregatDasar AS (
		SELECT 
			kode_cacat, 
			SUM(kuantitas) AS jumlah_cacat
		FROM buku_besar_cacats
		WHERE EXTRACT(MONTH FROM tanggal) = ? AND EXTRACT(YEAR FROM tanggal) = ?
		GROUP BY kode_cacat
	),
	KalkulasiKumulatif AS (
		SELECT 
			kode_cacat,
			jumlah_cacat,
			SUM(jumlah_cacat) OVER (ORDER BY jumlah_cacat DESC) AS kumulatif,
			SUM(jumlah_cacat) OVER () AS total_semua
		FROM AgregatDasar
	)
	SELECT 
		kode_cacat,
		jumlah_cacat,
		CASE WHEN total_semua > 0 THEN (jumlah_cacat::float / total_semua::float) * 100 ELSE 0 END AS persentase,
		CASE WHEN total_semua > 0 THEN (kumulatif::float / total_semua::float) * 100 ELSE 0 END AS persentase_kumulatif
	FROM KalkulasiKumulatif
	ORDER BY jumlah_cacat DESC;
	`

	if err := r.db.Raw(kueri, bulan, tahun).Scan(&hasil).Error; err != nil {
		return nil, err
	}

	return hasil, nil
}
