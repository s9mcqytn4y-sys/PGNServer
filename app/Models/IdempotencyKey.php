<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

final class IdempotencyKey extends Model
{
    public const STATUS_DITERIMA = 'DITERIMA';

    public const STATUS_BERHASIL = 'BERHASIL';

    public const STATUS_GAGAL = 'GAGAL';

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'idempotency_keys';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'kunci_idempotency',
        'metode_http',
        'endpoint',
        'hash_payload',
        'status_pemrosesan',
        'response_status_http',
        'response_body',
        'sumber_aplikasi',
        'diproses_pada',
    ];

    /**
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'response_body' => 'array',
            'diproses_pada' => 'datetime',
        ];
    }
}
