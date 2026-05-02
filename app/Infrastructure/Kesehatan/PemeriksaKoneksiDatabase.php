<?php

declare(strict_types=1);

namespace App\Infrastructure\Kesehatan;

use App\Domain\Kesehatan\StatusKoneksiDatabase;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Throwable;

class PemeriksaKoneksiDatabase
{
    public function periksa(): StatusKoneksiDatabase
    {
        $driver = (string) config('database.default', 'pgsql');

        try {
            $koneksi = DB::connection();
            $koneksi->select('select 1');

            return new StatusKoneksiDatabase(
                status: 'terhubung',
                driver: (string) $koneksi->getDriverName(),
            );
        } catch (Throwable $throwable) {
            Log::warning('Pemeriksaan kesehatan gagal menghubungkan basis data.', [
                'pesan' => $throwable->getMessage(),
                'driver' => $driver,
            ]);

            return new StatusKoneksiDatabase(
                status: 'tidakTerhubung',
                driver: $driver,
            );
        }
    }
}
