<?php

declare(strict_types=1);

use App\Http\Controllers\Api\V1\Autentikasi\KeluarSesiController;
use App\Http\Controllers\Api\V1\Autentikasi\MasukSesiController;
use App\Http\Controllers\Api\V1\Autentikasi\ProfilSayaController;
use App\Http\Controllers\Api\V1\PemeriksaanKesehatanController;
use App\Http\Controllers\Api\V1\QControl\DaftarPemeriksaanHarianController;
use App\Http\Controllers\Api\V1\QControl\DetailPemeriksaanHarianController;
use App\Http\Controllers\Api\V1\QControl\LaporanBulananRecordingDefectController;
use App\Http\Controllers\Api\V1\QControl\MembacaMasterDataQControlController;
use App\Http\Controllers\Api\V1\QControl\PenerimaanContohSinkronisasiController;
use App\Http\Controllers\Api\V1\QControl\SimpanPemeriksaanHarianController;
use Illuminate\Support\Facades\Route;

Route::prefix('v1')->group(function (): void {
    Route::get('/kesehatan', PemeriksaanKesehatanController::class)
        ->name('api.v1.kesehatan');

    Route::middleware(['auth:sanctum', 'headqc'])->group(function (): void {
        Route::prefix('qcontrol')->name('api.v1.qcontrol.')->group(function (): void {
            Route::post('/contoh', PenerimaanContohSinkronisasiController::class)
                ->name('contoh-sinkronisasi');

            Route::get('/master-data', MembacaMasterDataQControlController::class)
                ->name('master-data');

            Route::get('/pemeriksaan-harian', DaftarPemeriksaanHarianController::class)
                ->name('daftar-pemeriksaan-harian');

            Route::get('/pemeriksaan-harian/{pemeriksaanHarian}', DetailPemeriksaanHarianController::class)
                ->name('detail-pemeriksaan-harian');

            Route::post('/pemeriksaan-harian', SimpanPemeriksaanHarianController::class)
                ->name('pemeriksaan-harian');

            Route::get('/laporan-bulanan/recording-defect', LaporanBulananRecordingDefectController::class)
                ->name('laporan-bulanan.recording-defect');
        });
    });

    Route::post('/login', MasukSesiController::class)
        ->name('api.v1.login');

    Route::middleware(['auth:sanctum', 'headqc'])->group(function (): void {
        Route::post('/logout', KeluarSesiController::class)
            ->name('api.v1.logout');

        Route::get('/profil-saya', ProfilSayaController::class)
            ->name('api.v1.profil-saya');
    });
});
