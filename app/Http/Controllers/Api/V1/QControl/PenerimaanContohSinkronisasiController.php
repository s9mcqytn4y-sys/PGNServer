<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1\QControl;

use App\Http\Controllers\Controller;
use App\Http\Requests\Api\V1\QControl\PenerimaanContohSinkronisasiRequest;
use App\Http\Resources\Api\V1\QControl\ContohSinkronisasiResource;
use App\Support\Api\ResponApi;
use App\Support\Errors\KodeKesalahanApi;
use Illuminate\Http\JsonResponse;

final class PenerimaanContohSinkronisasiController extends Controller
{
    public function __invoke(PenerimaanContohSinkronisasiRequest $permintaan): JsonResponse
    {
        $idempotencyKey = $permintaan->idempotencyKey();

        if ($idempotencyKey === null) {
            return ResponApi::gagal(
                pesan: 'Header X-Idempotency-Key wajib diisi',
                kodeKesalahan: KodeKesalahanApi::VALIDASI_GAGAL,
                detailKesalahan: [
                    [
                        'field' => 'X-Idempotency-Key',
                        'pesan' => 'Header X-Idempotency-Key wajib diisi',
                    ],
                ],
                data: [
                    'payloadDiterima' => $permintaan->payloadDiterima(),
                ],
                statusHttp: 422,
            );
        }

        $data = (new ContohSinkronisasiResource([
            'diterima' => true,
            'idempotencyKey' => $idempotencyKey,
            'endpoint' => '/api/v1/qcontrol/contoh',
            'mode' => 'kontrak_awal',
        ]))->resolve($permintaan);

        return ResponApi::berhasil(
            pesan: 'Payload sinkronisasi QControl diterima',
            data: $data,
        );
    }
}
