<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1\QControl;

use App\Application\QControl\MembacaDetailPemeriksaanHarian;
use App\Http\Controllers\Controller;
use App\Http\Resources\Api\V1\QControl\DetailPemeriksaanHarianResource;
use App\Models\QControlPemeriksaanHarian;
use App\Support\Api\ResponApi;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

final class DetailPemeriksaanHarianController extends Controller
{
    public function __construct(
        private MembacaDetailPemeriksaanHarian $membacaDetailPemeriksaanHarian,
    ) {}

    public function __invoke(Request $permintaan, QControlPemeriksaanHarian $pemeriksaanHarian): JsonResponse
    {
        $pemeriksaanHarian = $this->membacaDetailPemeriksaanHarian->jalankan($pemeriksaanHarian);

        return ResponApi::berhasil(
            pesan: 'Detail pemeriksaan harian QControl berhasil dimuat',
            data: (new DetailPemeriksaanHarianResource($pemeriksaanHarian))->resolve($permintaan),
        );
    }
}
