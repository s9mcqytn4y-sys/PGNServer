<?php

declare(strict_types=1);

namespace App\Application\QControl;

use App\Models\QControlJenisDefect;
use App\Models\QControlLineProduksi;
use App\Models\QControlPart;
use App\Models\QControlPartJenisDefect;
use App\Models\QControlSlotWaktu;
use App\Models\User;

/**
 * Menjalankan audit master data QControl agar sinkron dengan template QC sumber.
 */
final class MemvalidasiMasterDataQControl
{
    /**
     * @return array{
     *     valid: bool,
     *     ringkasan: array<string, int>,
     *     temuan: list<string>
     * }
     */
    public function jalankan(): array
    {
        $temuan = [];

        $this->validasiRoleHeadQc($temuan);
        $this->validasiLineProduksi($temuan);
        $this->validasiSlotWaktu($temuan);
        $this->validasiPart($temuan);
        $this->validasiJenisDefect($temuan);
        $this->validasiRelasiPartDefect($temuan);
        $this->validasiTemplateExcel($temuan);

        return [
            'valid' => $temuan === [],
            'ringkasan' => [
                'jumlahLineProduksiAktif' => QControlLineProduksi::query()->where('aktif', true)->count(),
                'jumlahSlotWaktuAktif' => QControlSlotWaktu::query()->where('aktif', true)->count(),
                'jumlahPartAktif' => QControlPart::query()->where('aktif', true)->count(),
                'jumlahJenisDefectAktif' => QControlJenisDefect::query()->where('aktif', true)->count(),
                'jumlahRelasiPartDefectAktif' => QControlPartJenisDefect::query()->where('aktif', true)->count(),
            ],
            'temuan' => $temuan,
        ];
    }

    /**
     * @param  list<string>  $temuan
     */
    private function validasiRoleHeadQc(array &$temuan): void
    {
        $peranHeadQc = (string) config('qcontrol.headqc.peran');

        if (! User::query()->where('peran', $peranHeadQc)->exists()) {
            $temuan[] = 'Role HeadQC belum tersedia pada data pengguna.';
        }
    }

    /**
     * @param  list<string>  $temuan
     */
    private function validasiLineProduksi(array &$temuan): void
    {
        $lineAktif = QControlLineProduksi::query()
            ->where('aktif', true)
            ->orderBy('urutan_tampil')
            ->pluck('kode_line')
            ->all();

        if ($lineAktif !== ['PRESS', 'SEWING']) {
            $temuan[] = 'Line produksi aktif wajib tepat PRESS dan SEWING dengan urutan yang benar.';
        }
    }

    /**
     * @param  list<string>  $temuan
     */
    private function validasiSlotWaktu(array &$temuan): void
    {
        $slotAktif = QControlSlotWaktu::query()
            ->where('aktif', true)
            ->orderBy('urutan_tampil')
            ->get(['kode_slot', 'label_slot']);

        $slotDiharapkan = [
            ['kode_slot' => 'SLOT_0800_1200', 'label_slot' => '08.00 - 12.00'],
            ['kode_slot' => 'SLOT_1300_1530', 'label_slot' => '13.00 - 15.30'],
            ['kode_slot' => 'SLOT_1600_1730', 'label_slot' => '16.00 - 17.30'],
            ['kode_slot' => 'SLOT_1830_SELESAI', 'label_slot' => '18.30 - Selesai'],
        ];

        if ($slotAktif->map(fn (QControlSlotWaktu $slot): array => [
            'kode_slot' => $slot->kode_slot,
            'label_slot' => $slot->label_slot,
        ])->values()->all() !== $slotDiharapkan) {
            $temuan[] = 'Slot waktu aktif harus mengikuti urutan 08.00-12.00, 13.00-15.30, 16.00-17.30, lalu 18.30-Selesai.';
        }
    }

    /**
     * @param  list<string>  $temuan
     */
    private function validasiPart(array &$temuan): void
    {
        $partTidakLengkap = QControlPart::query()
            ->where(function ($query): void {
                $query->whereNull('kode_unik_part')
                    ->orWhere('kode_unik_part', '')
                    ->orWhereNull('nama_part')
                    ->orWhere('nama_part', '')
                    ->orWhereNull('nomor_part')
                    ->orWhere('nomor_part', '')
                    ->orWhereNull('line_default_id')
                    ->orWhere('aktif', false);
            })
            ->pluck('kode_unik_part')
            ->all();

        if ($partTidakLengkap !== []) {
            $temuan[] = 'Masih ada part yang belum lengkap atau tidak aktif: '.implode(', ', $partTidakLengkap);
        }

        $partTanpaTemplate = QControlPart::query()
            ->where('aktif', true)
            ->whereDoesntHave('daftarRelasiJenisDefect', fn ($query) => $query->where('aktif', true))
            ->pluck('kode_unik_part')
            ->all();

        if ($partTanpaTemplate !== []) {
            $temuan[] = 'Masih ada part aktif tanpa template defect: '.implode(', ', $partTanpaTemplate);
        }

        $duplikasiKodePart = QControlPart::query()
            ->select('kode_unik_part')
            ->groupBy('kode_unik_part')
            ->havingRaw('COUNT(*) > 1')
            ->pluck('kode_unik_part')
            ->all();

        if ($duplikasiKodePart !== []) {
            $temuan[] = 'Ditemukan duplikasi kode unik part: '.implode(', ', $duplikasiKodePart);
        }
    }

