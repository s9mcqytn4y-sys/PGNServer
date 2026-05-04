<?php

declare(strict_types=1);

use App\Http\Controllers\Api\V1\PemeriksaanKesehatanController;
use App\Http\Controllers\Api\V1\QControl\PenerimaanContohSinkronisasiController;
use Illuminate\Support\Facades\Route;

Route::prefix('v1')->group(function (): void {
    Route::get('/kesehatan', PemeriksaanKesehatanController::class)
        ->name('api.v1.kesehatan');

    Route::prefix('qcontrol')->name('api.v1.qcontrol.')->group(function (): void {
        Route::post('/contoh', PenerimaanContohSinkronisasiController::class)
            ->name('contoh-sinkronisasi');
    });
});
