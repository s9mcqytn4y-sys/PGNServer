<?php

declare(strict_types=1);

namespace App\Application\QControl;

use App\Models\QControlLineProduksi;
use App\Models\QControlPemeriksaanDefectSlot;
use App\Models\QControlPemeriksaanHarian;
use App\Support\Errors\KodeKesalahanApi;
use Illuminate\Support\Carbon;

/**
 * Menyusun read model bulanan recording defect dari transaksi daily QControl.
 */
final class MembacaLaporanBulananRecordingDefect
{
    /**
     * @param  array<string, mixed>  $filter
     * @return array<string, mixed>
     */
    public function jalankan(array $filter): array
    {
        $lineProduksi = QControlLineProduksi::query()->find($filter['lineProduksiId']);

        if (! $lineProduksi instanceof QControlLineProduksi || ! $lineProduksi->aktif) {
            throw new PengecualianPemeriksaanHarian(
                pesan: 'Line produksi QControl tidak aktif atau tidak tersedia',
                kodeKesalahan: KodeKesalahanApi::VALIDASI_GAGAL,
                detailKesalahan: [
                    [
                        'field' => 'lineProduksiId',
                        'pesan' => 'Line produksi QControl tidak aktif atau tidak tersedia',
                    ],
                ],
            );
        }

        $tahun = (int) $filter['tahun'];
        $bulan = (int) $filter['bulan'];
        $tanggalAwal = Carbon::createMidnightDate($tahun, $bulan, 1);
        $jumlahHari = $tanggalAwal->daysInMonth;
        $daftarTanggal = range(1, $jumlahHari);
        $templatePerTanggal = $this->buatTemplatePerTanggal($daftarTanggal);

        $daftarPemeriksaanHarian = QControlPemeriksaanHarian::query()
            ->whereYear('tanggal_produksi', $tahun)
            ->whereMonth('tanggal_produksi', $bulan)
            ->where('line_produksi_id', $lineProduksi->id)
            ->with([
                'daftarPemeriksaanPart' => function ($query) use ($filter): void {
                    $query->orderBy('urutan_tampil');

                    if (isset($filter['partId'])) {
                        $query->where('part_id', (string) $filter['partId']);
                    }

                    if (isset($filter['materialId'])) {
                        $query->where('material_id_snapshot', (string) $filter['materialId']);
                    }
                },
                'daftarPemeriksaanPart.daftarDefectSlot' => function ($query) use ($filter): void {
                    $query->orderBy('kode_tampilan_defect_snapshot')
                        ->orderBy('dibuat_pada');

                    if (isset($filter['jenisDefectId'])) {
                        $query->where('jenis_defect_id', (string) $filter['jenisDefectId']);
                    }
                },
            ])
            ->orderBy('tanggal_produksi')
            ->orderBy('dibuat_pada')
            ->get();

        /** @var array<string, array<string, mixed>> $daftarPart */
        $daftarPart = [];
        $totalHarian = $templatePerTanggal;

        foreach ($daftarPemeriksaanHarian as $pemeriksaanHarian) {
            $tanggalKe = (int) $pemeriksaanHarian->tanggal_produksi?->day;
            $kunciTanggal = (string) $tanggalKe;

            foreach ($pemeriksaanHarian->daftarPemeriksaanPart as $pemeriksaanPart) {
                $defectTersaring = $pemeriksaanPart->daftarDefectSlot;

                if (isset($filter['jenisDefectId']) && $defectTersaring->isEmpty()) {
                    continue;
                }

                $kunciPart = (string) $pemeriksaanPart->part_id;

                if (! array_key_exists($kunciPart, $daftarPart)) {
                    $daftarPart[$kunciPart] = [
                        'partId' => $kunciPart,
                        'kodeUnikPart' => (string) ($pemeriksaanPart->kode_unik_part_snapshot ?? $pemeriksaanPart->partTerkait?->kode_unik_part),
                        'namaPart' => (string) ($pemeriksaanPart->nama_part_snapshot ?? $pemeriksaanPart->partTerkait?->nama_part),
                        'nomorPart' => (string) ($pemeriksaanPart->nomor_part_snapshot ?? $pemeriksaanPart->partTerkait?->nomor_part),
                        'namaMaterial' => (string) ($pemeriksaanPart->nama_material_snapshot ?? $pemeriksaanPart->partTerkait?->materialTerkait?->nama_material),
                        'totalCheck' => 0,
                        'totalOk' => 0,
                        'totalDefect' => 0,
                        'rasioDefect' => 0.0,
                        'subtotalPerTanggal' => $templatePerTanggal,
                        'daftarDefectMap' => [],
                        'urutanSort' => count($daftarPart) + 1,
                    ];
                }

                $totalDefectPartHari = (int) $defectTersaring
                    ->sum(fn (QControlPemeriksaanDefectSlot $defectSlot): int => (int) $defectSlot->jumlah_defect);
                $totalCheckPartHari = (int) $pemeriksaanPart->total_check;

                $daftarPart[$kunciPart]['totalCheck'] += $totalCheckPartHari;
                $daftarPart[$kunciPart]['totalDefect'] += $totalDefectPartHari;
                $daftarPart[$kunciPart]['subtotalPerTanggal'][$kunciTanggal] += $totalDefectPartHari;
                $totalHarian[$kunciTanggal] += $totalDefectPartHari;

                foreach ($defectTersaring as $defectSlot) {
                    $kunciDefect = (string) ($defectSlot->relasi_part_defect_id ?? $defectSlot->jenis_defect_id ?? $defectSlot->id);

                    if (! array_key_exists($kunciDefect, $daftarPart[$kunciPart]['daftarDefectMap'])) {
                        $daftarPart[$kunciPart]['daftarDefectMap'][$kunciDefect] = [
                            'kodeTampilanDefect' => (string) $defectSlot->kode_tampilan_defect_snapshot,
                            'namaDefect' => (string) $defectSlot->nama_defect_snapshot,
                            'kategoriDefect' => $defectSlot->kategori_defect_snapshot,
                            'jumlahPerTanggal' => $templatePerTanggal,
                            'totalDefect' => 0,
                        ];
                    }

                    $jumlahDefect = (int) $defectSlot->jumlah_defect;

                    $daftarPart[$kunciPart]['daftarDefectMap'][$kunciDefect]['jumlahPerTanggal'][$kunciTanggal] += $jumlahDefect;
                    $daftarPart[$kunciPart]['daftarDefectMap'][$kunciDefect]['totalDefect'] += $jumlahDefect;
                }
            }
        }

        $daftarPartFinal = collect($daftarPart)
            ->sortBy('urutanSort')
            ->map(function (array $part): array {
                $part['totalOk'] = max(0, $part['totalCheck'] - $part['totalDefect']);
                $part['rasioDefect'] = $part['totalCheck'] === 0
                    ? 0.0
                    : round(($part['totalDefect'] / $part['totalCheck']) * 100, 2);

                $part['daftarDefect'] = collect($part['daftarDefectMap'])
                    ->sortBy('kodeTampilanDefect')
                    ->map(function (array $defect): array {
                        $defect['jumlahPerTanggal'] = (object) $defect['jumlahPerTanggal'];

                        return $defect;
                    })
                    ->values()
                    ->all();

                $part['subtotalPerTanggal'] = (object) $part['subtotalPerTanggal'];
                unset($part['daftarDefectMap'], $part['urutanSort']);

                return $part;
            })
            ->values()
            ->all();

        return [
            'bulan' => $bulan,
            'tahun' => $tahun,
            'line' => [
                'id' => $lineProduksi->id,
                'kodeLine' => $lineProduksi->kode_line,
                'namaLine' => $lineProduksi->nama_line,
            ],
            'daftarTanggal' => $daftarTanggal,
            'daftarPart' => $daftarPartFinal,
            'totalHarian' => (object) $totalHarian,
            'totalBulanan' => array_sum($totalHarian),
        ];
    }

    /**
     * @param  list<int>  $daftarTanggal
     * @return array<string, int>
     */
    private function buatTemplatePerTanggal(array $daftarTanggal): array
    {
        $template = [];

        foreach ($daftarTanggal as $tanggal) {
            $template[(string) $tanggal] = 0;
        }

        return $template;
    }
}
