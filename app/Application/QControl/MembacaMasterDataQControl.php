<?php

declare(strict_types=1);

namespace App\Application\QControl;

use App\Models\QControlJenisDefect;
use App\Models\QControlKategoriDefect;
use App\Models\QControlLineProduksi;
use App\Models\QControlMaterial;
use App\Models\QControlPart;
use App\Models\QControlPartJenisDefect;
use App\Models\QControlSlotWaktu;

final class MembacaMasterDataQControl
{
    private const VERSI_MASTER_DATA = '2026.05.2F-A';

    /**
     * @return array{
     *     data: array<string, mixed>,
     *     metadata: array<string, int>
     * }
     */
    public function jalankan(): array
    {
        $lineProduksi = QControlLineProduksi::query()
            ->orderBy('urutan_tampil')
            ->orderBy('kode_line')
            ->get();

        $slotWaktu = QControlSlotWaktu::query()
            ->orderBy('urutan_tampil')
            ->orderBy('kode_slot')
            ->get();

        $material = QControlMaterial::query()
            ->orderBy('nama_material')
            ->get();

        $kategoriDefect = QControlKategoriDefect::query()
            ->orderBy('urutan_tampil')
            ->orderBy('kode_kategori')
            ->get();

        $jenisDefect = QControlJenisDefect::query()
            ->with('kategoriDefectTerkait')
            ->orderBy('kode_defect')
            ->get();

        $part = QControlPart::query()
            ->with(['materialTerkait', 'lineProduksiDefault'])
            ->orderBy('kode_unik_part')
            ->get();

        $relasiPartDefect = QControlPartJenisDefect::query()
            ->with(['partTerkait', 'jenisDefectTerkait.kategoriDefectTerkait'])
            ->orderBy('urutan_tampil')
            ->get();

        return [
            'data' => [
                'versiMasterData' => self::VERSI_MASTER_DATA,
                'lineProduksi' => $lineProduksi,
                'slotWaktu' => $slotWaktu,
                'material' => $material,
                'part' => $part,
                'kategoriDefect' => $kategoriDefect,
                'jenisDefect' => $jenisDefect,
                'relasiPartDefect' => $relasiPartDefect,
            ],
            'metadata' => [
                'jumlahLineProduksi' => $lineProduksi->count(),
                'jumlahSlotWaktu' => $slotWaktu->count(),
                'jumlahMaterial' => $material->count(),
                'jumlahPart' => $part->count(),
                'jumlahJenisDefect' => $jenisDefect->count(),
                'jumlahRelasiPartDefect' => $relasiPartDefect->count(),
                'jumlahShiftOperasional' => 1,
            ],
        ];
    }
}