    /**
     * @param  list<string>  $temuan
     */
    private function validasiJenisDefect(array &$temuan): void
    {
        $defectTidakLengkap = QControlJenisDefect::query()
            ->where(function ($query): void {
                $query->whereNull('kode_defect')
                    ->orWhere('kode_defect', '')
                    ->orWhereNull('nama_defect')
                    ->orWhere('nama_defect', '')
                    ->orWhereNull('kategori_defect_id')
                    ->orWhere('aktif', false);
            })
            ->pluck('kode_defect')
            ->all();

        if ($defectTidakLengkap !== []) {
            $temuan[] = 'Masih ada jenis defect yang belum lengkap atau tidak aktif: '.implode(', ', $defectTidakLengkap);
        }

        $duplikasiKodeDefect = QControlJenisDefect::query()
            ->select('kode_defect')
            ->groupBy('kode_defect')
            ->havingRaw('COUNT(*) > 1')
            ->pluck('kode_defect')
            ->all();

        if ($duplikasiKodeDefect !== []) {
            $temuan[] = 'Ditemukan duplikasi kode defect: '.implode(', ', $duplikasiKodeDefect);
        }
    }

    /**
     * @param  list<string>  $temuan
     */
    private function validasiRelasiPartDefect(array &$temuan): void
    {
        $relasiTidakLengkap = QControlPartJenisDefect::query()
            ->where(function ($query): void {
                $query->whereNull('part_id')
                    ->orWhereNull('jenis_defect_id')
                    ->orWhereNull('kode_tampilan_defect')
                    ->orWhere('kode_tampilan_defect', '')
                    ->orWhereNull('urutan_tampil')
                    ->orWhere('aktif', false);
            })
            ->count();

        if ($relasiTidakLengkap > 0) {
            $temuan[] = 'Masih ada relasi part-defect yang belum lengkap atau tidak aktif.';
        }

        $duplikasiRelasi = QControlPartJenisDefect::query()
            ->select('part_id', 'jenis_defect_id')
            ->groupBy('part_id', 'jenis_defect_id')
            ->havingRaw('COUNT(*) > 1')
            ->count();

        if ($duplikasiRelasi > 0) {
            $temuan[] = 'Ditemukan relasi part-defect ganda pada master data.';
        }

        $duplikasiKodeTampilan = QControlPartJenisDefect::query()
            ->select('part_id', 'kode_tampilan_defect')
            ->groupBy('part_id', 'kode_tampilan_defect')
            ->havingRaw('COUNT(*) > 1')
            ->count();

        if ($duplikasiKodeTampilan > 0) {
            $temuan[] = 'Ditemukan kode tampilan defect duplikat pada part yang sama.';
        }
    }

    /**
     * @param  list<string>  $temuan
     */
    private function validasiTemplateExcel(array &$temuan): void
    {
        foreach ($this->templateBerdasarkanPart() as $kodePart => $template) {
            $part = QControlPart::query()
                ->with(['lineProduksiDefault', 'daftarRelasiJenisDefect' => fn ($query) => $query->with('jenisDefectTerkait')->where('aktif', true)])
                ->where('kode_unik_part', $kodePart)
                ->first();

            if (! $part instanceof QControlPart) {
                $temuan[] = 'Part template '.$kodePart.' tidak ditemukan pada master data.';

                continue;
            }

            if (! $part->aktif) {
                $temuan[] = 'Part template '.$kodePart.' harus aktif.';
            }

            if ($part->lineProduksiDefault?->kode_line !== $template['kodeLine']) {
                $temuan[] = 'Part '.$kodePart.' harus berada pada line '.$template['kodeLine'].'.';
            }

            $templateAktif = $part->daftarRelasiJenisDefect
                ->sortBy('urutan_tampil')
                ->mapWithKeys(fn (QControlPartJenisDefect $relasi): array => [
                    (string) $relasi->kode_tampilan_defect => (string) $relasi->jenisDefectTerkait?->kode_defect,
                ])
                ->all();

            if ($templateAktif !== $template['daftarDefect']) {
                $temuan[] = 'Template defect untuk part '.$kodePart.' tidak cocok dengan struktur Excel.';
            }
        }
    }

