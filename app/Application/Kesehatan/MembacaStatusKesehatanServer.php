<?php

declare(strict_types=1);

namespace App\Application\Kesehatan;

use App\Domain\Kesehatan\StatusKesehatanServer;
use App\Infrastructure\Kesehatan\PemeriksaKoneksiDatabase;

final class MembacaStatusKesehatanServer
{
    public function __construct(
        private PemeriksaKoneksiDatabase $pemeriksaKoneksiDatabase,
    ) {
    }

    public function jalankan(): StatusKesehatanServer
    {
        $statusKoneksiDatabase = $this->pemeriksaKoneksiDatabase->periksa();
        $zonaWaktu = (string) config('app.timezone', 'Asia/Jakarta');

        return new StatusKesehatanServer(
            status: $statusKoneksiDatabase->status === 'terhubung' ? 'sehat' : 'terganggu',
            namaAplikasi: (string) config('app.name', 'PGNServer'),
            versiApi: 'v1',
            waktuServer: now()->setTimezone($zonaWaktu)->toIso8601String(),
            zonaWaktu: $zonaWaktu,
            koneksiDatabase: $statusKoneksiDatabase,
        );
    }
}
