<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1\QControl;

use App\Application\QControl\MembacaLaporanBulananRecordingDefect;
use App\Application\QControl\PengecualianPemeriksaanHarian;
use App\Http\Controllers\Controller;
use App\Http\Requests\Api\V1\QControl\LaporanBulananRecordingDefectRequest;
use App\Http\Resources\Api\V1\QControl\LaporanBulananRecordingDefectResource;
use App\Support\Api\ResponApi;
use Illuminate\Http\JsonResponse;

/**
 * Menyajikan read model bulanan recording defect hasil agregasi transaksi daily.
 */
final class LaporanBulananRecordingDefectController extends Controller
{
    public function __construct(
        private readonly MembacaLaporanBulananRecordingDefect $membacaLaporanBulananRecordingDefect,
    ) {}

    public function __invoke(LaporanBulananRecordingDefectRequest $permintaan): JsonResponse
    {
        try {
            $hasilPembacaan = $this->membacaLaporanBulananRecordingDefect->jalankan(
                $permintaan->filterTervalidasi(),
            );
        } catch (PengecualianPemeriksaanHarian $pengecualian) {
            return ResponApi::gagal(
                pesan: $pengecualian->getMessage(),
                kodeKesalahan: $pengecualian->kodeKesalahan,
                detailKesalahan: $pengecualian->detailKesalahan,
                statusHttp: $pengecualian->statusHttp,
            );
        }

        return ResponApi::berhasil(
            pesan: 'Laporan bulanan recording defect QControl berhasil dimuat',
            data: (new LaporanBulananRecordingDefectResource($hasilPembacaan))->resolve($permintaan),
        );
    }
}
