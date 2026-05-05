<?php

declare(strict_types=1);

namespace App\Application\QControl;

use App\Support\Errors\KodeKesalahanApi;
use RuntimeException;

final class PengecualianPemeriksaanHarian extends RuntimeException
{
    /**
     * @param  array<int, array{field: string|null, pesan: string}>  $detailKesalahan
     */
    public function __construct(
        string $pesan,
        public readonly KodeKesalahanApi $kodeKesalahan,
        public readonly array $detailKesalahan = [],
        public readonly int $statusHttp = 422,
    ) {
        parent::__construct($pesan);
    }
}
