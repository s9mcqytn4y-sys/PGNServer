<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1\QControl;

use App\Models\QControlPemeriksaanHarian;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

final class PemeriksaanHarianResource extends JsonResource
{
    /**
     * @return array{
     *     pemeriksaanHarianId: string,
     *     tanggalProduksi: string,
     *     lineProduksi: array{id: string, kodeLine: string, namaLine: string},
     *     totalCheck: int,
     *     totalOk: int,
     *     totalDefect: int,
     *     rasioDefect: float,
     *     jumlahPart: int,
     *     jumlahBarisDefect: int,
     *     duplikat: bool
     * }
     */
    public function toArray(Request $request): array
    {
        /** @var array{
         *     pemeriksaanHarian: QControlPemeriksaanHarian,
         *     jumlahPart: int,
         *     jumlahBarisDefect: int,
         *     duplikat: bool
         * } $data
         */
        $data = $this->resource;
        $pemeriksaanHarian = $data['pemeriksaanHarian'];

        return [
            'pemeriksaanHarianId' => $pemeriksaanHarian->id,
            'tanggalProduksi' => $pemeriksaanHarian->tanggal_produksi?->format('Y-m-d') ?? '',
            'lineProduksi' => [
                'id' => (string) $pemeriksaanHarian->lineProduksi?->id,
                'kodeLine' => (string) $pemeriksaanHarian->lineProduksi?->kode_line,
                'namaLine' => (string) $pemeriksaanHarian->lineProduksi?->nama_line,
            ],
            'totalCheck' => (int) $pemeriksaanHarian->total_check,
            'totalOk' => (int) $pemeriksaanHarian->total_ok,
            'totalDefect' => (int) $pemeriksaanHarian->total_defect,
            'rasioDefect' => round((float) $pemeriksaanHarian->rasio_defect, 2),
            'jumlahPart' => $data['jumlahPart'],
            'jumlahBarisDefect' => $data['jumlahBarisDefect'],
            'duplikat' => $data['duplikat'],
        ];
    }
}
