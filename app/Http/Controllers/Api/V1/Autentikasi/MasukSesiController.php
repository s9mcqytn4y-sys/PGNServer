<?php

declare(strict_types=1);

namespace App\Http\Controllers\Api\V1\Autentikasi;

use App\Http\Controllers\Controller;
use App\Http\Requests\Api\V1\Autentikasi\MasukSesiRequest;
use App\Http\Resources\Api\V1\Autentikasi\AutentikasiResource;
use App\Models\User;
use App\Support\Api\ResponApi;
use App\Support\Errors\KodeKesalahanApi;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Hash;

final class MasukSesiController extends Controller
{
    public function __invoke(MasukSesiRequest $permintaan): JsonResponse
    {
        $pengguna = User::query()
            ->where('email', $permintaan->emailMasuk())
            ->first();

        if ($pengguna === null || ! Hash::check($permintaan->kataSandi(), $pengguna->password)) {
            return ResponApi::gagal(
                pesan: 'Email atau password tidak sesuai',
                kodeKesalahan: KodeKesalahanApi::AUTENTIKASI_GAGAL,
                statusHttp: 401,
            );
        }

        if ($pengguna->peran !== (string) config('qcontrol.headqc.peran')) {
            return ResponApi::gagal(
                pesan: 'Hanya HeadQC yang diizinkan masuk ke sistem QControl',
                kodeKesalahan: KodeKesalahanApi::AKSES_DITOLAK,
                statusHttp: 403,
            );
        }

        $pengguna->tokens()->delete();

        $tokenBaru = $pengguna->createToken('qcontrol-desktop')->plainTextToken;
        $data = (new AutentikasiResource([
            'token' => $tokenBaru,
            'pengguna' => $pengguna,
        ]))->resolve($permintaan);

        return ResponApi::berhasil(
            pesan: 'Sesi berhasil dibuat',
            data: $data,
        );
    }
}
