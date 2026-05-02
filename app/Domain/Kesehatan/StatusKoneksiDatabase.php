<?php

declare(strict_types=1);

namespace App\Domain\Kesehatan;

final readonly class StatusKoneksiDatabase
{
    public function __construct(
        public string $status,
        public string $driver,
    ) {
    }

    /**
     * @return array{status: string, driver: string}
     */
    public function keArray(): array
    {
        return [
            'status' => $this->status,
            'driver' => $this->driver,
        ];
    }
}
