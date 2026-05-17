// Package analitik menangani agregasi dan kalkulasi data pelaporan 7 QC Tools.
package analitik

// LayananAnalitik menyediakan logika pelaporan.
type LayananAnalitik interface {
	DapatkanParetoBulanan(bulan, tahun int) ([]DTOParetoMetrik, error)
}

type layananAnalitik struct {
	repo RepositoriAnalitik
}

// KonstruksiLayananBaru membuat objek LayananAnalitik.
func KonstruksiLayananBaru(repo RepositoriAnalitik) LayananAnalitik {
	return &layananAnalitik{repo: repo}
}

// DapatkanParetoBulanan mengorkestrasi penarikan agregasi Pareto.
func (l *layananAnalitik) DapatkanParetoBulanan(bulan, tahun int) ([]DTOParetoMetrik, error) {
	return l.repo.KalkulasiParetoBulanan(bulan, tahun)
}
