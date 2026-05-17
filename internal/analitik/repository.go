// Package analitik menangani agregasi dan kalkulasi data pelaporan 7 QC Tools.
package analitik

import (
	"gorm.io/gorm"
)

// RepositoriAnalitik mendefinisikan antarmuka pengumpulan data analitik.
type RepositoriAnalitik interface {
	KalkulasiParetoBulanan(bulan, tahun int) ([]DTOParetoMetrik, error)
	KalkulasiPareto(tanggalMulai, tanggalSelesai string, zonaLini string) ([]DTOParetoMetrik, error)
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
			d.kode_cacat, 
			SUM(d.rasio_cacat) AS jumlah_cacat
		FROM detail_inspeksis d
		JOIN lembar_periksas l ON d.lembar_periksa_id = l.id
		WHERE EXTRACT(MONTH FROM l.tanggal) = ? AND EXTRACT(YEAR FROM l.tanggal) = ?
		GROUP BY d.kode_cacat
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

// KalkulasiPareto menghitung Pareto dengan filter tanggalMulai, tanggalSelesai, dan zonaLini
func (r *repositoriAnalitik) KalkulasiPareto(tanggalMulai, tanggalSelesai string, zonaLini string) ([]DTOParetoMetrik, error) {
	var hasil []DTOParetoMetrik

	// Build raw query with dynamic parameters
	kueriAgregat := `
		SELECT 
			d.kode_cacat, 
			SUM(d.rasio_cacat) AS jumlah_cacat
		FROM detail_inspeksis d
		JOIN lembar_periksas l ON d.lembar_periksa_id = l.id
		WHERE 1=1
	`
	var args []interface{}

	if tanggalMulai != "" {
		kueriAgregat += " AND l.tanggal >= ?"
		args = append(args, tanggalMulai)
	}
	if tanggalSelesai != "" {
		kueriAgregat += " AND l.tanggal <= ?"
		args = append(args, tanggalSelesai)
	}
	if zonaLini != "" {
		kueriAgregat += " AND l.zona_lini = ?"
		args = append(args, zonaLini)
	}

	kueriAgregat += " GROUP BY d.kode_cacat"

	kueriUtama := `
	WITH AgregatDasar AS (` + kueriAgregat + `),
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

	if err := r.db.Raw(kueriUtama, args...).Scan(&hasil).Error; err != nil {
		return nil, err
	}

	return hasil, nil
}
