<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1\QControl;

use App\Application\QControl\MembacaDaftarPemeriksaanHarian;
use App\Http\Controllers\Controller;
use App\Http\Requests\Api\V1\QControl\DaftarPemeriksaanHarianRequest;
use App\Http\Resources\Api\V1\QControl\RingkasanPemeriksaanHarianResource;
use App\Support\Api\ResponApi;
use Illuminate\Http\JsonResponse;

final class DaftarPemeriksaanHarianController extends Controller
{
    public function __construct(
        private MembacaDaftarPemeriksaanHarian $membacaDaftarPemeriksaanHarian,
    ) {}

    public function __invoke(DaftarPemeriksaanHarianRequest $permintaan): JsonResponse
    {
        $hasilPembacaan = $this->membacaDaftarPemeriksaanHarian->jalankan(
            $permintaan->filterTervalidasi(),
        );

        return ResponApi::berhasil(
            pesan: 'Daftar pemeriksaan harian QControl berhasil dimuat',
            data: [
                'daftarPemeriksaanHarian' => RingkasanPemeriksaanHarianResource::collection(
                    $hasilPembacaan['daftarPemeriksaanHarian'],
                )->resolve($permintaan),
            ],
            metadata: $hasilPembacaan['metadata'],
        );
    }
}
