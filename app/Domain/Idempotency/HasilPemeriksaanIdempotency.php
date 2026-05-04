<?php

declare(strict_types=1);

namespace App\Domain\Idempotency;

final class HasilPemeriksaanIdempotency
{
    /**
     * @param  array<string, mixed>|null  $dataTersimpan
     */
    public function __construct(
        public bool $sudahAda,
        public ?array $dataTersimpan = null,
        public ?int $statusHttpTersimpan = null,
    ) {}
}
