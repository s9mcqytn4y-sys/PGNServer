<?php

declare(strict_types=1);

namespace App\Support\Api;

use App\Support\Errors\KodeKesalahanApi;
use Illuminate\Http\JsonResponse;

final class ResponApi
{
    /**
     * @param  array<string, mixed>|null  $data
     * @param  array<string, mixed>|null  $metadata
     */
    public static function berhasil(
        string $pesan,
        ?array $data = null,
        ?array $metadata = null,
        int $statusHttp = 200,
    ): JsonResponse {
        return response()->json([
            'berhasil' => true,
            'pesan' => $pesan,
            'data' => $data,
            'metadata' => $metadata,
            'kesalahan' => null,
        ], $statusHttp);
    }

    /**
     * @param  array<int, array{field: string|null, pesan: string}>  $detailKesalahan
     * @param  array<string, mixed>|null  $data
     * @param  array<string, mixed>|null  $metadata
     */
    public static function gagal(
        string $pesan,
        KodeKesalahanApi $kodeKesalahan,
        array $detailKesalahan = [],
        ?array $data = null,
        ?array $metadata = null,
        int $statusHttp = 400,
    ): JsonResponse {
        return response()->json([
            'berhasil' => false,
            'pesan' => $pesan,
            'data' => $data,
            'metadata' => $metadata,
            'kesalahan' => [
                'kode' => $kodeKesalahan->value,
                'detail' => $detailKesalahan,
            ],
        ], $statusHttp);
    }
}
