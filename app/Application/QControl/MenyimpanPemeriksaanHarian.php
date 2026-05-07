<?php

declare(strict_types=1);

namespace App\Application\QControl;

use App\Models\QControlLineProduksi;
use App\Models\QControlPart;
use App\Models\QControlPartJenisDefect;
use App\Models\QControlPemeriksaanDefectSlot;
use App\Models\QControlPemeriksaanHarian;
use App\Models\QControlPemeriksaanPart;
use App\Models\QControlPemeriksaanProduksiTanpaNg;
use App\Models\QControlSlotWaktu;
use App\Models\User;
use App\Support\Errors\KodeKesalahanApi;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Collection as KoleksiEloquent;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Str;

final class MenyimpanPemeriksaanHarian
{
    /**
     * @param  array<string, mixed>  $payload
     * @return array{
     *     pemeriksaanHarian: QControlPemeriksaanHarian,
     *     jumlahPart: int,
     *     jumlahBarisDefect: int,
     *     duplikat: bool
     * }
     */
    public function jalankan(
        User $penggunaHeadQC,
        array $payload,
        string $kunciIdempotency,
        string $hashPayload,
    ): array {
        $lineProduksi = QControlLineProduksi::query()
            ->find($payload['lineProduksiId']);

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

        /** @var array<int, array<string, mixed>> $daftarPartPayload */
        $daftarPartPayload = $payload['daftarPart'];

        $this->validasiPartTidakDuplikat($daftarPartPayload);
        $this->validasiDefectSlotTidakDuplikat($daftarPartPayload);

        $daftarProduksiTanpaNgPayload = $payload['daftarProduksiTanpaNg'] ?? [];
        $this->validasiProduksiTanpaNgTidakDuplikat($daftarProduksiTanpaNgPayload);

        $daftarPartId = collect($daftarPartPayload)
            ->pluck('partId')
            ->merge(collect($daftarProduksiTanpaNgPayload)->pluck('partId'))
            ->filter(fn (mixed $partId): bool => is_string($partId))
            ->unique()
            ->values()
            ->all();

        $daftarSlotWaktuId = collect($daftarPartPayload)
            ->flatMap(fn (array $part): array => $part['daftarDefect'] ?? [])
            ->pluck('slotWaktuId')
            ->filter(fn (mixed $slotWaktuId): bool => is_string($slotWaktuId))
            ->unique()
            ->values()
            ->all();

        $daftarRelasiPartDefectId = collect($daftarPartPayload)
            ->flatMap(fn (array $part): array => $part['daftarDefect'] ?? [])
            ->pluck('relasiPartDefectId')
            ->filter(fn (mixed $relasiPartDefectId): bool => is_string($relasiPartDefectId))
            ->unique()
            ->values()
            ->all();

        /** @var KoleksiEloquent<int, QControlPart> $koleksiPart */
        $koleksiPart = QControlPart::query()
            ->whereIn('id', $daftarPartId)
            ->with('materialTerkait')
            ->get();

        /** @var KoleksiEloquent<int, QControlSlotWaktu> $koleksiSlotWaktu */
        $koleksiSlotWaktu = QControlSlotWaktu::query()
            ->whereIn('id', $daftarSlotWaktuId)
            ->get();

        /** @var KoleksiEloquent<int, QControlPartJenisDefect> $koleksiRelasiPartDefect */
        $koleksiRelasiPartDefect = QControlPartJenisDefect::query()
            ->whereIn('id', $daftarRelasiPartDefectId)
            ->with('jenisDefectTerkait.kategoriDefectTerkait')
            ->get();

        /** @var array<string, QControlPart> $partBerdasarkanId */
        $partBerdasarkanId = $koleksiPart->keyBy('id')->all();

        /** @var array<string, QControlSlotWaktu> $slotWaktuBerdasarkanId */
        $slotWaktuBerdasarkanId = $koleksiSlotWaktu->keyBy('id')->all();

        /** @var array<string, QControlPartJenisDefect> $relasiPartDefectBerdasarkanId */
        $relasiPartDefectBerdasarkanId = $koleksiRelasiPartDefect->keyBy('id')->all();

        $ringkasanPart = [];
        $totalCheckHarian = 0;
        $totalDefectHarian = 0;
        $jumlahBarisDefect = 0;

        foreach ($daftarPartPayload as $indeksPart => $dataPart) {
            $part = $partBerdasarkanId[$dataPart['partId']] ?? null;

            if (! $part instanceof QControlPart || ! $part->aktif) {
                throw new PengecualianPemeriksaanHarian(
                    pesan: 'Part QControl tidak aktif atau tidak tersedia',
                    kodeKesalahan: KodeKesalahanApi::VALIDASI_GAGAL,
                    detailKesalahan: [
                        [
                            'field' => "daftarPart.$indeksPart.partId",
                            'pesan' => 'Part QControl tidak aktif atau tidak tersedia',
                        ],
                    ],
                );
            }

            if ($part->line_default_id !== $lineProduksi->id) {
                throw new PengecualianPemeriksaanHarian(
                    pesan: 'Part tidak sesuai dengan line produksi yang dipilih',
                    kodeKesalahan: KodeKesalahanApi::TEMPLATE_DEFECT_TIDAK_VALID,
                    detailKesalahan: [
                        [
                            'field' => "daftarPart.$indeksPart.partId",
                            'pesan' => 'Part tidak sesuai dengan line produksi yang dipilih',
                        ],
                    ],
                );
            }

            $totalCheckPart = (int) $dataPart['totalCheck'];
            /** @var array<int, array<string, mixed>> $daftarDefectPart */
            $daftarDefectPart = $dataPart['daftarDefect'] ?? [];
            $totalDefectPart = 0;
            $barisDefectPart = [];

            foreach ($daftarDefectPart as $indeksDefect => $dataDefect) {
                $slotWaktu = $slotWaktuBerdasarkanId[$dataDefect['slotWaktuId']] ?? null;

                if (! $slotWaktu instanceof QControlSlotWaktu || ! $slotWaktu->aktif) {
                    throw new PengecualianPemeriksaanHarian(
                        pesan: 'Slot waktu QControl tidak aktif atau tidak tersedia',
                        kodeKesalahan: KodeKesalahanApi::TEMPLATE_DEFECT_TIDAK_VALID,
                        detailKesalahan: [
                            [
                                'field' => "daftarPart.$indeksPart.daftarDefect.$indeksDefect.slotWaktuId",
                                'pesan' => 'Slot waktu QControl tidak aktif atau tidak tersedia',
                            ],
                        ],
                    );
                }

                $relasiPartDefect = $relasiPartDefectBerdasarkanId[$dataDefect['relasiPartDefectId']] ?? null;

                if (
                    ! $relasiPartDefect instanceof QControlPartJenisDefect
                    || ! $relasiPartDefect->aktif
                    || $relasiPartDefect->part_id !== $part->id
                    || ! $relasiPartDefect->jenisDefectTerkait?->aktif
                ) {
                    throw new PengecualianPemeriksaanHarian(
                        pesan: 'Template defect QControl tidak valid untuk part ini',
                        kodeKesalahan: KodeKesalahanApi::TEMPLATE_DEFECT_TIDAK_VALID,
                        detailKesalahan: [
                            [
                                'field' => "daftarPart.$indeksPart.daftarDefect.$indeksDefect.relasiPartDefectId",
                                'pesan' => 'Template defect QControl tidak valid untuk part ini',
                            ],
                        ],
                    );
                }

                $jumlahDefect = (int) $dataDefect['jumlahDefect'];
                $totalDefectPart += $jumlahDefect;
                $jumlahBarisDefect++;

                $barisDefectPart[] = [
                    'relasiPartDefect' => $relasiPartDefect,
                    'slotWaktu' => $slotWaktu,
                    'kategoriDefectSnapshot' => $relasiPartDefect->jenisDefectTerkait?->kategoriDefectTerkait?->nama_kategori,
                    'jumlahDefect' => $jumlahDefect,
                ];
            }

            if ($totalDefectPart > $totalCheckPart) {
                throw new PengecualianPemeriksaanHarian(
                    pesan: 'Total defect tidak boleh melebihi total check',
                    kodeKesalahan: KodeKesalahanApi::TOTAL_DEFECT_MELEBIHI_TOTAL_CHECK,
                    detailKesalahan: [
                        [
                            'field' => "daftarPart.$indeksPart.totalCheck",
                            'pesan' => 'Total defect tidak boleh melebihi total check',
                        ],
                    ],
                );
            }

            $totalOkPart = $totalCheckPart - $totalDefectPart;
            $rasioDefectPart = $this->hitungRasioDefect($totalDefectPart, $totalCheckPart);
            $kategoriNgSnapshot = collect($barisDefectPart)
                ->pluck('kategoriDefectSnapshot')
                ->filter(fn (mixed $kategori): bool => is_string($kategori) && $kategori !== '')
                ->unique()
                ->implode(', ');

            $totalCheckHarian += $totalCheckPart;
            $totalDefectHarian += $totalDefectPart;

            $ringkasanPart[] = [
                'part' => $part,
                'totalCheck' => $totalCheckPart,
                'totalOk' => $totalOkPart,
                'totalDefect' => $totalDefectPart,
                'rasioDefect' => $rasioDefectPart,
                'kategoriNgSnapshot' => $kategoriNgSnapshot !== '' ? $kategoriNgSnapshot : null,
                'urutanTampil' => $indeksPart + 1,
                'daftarDefect' => $barisDefectPart,
            ];
        }

        $totalOkHarian = $totalCheckHarian - $totalDefectHarian;
        $rasioDefectHarian = $this->hitungRasioDefect($totalDefectHarian, $totalCheckHarian);

        $ringkasanProduksiTanpaNg = [];
        foreach ($daftarProduksiTanpaNgPayload as $indeksTanpaNg => $dataTanpaNg) {
            $part = $partBerdasarkanId[$dataTanpaNg['partId']] ?? null;

            if (! $part instanceof QControlPart || ! $part->aktif) {
                throw new PengecualianPemeriksaanHarian(
                    pesan: 'Part QControl tidak aktif atau tidak tersedia (Seksi Tanpa NG)',
                    kodeKesalahan: KodeKesalahanApi::VALIDASI_GAGAL,
                    detailKesalahan: [
                        [
                            'field' => "daftarProduksiTanpaNg.$indeksTanpaNg.partId",
                            'pesan' => 'Part QControl tidak aktif atau tidak tersedia',
                        ],
                    ],
                );
            }

            if ($part->line_default_id !== $lineProduksi->id) {
                throw new PengecualianPemeriksaanHarian(
                    pesan: 'Part tidak sesuai dengan line produksi yang dipilih (Seksi Tanpa NG)',
                    kodeKesalahan: KodeKesalahanApi::TEMPLATE_DEFECT_TIDAK_VALID,
                    detailKesalahan: [
                        [
                            'field' => "daftarProduksiTanpaNg.$indeksTanpaNg.partId",
                            'pesan' => 'Part tidak sesuai dengan line produksi yang dipilih',
                        ],
                    ],
                );
            }

            $ringkasanProduksiTanpaNg[] = [
                'part' => $part,
                'totalProduksi' => (int) $dataTanpaNg['totalProduksi'],
                'catatan' => $dataTanpaNg['catatan'] ?? null,
            ];
        }

        return DB::transaction(function () use (
            $payload,
            $penggunaHeadQC,
            $lineProduksi,
            $ringkasanPart,
            $ringkasanProduksiTanpaNg,
            $kunciIdempotency,
            $hashPayload,
            $totalCheckHarian,
            $totalOkHarian,
            $totalDefectHarian,
            $rasioDefectHarian,
            $jumlahBarisDefect,
        ): array {
            $pemeriksaanHarianTersimpan = QControlPemeriksaanHarian::query()
                ->with([
                    'lineProduksi',
                    'daftarPemeriksaanPart' => fn (Builder $query): Builder => $query->with('daftarDefectSlot'),
                ])
                ->where('idempotency_key', $kunciIdempotency)
                ->lockForUpdate()
                ->first();

            if ($pemeriksaanHarianTersimpan instanceof QControlPemeriksaanHarian) {
                if ($pemeriksaanHarianTersimpan->hash_payload === $hashPayload) {
                    return [
                        'pemeriksaanHarian' => $pemeriksaanHarianTersimpan,
                        'jumlahPart' => $pemeriksaanHarianTersimpan->daftarPemeriksaanPart->count(),
                        'jumlahBarisDefect' => $pemeriksaanHarianTersimpan->daftarPemeriksaanPart
                            ->sum(fn (QControlPemeriksaanPart $pemeriksaanPart): int => $pemeriksaanPart->daftarDefectSlot->count()),
                        'duplikat' => true,
                    ];
                }

                throw new PengecualianPemeriksaanHarian(
                    pesan: 'Idempotency key sudah digunakan untuk payload berbeda',
                    kodeKesalahan: KodeKesalahanApi::KONFLIK_IDEMPOTENCY,
                    detailKesalahan: [
                        [
                            'field' => 'X-Idempotency-Key',
                            'pesan' => 'Idempotency key sudah digunakan untuk payload berbeda',
                        ],
                    ],
                    statusHttp: 409,
                );
            }

            $namaPicSnapshot = $penggunaHeadQC->name;
            $emailPicSnapshot = $penggunaHeadQC->email;

            $pemeriksaanHarian = QControlPemeriksaanHarian::query()->create([
                'id' => (string) Str::uuid(),
                'tanggal_produksi' => $payload['tanggalProduksi'],
                'line_produksi_id' => $lineProduksi->id,
                'kode_line_snapshot' => $lineProduksi->kode_line,
                'nama_line_snapshot' => $lineProduksi->nama_line,
                'nomor_dokumen_snapshot' => $payload['nomorDokumen'] ?? 'FM-QA-025',
                'revisi_dokumen_snapshot' => $payload['revisi'] ?? '1',
                'pengguna_headqc_id' => $penggunaHeadQC->id,
                'nama_pic_snapshot' => $namaPicSnapshot,
                'email_pic_snapshot' => $emailPicSnapshot,
                'client_draft_id' => $payload['clientDraftId'] ?? null,
                'idempotency_key' => $kunciIdempotency,
                'hash_payload' => $hashPayload,
                'status' => 'DITERIMA',
                'total_check' => $totalCheckHarian,
                'total_ok' => $totalOkHarian,
                'total_defect' => $totalDefectHarian,
                'rasio_defect' => $rasioDefectHarian,
                'catatan' => $payload['catatan'] ?? null,
                'disiapkan_oleh_snapshot' => $namaPicSnapshot,
                'diperiksa_oleh_snapshot' => $namaPicSnapshot,
                'disetujui_oleh_snapshot' => $namaPicSnapshot,
                'diterima_pada' => Carbon::now(),
            ]);

            foreach ($ringkasanPart as $dataPart) {
                $pemeriksaanPart = QControlPemeriksaanPart::query()->create([
                    'id' => (string) Str::uuid(),
                    'pemeriksaan_harian_id' => $pemeriksaanHarian->id,
                    'part_id' => $dataPart['part']->id,
                    'kode_unik_part_snapshot' => $dataPart['part']->kode_unik_part,
                    'nomor_part_snapshot' => $dataPart['part']->nomor_part,
                    'nama_part_snapshot' => $dataPart['part']->nama_part,
                    'material_id_snapshot' => $dataPart['part']->material_id,
                    'nama_material_snapshot' => $dataPart['part']->materialTerkait?->nama_material,
                    'kategori_ng_snapshot' => $dataPart['kategoriNgSnapshot'],
                    'total_check' => $dataPart['totalCheck'],
                    'total_ok' => $dataPart['totalOk'],
                    'total_defect' => $dataPart['totalDefect'],
                    'rasio_defect' => $dataPart['rasioDefect'],
                    'urutan_tampil' => $dataPart['urutanTampil'],
                ]);

                foreach ($dataPart['daftarDefect'] as $dataDefect) {
                    QControlPemeriksaanDefectSlot::query()->create([
                        'id' => (string) Str::uuid(),
                        'pemeriksaan_part_id' => $pemeriksaanPart->id,
                        'relasi_part_defect_id' => $dataDefect['relasiPartDefect']->id,
                        'kode_tampilan_defect_snapshot' => $dataDefect['relasiPartDefect']->kode_tampilan_defect,
                        'jenis_defect_id' => $dataDefect['relasiPartDefect']->jenis_defect_id,
                        'kode_defect_snapshot' => $dataDefect['relasiPartDefect']->jenisDefectTerkait?->kode_defect,
                        'nama_defect_snapshot' => $dataDefect['relasiPartDefect']->jenisDefectTerkait?->nama_defect,
                        'kategori_defect_snapshot' => $dataDefect['kategoriDefectSnapshot'],
                        'slot_waktu_id' => $dataDefect['slotWaktu']->id,
                        'kode_slot_snapshot' => $dataDefect['slotWaktu']->kode_slot,
                        'label_slot_snapshot' => $dataDefect['slotWaktu']->label_slot,
                        'jam_mulai_snapshot' => $dataDefect['slotWaktu']->jam_mulai?->format('H:i:s'),
                        'jam_selesai_snapshot' => $dataDefect['slotWaktu']->jam_selesai?->format('H:i:s'),
                        'jumlah_defect' => $dataDefect['jumlahDefect'],
                    ]);
                }
            }

            foreach ($ringkasanProduksiTanpaNg as $dataTanpaNg) {
                QControlPemeriksaanProduksiTanpaNg::query()->create([
                    'id' => (string) Str::uuid(),
                    'pemeriksaan_harian_id' => $pemeriksaanHarian->id,
                    'part_id' => $dataTanpaNg['part']->id,
                    'uniq_no_part' => $dataTanpaNg['part']->kode_unik_part,
                    'nomor_part_snapshot' => $dataTanpaNg['part']->nomor_part,
                    'nama_part_snapshot' => $dataTanpaNg['part']->nama_part,
                    'total_produksi' => $dataTanpaNg['totalProduksi'],
                    'catatan' => $dataTanpaNg['catatan'],
                ]);
            }

            $pemeriksaanHarian->load('lineProduksi');

            return [
                'pemeriksaanHarian' => $pemeriksaanHarian,
                'jumlahPart' => count($ringkasanPart),
                'jumlahBarisDefect' => $jumlahBarisDefect,
                'duplikat' => false,
            ];
        });
    }

