<?php

declare(strict_types=1);

namespace App\Application\Idempotency;

use App\Domain\Idempotency\HasilPemeriksaanIdempotency;
use App\Models\IdempotencyKey;
use Illuminate\Database\QueryException;
use Illuminate\Support\Carbon;

final class MengelolaIdempotencyKey
{
    public function periksaAtauSimpanTransaksiSerius(
        string $kunciIdempotency,
        string $metodeHttp,
        string $endpoint,
        string $hashPayload,
        ?string $sumberAplikasi = null,
    ): HasilPemeriksaanIdempotency {
        $idempotencyKeyTersimpan = IdempotencyKey::query()
            ->where('kunci_idempotency', $kunciIdempotency)
            ->first();

        if ($idempotencyKeyTersimpan !== null) {
            return $this->buatHasilPemeriksaanTransaksiSerius(
                $idempotencyKeyTersimpan,
                $hashPayload,
            );
        }

        try {
            IdempotencyKey::query()->create([
                'kunci_idempotency' => $kunciIdempotency,
                'metode_http' => $metodeHttp,
                'endpoint' => $endpoint,
                'hash_payload' => $hashPayload,
                'status_pemrosesan' => IdempotencyKey::STATUS_DITERIMA,
                'sumber_aplikasi' => $sumberAplikasi,
            ]);

            return new HasilPemeriksaanIdempotency(sudahAda: false);
        } catch (QueryException) {
            $idempotencyKeyTersimpan = IdempotencyKey::query()
                ->where('kunci_idempotency', $kunciIdempotency)
                ->first();

            if ($idempotencyKeyTersimpan === null) {
                return new HasilPemeriksaanIdempotency(sudahAda: false);
            }

            return $this->buatHasilPemeriksaanTransaksiSerius(
                $idempotencyKeyTersimpan,
                $hashPayload,
            );
        }
    }

    public function periksaAtauSimpanPenerimaan(
        string $kunciIdempotency,
        string $metodeHttp,
        string $endpoint,
        string $hashPayload,
        ?string $sumberAplikasi = null,
    ): HasilPemeriksaanIdempotency {
        $idempotencyKeyTersimpan = IdempotencyKey::query()
            ->where('kunci_idempotency', $kunciIdempotency)
            ->first();

        if ($idempotencyKeyTersimpan !== null) {
            return new HasilPemeriksaanIdempotency(
                sudahAda: true,
                dataTersimpan: $idempotencyKeyTersimpan->response_body,
                statusHttpTersimpan: $idempotencyKeyTersimpan->response_status_http,
            );
        }

        try {
            IdempotencyKey::query()->create([
                'kunci_idempotency' => $kunciIdempotency,
                'metode_http' => $metodeHttp,
                'endpoint' => $endpoint,
                'hash_payload' => $hashPayload,
                'status_pemrosesan' => IdempotencyKey::STATUS_DITERIMA,
                'sumber_aplikasi' => $sumberAplikasi,
            ]);

            return new HasilPemeriksaanIdempotency(sudahAda: false);
        } catch (QueryException) {
            $idempotencyKeyTersimpan = IdempotencyKey::query()
                ->where('kunci_idempotency', $kunciIdempotency)
                ->first();

            return new HasilPemeriksaanIdempotency(
                sudahAda: true,
                dataTersimpan: $idempotencyKeyTersimpan?->response_body,
                statusHttpTersimpan: $idempotencyKeyTersimpan?->response_status_http,
            );
        }
    }

    /**
     * @param  array<string, mixed>  $responseBody
     */
    public function simpanResponseBerhasil(
        string $kunciIdempotency,
        int $statusHttp,
        array $responseBody,
    ): void {
        IdempotencyKey::query()
            ->where('kunci_idempotency', $kunciIdempotency)
            ->update([
                'status_pemrosesan' => IdempotencyKey::STATUS_BERHASIL,
                'response_status_http' => $statusHttp,
                'response_body' => $responseBody,
                'diproses_pada' => Carbon::now(),
            ]);
    }

    public function batalkanPenerimaan(string $kunciIdempotency): void
    {
        IdempotencyKey::query()
            ->where('kunci_idempotency', $kunciIdempotency)
            ->where('status_pemrosesan', IdempotencyKey::STATUS_DITERIMA)
            ->delete();
    }

    private function buatHasilPemeriksaanTransaksiSerius(
        IdempotencyKey $idempotencyKeyTersimpan,
        string $hashPayload,
    ): HasilPemeriksaanIdempotency {
        if (
            $idempotencyKeyTersimpan->hash_payload !== null
            && $idempotencyKeyTersimpan->hash_payload !== $hashPayload
        ) {
            return new HasilPemeriksaanIdempotency(
                sudahAda: true,
                konflikPayload: true,
            );
        }

        if (
            is_array($idempotencyKeyTersimpan->response_body)
            && is_int($idempotencyKeyTersimpan->response_status_http)
        ) {
            return new HasilPemeriksaanIdempotency(
                sudahAda: true,
                dataTersimpan: $idempotencyKeyTersimpan->response_body,
                statusHttpTersimpan: $idempotencyKeyTersimpan->response_status_http,
            );
        }

        return new HasilPemeriksaanIdempotency(
            sudahAda: false,
            perluDiprosesUlang: true,
        );
    }
}
