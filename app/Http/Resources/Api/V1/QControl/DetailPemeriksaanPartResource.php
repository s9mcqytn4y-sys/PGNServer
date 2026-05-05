<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1\QControl;

use App\Models\QControlPemeriksaanPart;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

final class DetailPemeriksaanPartResource extends JsonResource
{
    /**
     * @return array{
     *     id: string,
     *     partId: string,
     *     kodeUnikPartSnapshot: string,
     *     nomorPartSnapshot: string,
     *     namaPartSnapshot: string,
     *     namaMaterialSnapshot: string,
     *     totalCheck: int,
     *     totalOk: int,
     *     totalDefect: int,
     *     rasioDefect: float,
     *     urutanTampil: int,
     *     daftarDefectSlot: array<int, array<string, mixed>>
     * }
     */
    public function toArray(Request $request): array
    {
        /** @var QControlPemeriksaanPart $pemeriksaanPart */
        $pemeriksaanPart = $this->resource;

        return [
            'id' => $pemeriksaanPart->id,
            'partId' => (string) $pemeriksaanPart->part_id,
            'kodeUnikPartSnapshot' => (string) ($pemeriksaanPart->kode_unik_part_snapshot ?? $pemeriksaanPart->partTerkait?->kode_unik_part),
            'nomorPartSnapshot' => (string) ($pemeriksaanPart->nomor_part_snapshot ?? $pemeriksaanPart->partTerkait?->nomor_part),
            'namaPartSnapshot' => (string) ($pemeriksaanPart->nama_part_snapshot ?? $pemeriksaanPart->partTerkait?->nama_part),
            'namaMaterialSnapshot' => (string) ($pemeriksaanPart->nama_material_snapshot ?? $pemeriksaanPart->partTerkait?->materialTerkait?->nama_material),
            'totalCheck' => (int) $pemeriksaanPart->total_check,
            'totalOk' => (int) $pemeriksaanPart->total_ok,
            'totalDefect' => (int) $pemeriksaanPart->total_defect,
            'rasioDefect' => round((float) $pemeriksaanPart->rasio_defect, 2),
            'urutanTampil' => (int) $pemeriksaanPart->urutan_tampil,
            'daftarDefectSlot' => DetailPemeriksaanDefectSlotResource::collection($pemeriksaanPart->daftarDefectSlot)
                ->resolve($request),
        ];
    }
}
