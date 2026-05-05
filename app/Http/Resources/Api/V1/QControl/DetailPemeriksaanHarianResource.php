<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1\QControl;

use App\Models\QControlPemeriksaanHarian;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

final class DetailPemeriksaanHarianResource extends JsonResource
{
    /**
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        /** @var QControlPemeriksaanHarian $pemeriksaanHarian */
        $pemeriksaanHarian = $this->resource;

        return [
            'id' => $pemeriksaanHarian->id,
            'tanggalProduksi' => $pemeriksaanHarian->tanggal_produksi?->format('Y-m-d') ?? '',
            'nomorDokumen' => $pemeriksaanHarian->nomor_dokumen,
            'revisi' => $pemeriksaanHarian->revisi,
            'clientDraftId' => $pemeriksaanHarian->client_draft_id,
            'idempotencyKey' => $pemeriksaanHarian->idempotency_key,
            'status' => (string) $pemeriksaanHarian->status,
            'catatan' => $pemeriksaanHarian->catatan,
            'diterimaPada' => $pemeriksaanHarian->diterima_pada?->toIso8601String(),
            'lineProduksi' => [
                'id' => (string) $pemeriksaanHarian->line_produksi_id,
                'kodeLine' => (string) ($pemeriksaanHarian->kode_line_snapshot ?? $pemeriksaanHarian->lineProduksi?->kode_line),
                'namaLine' => (string) ($pemeriksaanHarian->nama_line_snapshot ?? $pemeriksaanHarian->lineProduksi?->nama_line),
            ],
            'totalCheck' => (int) $pemeriksaanHarian->total_check,
            'totalOk' => (int) $pemeriksaanHarian->total_ok,
            'totalDefect' => (int) $pemeriksaanHarian->total_defect,
            'rasioDefect' => round((float) $pemeriksaanHarian->rasio_defect, 2),
            'daftarPart' => DetailPemeriksaanPartResource::collection($pemeriksaanHarian->daftarPemeriksaanPart)
                ->resolve($request),
        ];
    }
}