    /**
     * @return array<string, array{kodeLine: string, daftarDefect: array<string, string>}>
     */
    private function templateBerdasarkanPart(): array
    {
        $templatePressCarpet = [
            'A' => 'PENYOK',
            'B' => 'GALER',
            'C' => 'KARPET_TIPIS',
            'D' => 'BELANG',
            'E' => 'HOLE_TA',
            'F' => 'POTONGAN_BERLEBIH',
            'G' => 'SOBEK',
            'H' => 'TERLIPAT',
        ];

        $templatePressUmum = [
            'A' => 'LAMINASI_BERKERUT',
            'B' => 'LAMINASI_BOLONG',
            'C' => 'LAMINASI_TIDAK_MATANG',
            'D' => 'TERDAPAT_BENDA_ASING',
            'E' => 'BAHAN_TIPIS',
            'F' => 'POTONGAN_BERLEBIH',
            'G' => 'LAMINASI_TERSOBEK',
            'H' => 'DIMENSI_TIDAK_SESUAI',
        ];

        $templateCb9 = [
            'A' => 'TERDAPAT_BENDA_ASING',
            'B' => 'PENYOK',
            'C' => 'KARPET_BERJAMUR',
            'D' => 'KARPET_TIPIS',
            'E' => 'POTONGAN_BERLEBIH',
            'F' => 'SOBEK',
            'G' => 'DIMENSI_TIDAK_SESUAI',
        ];

        $templateFeltSewing = [
            'A' => 'SOBEK',
            'B' => 'BRUDUL',
            'C' => 'SPUNBOND_TIDAK_MEREKAT',
            'D' => 'SPUNBOND_TERLIPAT',
            'E' => 'LAMINATING_TIDAK_MATANG',
            'F' => 'LAMINATING_BOLONG',
            'G' => 'SPUNBOND_TERPOTONG',
            'H' => 'TERBALIK',
            'I' => 'OVERCUTTING',
            'J' => 'SEWING_MIRING',
            'K' => 'MARGIN_OUT_DIMENSI',
            'L' => 'BACKSTITCH_KURANG_DARI_15MM',
        ];

        $templateCoverSewing = [
            'A' => 'SOBEK',
            'B' => 'BRUDUL',
            'C' => 'TERDAPAT_BENDA_ASING',
            'D' => 'KARPET_BERJAMUR',
            'E' => 'TERBALIK',
            'F' => 'POTONGAN_BERLEBIH',
            'G' => 'SEWING_MIRING',
            'H' => 'MAGIC_TAPE_TERBALIK',
            'I' => 'MARGIN_OUT_DIMENSI',
            'J' => 'SEWING_LONCAT',
            'K' => 'MAGIC_TAPE_MIRING',
            'L' => 'MAGIC_TAPE_TIDAK_TERSEWING',
        ];

        $templateProtector = [
            'A' => 'SOBEK',
            'B' => 'BRUDUL',
            'C' => 'SPUNBOND_TIDAK_MEREKAT',
            'D' => 'SPUNBOND_TERLIPAT',
            'E' => 'SPUNBOND_HARDEN',
            'F' => 'TERDAPAT_BENDA_ASING',
            'G' => 'LAMINATING_TIDAK_MATANG',
            'H' => 'LAMINATING_BOLONG',
            'I' => 'SPUNBOND_TERPOTONG',
            'J' => 'TERBALIK',
            'K' => 'OVERCUTTING',
            'L' => 'SEWING_MIRING',
            'M' => 'MARGIN_OUT_DIMENSI',
        ];

        return [
            'CR6' => ['kodeLine' => 'PRESS', 'daftarDefect' => $templatePressCarpet],
            'CL7' => ['kodeLine' => 'PRESS', 'daftarDefect' => $templatePressCarpet],
            'CB9' => ['kodeLine' => 'PRESS', 'daftarDefect' => $templateCb9],
            'BT136' => ['kodeLine' => 'PRESS', 'daftarDefect' => $templatePressUmum],
            'BT137' => ['kodeLine' => 'PRESS', 'daftarDefect' => $templatePressUmum],
            'BT144' => ['kodeLine' => 'PRESS', 'daftarDefect' => $templatePressUmum],
            'BM7' => ['kodeLine' => 'PRESS', 'daftarDefect' => $templatePressUmum],
            'BM8' => ['kodeLine' => 'PRESS', 'daftarDefect' => $templatePressUmum],
            'FJ0' => ['kodeLine' => 'PRESS', 'daftarDefect' => $templatePressUmum],
            'FJ1' => ['kodeLine' => 'PRESS', 'daftarDefect' => $templatePressUmum],
            'FSB' => ['kodeLine' => 'SEWING', 'daftarDefect' => $templateFeltSewing],
            'CFRSH' => ['kodeLine' => 'SEWING', 'daftarDefect' => $templateCoverSewing],
            'PRSB_RH_070' => ['kodeLine' => 'SEWING', 'daftarDefect' => $templateProtector],
            'PRSB_LH_080' => ['kodeLine' => 'SEWING', 'daftarDefect' => $templateProtector],
            'PRSB_RH_090' => ['kodeLine' => 'SEWING', 'daftarDefect' => $templateProtector],
            'PRSB_LH_100' => ['kodeLine' => 'SEWING', 'daftarDefect' => $templateProtector],
            'PRSB_RH_110' => ['kodeLine' => 'SEWING', 'daftarDefect' => $templateProtector],
            'PRSB_LH_120' => ['kodeLine' => 'SEWING', 'daftarDefect' => $templateProtector],
        ];
    }
}