    /**
     * @param  array<int, array<string, mixed>>  $daftarPartPayload
     */
    private function validasiPartTidakDuplikat(array $daftarPartPayload): void
    {
        $partYangSudahAda = [];

        foreach ($daftarPartPayload as $indeksPart => $dataPart) {
            $partId = (string) $dataPart['partId'];

            if (array_key_exists($partId, $partYangSudahAda)) {
                throw new PengecualianPemeriksaanHarian(
                    pesan: 'Part duplikat tidak diizinkan dalam satu payload pemeriksaan',
                    kodeKesalahan: KodeKesalahanApi::VALIDASI_GAGAL,
                    detailKesalahan: [
                        [
                            'field' => "daftarPart.$indeksPart.partId",
                            'pesan' => 'Part duplikat tidak diizinkan dalam satu payload pemeriksaan',
                        ],
                    ],
                );
            }

            $partYangSudahAda[$partId] = true;
        }
    }

    /**
     * @param  array<int, array<string, mixed>>  $daftarPartPayload
     */
    private function validasiDefectSlotTidakDuplikat(array $daftarPartPayload): void
    {
        foreach ($daftarPartPayload as $indeksPart => $dataPart) {
            $kombinasiYangSudahAda = [];

            foreach (($dataPart['daftarDefect'] ?? []) as $indeksDefect => $dataDefect) {
                $kombinasi = $dataDefect['relasiPartDefectId'].'|'.$dataDefect['slotWaktuId'];

                if (array_key_exists($kombinasi, $kombinasiYangSudahAda)) {
                    throw new PengecualianPemeriksaanHarian(
                        pesan: 'Kombinasi defect dan slot waktu duplikat tidak diizinkan dalam part yang sama',
                        kodeKesalahan: KodeKesalahanApi::VALIDASI_GAGAL,
                        detailKesalahan: [
                            [
                                'field' => "daftarPart.$indeksPart.daftarDefect.$indeksDefect.relasiPartDefectId",
                                'pesan' => 'Kombinasi defect dan slot waktu duplikat tidak diizinkan dalam part yang sama',
                            ],
                        ],
                    );
                }

                $kombinasiYangSudahAda[$kombinasi] = true;
            }
        }
    }

