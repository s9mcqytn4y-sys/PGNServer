<?php

declare(strict_types=1);

namespace App\Http\Middleware;

use App\Models\User;
use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Exception\AccessDeniedHttpException;

/**
 * Memastikan hanya pengguna HeadQC yang dapat mengakses endpoint terproteksi QControl.
 */
final class PastikanPenggunaHeadQC
{
    /**
     * @param  Closure(Request): Response  $lanjut
     */
    public function handle(Request $permintaan, Closure $lanjut): Response
    {
        $pengguna = $permintaan->user();

        if (! $pengguna instanceof User) {
            throw new AccessDeniedHttpException('Hanya HeadQC yang diizinkan mengakses endpoint ini');
        }

        if ($pengguna->peran !== (string) config('qcontrol.headqc.peran')) {
            throw new AccessDeniedHttpException('Hanya HeadQC yang diizinkan mengakses endpoint ini');
        }

        return $lanjut($permintaan);
    }
}
