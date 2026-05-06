<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1\QControl;

use App\Models\QControlPemeriksaanDefectSlot;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

final class DetailPemeriksaanDefectSlotResource extends JsonResource
{
    /**
     * @return array{
     *     id: string,
     *     relasiPartDefectId: string,
     *     jenisDefectId: string,
     *     slotWaktuId: string,
     *     kodeTampilanDefectSnapshot: string,
     *     kodeDefectSnapshot: string,
     *     namaDefectSnapshot: string,
     *     kategoriDefectSnapshot: string|null,
     *     kodeSlotSnapshot: string,
     *     labelSlotSnapshot: string,
     *     jamMulaiSnapshot: string|null,
     *     jamSelesaiSnapshot: string|null,
     *     jumlahDefect: int
     * }
     */
    public function toArray(Request $request): array
    {
        /** @var QControlPemeriksaanDefectSlot $defectSlot */
        $defectSlot = $this->resource;

        return [
            'id' => $defectSlot->id,
            'relasiPartDefectId' => (string) $defectSlot->relasi_part_defect_id,
            'jenisDefectId' => (string) $defectSlot->jenis_defect_id,
            'slotWaktuId' => (string) $defectSlot->slot_waktu_id,
            'kodeTampilanDefectSnapshot' => (string) ($defectSlot->kode_tampilan_defect_snapshot ?? $defectSlot->relasiPartDefect?->kode_tampilan_defect),
            'kodeDefectSnapshot' => (string) ($defectSlot->kode_defect_snapshot ?? $defectSlot->jenisDefect?->kode_defect ?? $defectSlot->relasiPartDefect?->jenisDefectTerkait?->kode_defect),
            'namaDefectSnapshot' => (string) ($defectSlot->nama_defect_snapshot ?? $defectSlot->jenisDefect?->nama_defect ?? $defectSlot->relasiPartDefect?->jenisDefectTerkait?->nama_defect),
            'kategoriDefectSnapshot' => $defectSlot->kategori_defect_snapshot,
            'kodeSlotSnapshot' => (string) ($defectSlot->kode_slot_snapshot ?? $defectSlot->slotWaktu?->kode_slot),
            'labelSlotSnapshot' => (string) ($defectSlot->label_slot_snapshot ?? $defectSlot->slotWaktu?->label_slot),
            'jamMulaiSnapshot' => $defectSlot->jam_mulai_snapshot ?? $defectSlot->slotWaktu?->jam_mulai?->format('H:i:s'),
            'jamSelesaiSnapshot' => $defectSlot->jam_selesai_snapshot ?? $defectSlot->slotWaktu?->jam_selesai?->format('H:i:s'),
            'jumlahDefect' => (int) $defectSlot->jumlah_defect,
        ];
    }
}
