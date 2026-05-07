<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1\QControl;

use App\Models\QControlPemeriksaanProduksiTanpaNg;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

final class PemeriksaanProduksiTanpaNgResource extends JsonResource
{
    /**
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        /** @var QControlPemeriksaanProduksiTanpaNg $produksiTanpaNg */
        $produksiTanpaNg = $this->resource;

        return [
            'id' => $produksiTanpaNg->id,
            'partId' => $produksiTanpaNg->part_id,
            'uniqNoPart' => $produksiTanpaNg->uniq_no_part,
            'nomorPartSnapshot' => $produksiTanpaNg->nomor_part_snapshot,
            'namaPartSnapshot' => $produksiTanpaNg->nama_part_snapshot,
            'totalProduksi' => (int) $produksiTanpaNg->total_produksi,
            'catatan' => $produksiTanpaNg->catatan,
        ];
    }
}
