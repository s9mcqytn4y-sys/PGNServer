<?php

declare(strict_types=1);

use App\Support\Api\ResponApi;
use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return ResponApi::berhasil(
        pesan: 'PGNServer REST API siap digunakan',
        data: [
            'namaAplikasi' => config('app.name'),
            'versiApi' => 'v1',
        ],
    );
});
