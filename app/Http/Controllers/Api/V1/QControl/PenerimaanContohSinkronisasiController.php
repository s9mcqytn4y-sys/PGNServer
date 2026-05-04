<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1\QControl;

use App\Application\Idempotency\MengelolaIdempotencyKey;
use App\Http\Controllers\Controller;
use App\Http\Requests\Api\V1\QControl\PenerimaanContohSinkronisasiRequest;
use App\Http\Resources\Api\V1\QControl\ContohSinkronisasiResource;
use App\Support\Api\ResponApi;
use App\Support\Errors\KodeKesalahanApi;
use Illuminate\Http\JsonResponse;

final class PenerimaanContohSinkronisasiController extends Controller
{
    public function __construct(
        private MengelolaIdempotencyKey $mengelolaIdempotencyKey,
    ) {}

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

        $endpoint = $permintaan->getPathInfo();
        $hasilPemeriksaanIdempotency = $this->mengelolaIdempotencyKey->periksaAtauSimpanPenerimaan(
            kunciIdempotency: $idempotencyKey,
            metodeHttp: $permintaan->method(),
            endpoint: $endpoint,
            hashPayload: $permintaan->hashPayload(),
            sumberAplikasi: 'QControl',
        );

        if ($hasilPemeriksaanIdempotency->sudahAda) {
            return ResponApi::berhasil(
                pesan: 'Payload sinkronisasi QControl sudah pernah diterima',
                data: (new ContohSinkronisasiResource([
                    'diterima' => true,
                    'duplikat' => true,
                    'idempotencyKey' => $idempotencyKey,
                    'endpoint' => $endpoint,
                    'mode' => 'kontrak_awal',
                ]))->resolve($permintaan),
            );
        }

        $data = (new ContohSinkronisasiResource([
            'diterima' => true,
            'duplikat' => false,
            'idempotencyKey' => $idempotencyKey,
            'endpoint' => $endpoint,
            'mode' => 'kontrak_awal',
        ]))->resolve($permintaan);

        $respon = ResponApi::berhasil(
            pesan: 'Payload sinkronisasi QControl diterima',
            data: $data,
        );

        $this->mengelolaIdempotencyKey->simpanResponseBerhasil(
            kunciIdempotency: $idempotencyKey,
            statusHttp: $respon->getStatusCode(),
            responseBody: $respon->getData(true),
        );

        return $respon;
    }
}
