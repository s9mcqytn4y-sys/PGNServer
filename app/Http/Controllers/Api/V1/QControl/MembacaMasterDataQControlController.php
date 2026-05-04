<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1\QControl;

use App\Application\QControl\MembacaMasterDataQControl;
use App\Http\Controllers\Controller;
use App\Http\Resources\Api\V1\QControl\MasterDataQControlResource;
use App\Support\Api\ResponApi;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

final class MembacaMasterDataQControlController extends Controller
{
    public function __construct(
        private MembacaMasterDataQControl $membacaMasterDataQControl,
    ) {}

    public function __invoke(Request $permintaan): JsonResponse
    {
        $hasilPembacaan = $this->membacaMasterDataQControl->jalankan();
        $data = (new MasterDataQControlResource($hasilPembacaan['data']))->resolve($permintaan);

        return ResponApi::berhasil(
            pesan: 'Master data QControl berhasil dimuat',
            data: $data,
            metadata: $hasilPembacaan['metadata'],
        );
    }
}
