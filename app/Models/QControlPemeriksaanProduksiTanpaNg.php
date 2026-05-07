<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

final class QControlPemeriksaanProduksiTanpaNg extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_pemeriksaan_produksi_tanpa_ng';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'pemeriksaan_harian_id',
        'part_id',
        'uniq_no_part',
        'nomor_part_snapshot',
        'nama_part_snapshot',
        'total_produksi',
        'catatan',
    ];

    /**
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'total_produksi' => 'integer',
        ];
    }

    /**
     * @return BelongsTo<QControlPemeriksaanHarian, $this>
     */
    public function pemeriksaanHarian(): BelongsTo
    {
        return $this->belongsTo(QControlPemeriksaanHarian::class, 'pemeriksaan_harian_id');
    }

    /**
     * @return BelongsTo<QControlPart, $this>
     */
    public function part(): BelongsTo
    {
        return $this->belongsTo(QControlPart::class, 'part_id');
    }
}
