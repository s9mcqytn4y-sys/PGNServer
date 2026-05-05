<?php

declare(strict_types=1);

namespace App\Application\QControl;

use App\Models\QControlLineProduksi;
use App\Models\QControlPart;
use App\Models\QControlPartJenisDefect;
use App\Models\QControlPemeriksaanDefectSlot;
use App\Models\QControlPemeriksaanHarian;
use App\Models\QControlPemeriksaanPart;
use App\Models\QControlSlotWaktu;
use App\Models\User;
use App\Support\Errors\KodeKesalahanApi;
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

        $daftarPartId = collect($daftarPartPayload)
            ->pluck('partId')
            ->filter(fn (mixed $partId): bool => is_string($partId))
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
            ->get();

        /** @var KoleksiEloquent<int, QControlSlotWaktu> $koleksiSlotWaktu */
        $koleksiSlotWaktu = QControlSlotWaktu::query()
            ->whereIn('id', $daftarSlotWaktuId)
            ->get();

        /** @var KoleksiEloquent<int, QControlPartJenisDefect> $koleksiRelasiPartDefect */
        $koleksiRelasiPartDefect = QControlPartJenisDefect::query()
            ->whereIn('id', $daftarRelasiPartDefectId)
            ->with('jenisDefectTerkait')
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

            $totalCheckHarian += $totalCheckPart;
            $totalDefectHarian += $totalDefectPart;

            $ringkasanPart[] = [
                'part' => $part,
                'totalCheck' => $totalCheckPart,
                'totalOk' => $totalOkPart,
                'totalDefect' => $totalDefectPart,
                'rasioDefect' => $rasioDefectPart,
                'urutanTampil' => $indeksPart + 1,
                'daftarDefect' => $barisDefectPart,
            ];
        }

        $totalOkHarian = $totalCheckHarian - $totalDefectHarian;
        $rasioDefectHarian = $this->hitungRasioDefect($totalDefectHarian, $totalCheckHarian);

        return DB::transaction(function () use (
            $payload,
            $penggunaHeadQC,
            $lineProduksi,
            $ringkasanPart,
            $kunciIdempotency,
            $hashPayload,
            $totalCheckHarian,
            $totalOkHarian,
            $totalDefectHarian,
            $rasioDefectHarian,
            $jumlahBarisDefect,
        ): array {
            $pemeriksaanHarianTersimpan = QControlPemeriksaanHarian::query()
                ->with(['lineProduksi', 'daftarPemeriksaanPart.daftarDefectSlot'])
                ->where('tanggal_produksi', $payload['tanggalProduksi'])
                ->where('line_produksi_id', $lineProduksi->id)
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
                    pesan: 'Pemeriksaan harian untuk tanggal produksi dan line tersebut sudah ada',
                    kodeKesalahan: KodeKesalahanApi::PEMERIKSAAN_HARIAN_SUDAH_ADA,
                    detailKesalahan: [
                        [
                            'field' => 'tanggalProduksi',
                            'pesan' => 'Pemeriksaan harian untuk tanggal produksi dan line tersebut sudah ada',
                        ],
                    ],
                    statusHttp: 409,
                );
            }

            $pemeriksaanHarian = QControlPemeriksaanHarian::query()->create([
                'id' => (string) Str::uuid(),
                'tanggal_produksi' => $payload['tanggalProduksi'],
                'line_produksi_id' => $lineProduksi->id,
                'nomor_dokumen' => $payload['nomorDokumen'] ?? 'FM-QA-025',
                'revisi' => $payload['revisi'] ?? '1',
                'pengguna_headqc_id' => $penggunaHeadQC->id,
                'client_draft_id' => $payload['clientDraftId'] ?? null,
                'idempotency_key' => $kunciIdempotency,
                'hash_payload' => $hashPayload,
                'status' => 'DITERIMA',
                'total_check' => $totalCheckHarian,
                'total_ok' => $totalOkHarian,
                'total_defect' => $totalDefectHarian,
                'rasio_defect' => $rasioDefectHarian,
                'catatan' => $payload['catatan'] ?? null,
                'diterima_pada' => Carbon::now(),
            ]);

            foreach ($ringkasanPart as $dataPart) {
                $pemeriksaanPart = QControlPemeriksaanPart::query()->create([
                    'id' => (string) Str::uuid(),
                    'pemeriksaan_harian_id' => $pemeriksaanHarian->id,
                    'part_id' => $dataPart['part']->id,
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
                        'jenis_defect_id' => $dataDefect['relasiPartDefect']->jenis_defect_id,
                        'slot_waktu_id' => $dataDefect['slotWaktu']->id,
                        'jumlah_defect' => $dataDefect['jumlahDefect'],
                    ]);
                }
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

    private function hitungRasioDefect(int $totalDefect, int $totalCheck): float
    {
        if ($totalCheck === 0) {
            return 0.0;
        }

        return round(($totalDefect / $totalCheck) * 100, 2);
    }
}
