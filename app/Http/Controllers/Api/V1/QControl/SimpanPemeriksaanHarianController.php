<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1\QControl;

use App\Application\Idempotency\MengelolaIdempotencyKey;
use App\Application\QControl\MenyimpanPemeriksaanHarian;
use App\Application\QControl\PengecualianPemeriksaanHarian;
use App\Http\Controllers\Controller;
use App\Http\Requests\Api\V1\QControl\SimpanPemeriksaanHarianRequest;
use App\Http\Resources\Api\V1\QControl\PemeriksaanHarianResource;
use App\Support\Api\ResponApi;
use App\Support\Errors\KodeKesalahanApi;
use Illuminate\Http\JsonResponse;

final class SimpanPemeriksaanHarianController extends Controller
{
    public function __construct(
        private MengelolaIdempotencyKey $mengelolaIdempotencyKey,
        private MenyimpanPemeriksaanHarian $menyimpanPemeriksaanHarian,
    ) {}

    public function __invoke(SimpanPemeriksaanHarianRequest $permintaan): JsonResponse
    {
        $idempotencyKey = (string) $permintaan->idempotencyKey();

        $hasilPemeriksaanIdempotency = $this->mengelolaIdempotencyKey->periksaAtauSimpanTransaksiSerius(
            kunciIdempotency: $idempotencyKey,
            metodeHttp: $permintaan->method(),
            endpoint: $permintaan->getPathInfo(),
            hashPayload: $permintaan->hashPayload(),
            sumberAplikasi: 'QControl',
        );

        if ($hasilPemeriksaanIdempotency->konflikPayload) {
            return ResponApi::gagal(
                pesan: 'Idempotency key sudah digunakan untuk payload berbeda',
                kodeKesalahan: KodeKesalahanApi::KONFLIK_IDEMPOTENCY,
                statusHttp: 409,
            );
        }

        if (
            $hasilPemeriksaanIdempotency->sudahAda
            && is_array($hasilPemeriksaanIdempotency->dataTersimpan)
            && is_int($hasilPemeriksaanIdempotency->statusHttpTersimpan)
        ) {
            $dataTersimpan = $hasilPemeriksaanIdempotency->dataTersimpan;

            if (isset($dataTersimpan['data']) && is_array($dataTersimpan['data'])) {
                $dataTersimpan['data']['duplikat'] = true;
            }

            return response()->json(
                $dataTersimpan,
                $hasilPemeriksaanIdempotency->statusHttpTersimpan,
            );
        }

        try {
            $hasilPenyimpanan = $this->menyimpanPemeriksaanHarian->jalankan(
                penggunaHeadQC: $permintaan->penggunaHeadQC(),
                payload: $permintaan->payloadTervalidasi(),
                kunciIdempotency: $idempotencyKey,
                hashPayload: $permintaan->hashPayload(),
            );
        } catch (PengecualianPemeriksaanHarian $pengecualian) {
            $this->mengelolaIdempotencyKey->batalkanPenerimaan($idempotencyKey);

            return ResponApi::gagal(
                pesan: $pengecualian->getMessage(),
                kodeKesalahan: $pengecualian->kodeKesalahan,
                detailKesalahan: $pengecualian->detailKesalahan,
                statusHttp: $pengecualian->statusHttp,
            );
        } catch (\Throwable $throwable) {
            $this->mengelolaIdempotencyKey->batalkanPenerimaan($idempotencyKey);

            throw $throwable;
        }

        $data = (new PemeriksaanHarianResource($hasilPenyimpanan))->resolve($permintaan);
        $respon = ResponApi::berhasil(
            pesan: 'Pemeriksaan harian QControl berhasil diterima',
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
