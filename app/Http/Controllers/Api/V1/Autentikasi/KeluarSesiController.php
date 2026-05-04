<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1\Autentikasi;

use App\Http\Controllers\Controller;
use App\Support\Api\ResponApi;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

final class KeluarSesiController extends Controller
{
    public function __invoke(Request $permintaan): JsonResponse
    {
        $permintaan->user()?->currentAccessToken()?->delete();
        Auth::guard('web')->logout();
        Auth::forgetGuards();

        return ResponApi::berhasil(
            pesan: 'Sesi berhasil ditutup',
        );
    }
}
