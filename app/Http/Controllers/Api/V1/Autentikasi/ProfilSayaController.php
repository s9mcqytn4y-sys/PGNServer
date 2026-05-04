<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1\Autentikasi;

use App\Http\Controllers\Controller;
use App\Http\Resources\Api\V1\Autentikasi\ProfilPenggunaResource;
use App\Support\Api\ResponApi;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

final class ProfilSayaController extends Controller
{
    public function __invoke(Request $permintaan): JsonResponse
    {
        $pengguna = $permintaan->user();
        $data = (new ProfilPenggunaResource($pengguna))->resolve($permintaan);

        return ResponApi::berhasil(
            pesan: 'Profil pengguna berhasil dimuat',
            data: $data,
        );
    }
}
