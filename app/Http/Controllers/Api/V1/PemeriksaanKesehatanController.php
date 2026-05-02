<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1;

use App\Application\Kesehatan\MembacaStatusKesehatanServer;
use App\Http\Controllers\Controller;
use App\Http\Requests\Api\V1\MembacaStatusKesehatanRequest;
use App\Http\Resources\Api\V1\StatusKesehatanResource;
use App\Support\Api\ResponApi;
use App\Support\Errors\KodeKesalahanApi;
use Illuminate\Http\JsonResponse;

final class PemeriksaanKesehatanController extends Controller
{
    public function __construct(
        private MembacaStatusKesehatanServer $membacaStatusKesehatanServer,
    ) {
    }

    public function __invoke(MembacaStatusKesehatanRequest $permintaan): JsonResponse
    {
        $statusKesehatanServer = $this->membacaStatusKesehatanServer->jalankan();
        $data = (new StatusKesehatanResource($statusKesehatanServer))->resolve($permintaan);

        if ($statusKesehatanServer->databaseTersedia()) {
            return ResponApi::berhasil(
                pesan: 'Server berjalan normal',
                data: $data,
            );
        }

        return ResponApi::gagal(
            pesan: 'Server berjalan, tetapi koneksi database belum tersedia',
            kodeKesalahan: KodeKesalahanApi::DATABASE_TIDAK_TERHUBUNG,
            detailKesalahan: [],
            data: $data,
            statusHttp: 503,
        );
    }
}
