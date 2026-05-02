<?php

declare(strict_types=1);

use App\Http\Controllers\Api\V1\PemeriksaanKesehatanController;
use Illuminate\Support\Facades\Route;

Route::prefix('v1')->group(function (): void {
    Route::get('/kesehatan', PemeriksaanKesehatanController::class)
        ->name('api.v1.kesehatan');
});
