<?php

declare(strict_types=1);

namespace App\Application\QControl;

use App\Models\QControlPemeriksaanHarian;

final class MembacaDetailPemeriksaanHarian
{
    public function jalankan(QControlPemeriksaanHarian $pemeriksaanHarian): QControlPemeriksaanHarian
    {
        $pemeriksaanHarian->load([
            'lineProduksi',
            'daftarPemeriksaanPart' => fn ($query) => $query->orderBy('urutan_tampil'),
            'daftarPemeriksaanPart.partTerkait.materialTerkait',
            'daftarPemeriksaanPart.daftarDefectSlot' => fn ($query) => $query->orderBy('dibuat_pada'),
            'daftarPemeriksaanPart.daftarDefectSlot.relasiPartDefect.jenisDefectTerkait',
            'daftarPemeriksaanPart.daftarDefectSlot.slotWaktu',
        ]);

        return $pemeriksaanHarian;
    }
}