    /**
     * @param  array<int, array<string, mixed>>  $daftarProduksiTanpaNgPayload
     */
    private function validasiProduksiTanpaNgTidakDuplikat(array $daftarProduksiTanpaNgPayload): void
    {
        $partYangSudahAda = [];

        foreach ($daftarProduksiTanpaNgPayload as $indeksTanpaNg => $dataTanpaNg) {
            $partId = (string) $dataTanpaNg['partId'];

            if (array_key_exists($partId, $partYangSudahAda)) {
                throw new PengecualianPemeriksaanHarian(
                    pesan: 'Part duplikat tidak diizinkan dalam seksi Produksi Tanpa NG',
                    kodeKesalahan: KodeKesalahanApi::VALIDASI_GAGAL,
                    detailKesalahan: [
                        [
                            'field' => "daftarProduksiTanpaNg.$indeksTanpaNg.partId",
                            'pesan' => 'Part duplikat tidak diizinkan dalam seksi Produksi Tanpa NG',
                        ],
                    ],
                );
            }

            $partYangSudahAda[$partId] = true;
        }
    }

    private function hitungRasioDefect(int $totalDefect, int $totalCheck): float
    {
        if ($totalCheck === 0) {
            return 0.0;
        }

        return round(($totalDefect / $totalCheck) * 100, 2);
    }
}
