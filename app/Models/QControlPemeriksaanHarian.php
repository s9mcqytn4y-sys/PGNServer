<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;

final class QControlPemeriksaanHarian extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_pemeriksaan_harian';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'tanggal_produksi',
        'line_produksi_id',
        'kode_line_snapshot',
        'nama_line_snapshot',
        'nomor_dokumen',
        'revisi',
        'pengguna_headqc_id',
        'client_draft_id',
        'idempotency_key',
        'hash_payload',
        'status',
        'total_check',
        'total_ok',
        'total_defect',
        'rasio_defect',
        'catatan',
        'diterima_pada',
    ];

    /**
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'tanggal_produksi' => 'date',
            'total_check' => 'integer',
            'total_ok' => 'integer',
            'total_defect' => 'integer',
            'rasio_defect' => 'decimal:2',
            'diterima_pada' => 'datetime',
        ];
    }

    /**
     * @return BelongsTo<QControlLineProduksi, $this>
     */
    public function lineProduksi(): BelongsTo
    {
        return $this->belongsTo(QControlLineProduksi::class, 'line_produksi_id');
    }

    /**
     * @return BelongsTo<User, $this>
     */
    public function penggunaHeadQC(): BelongsTo
    {
        return $this->belongsTo(User::class, 'pengguna_headqc_id');
    }

    /**
     * @return HasMany<QControlPemeriksaanPart, $this>
     */
    public function daftarPemeriksaanPart(): HasMany
    {
        return $this->hasMany(QControlPemeriksaanPart::class, 'pemeriksaan_harian_id');
    }
}
