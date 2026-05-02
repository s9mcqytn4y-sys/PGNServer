<?php

declare(strict_types=1);

namespace App\Domain\Kesehatan;

final readonly class StatusKesehatanServer
{
    public function __construct(
        public string $status,
        public string $namaAplikasi,
        public string $versiApi,
        public string $waktuServer,
        public string $zonaWaktu,
        public StatusKoneksiDatabase $koneksiDatabase,
    ) {
    }

    public function databaseTersedia(): bool
    {
        return $this->koneksiDatabase->status === 'terhubung';
    }
}
