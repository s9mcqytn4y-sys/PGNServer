<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;

final class QControlPemeriksaanPart extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_pemeriksaan_part';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'pemeriksaan_harian_id',
        'part_id',
        'total_check',
        'total_ok',
        'total_defect',
        'rasio_defect',
        'urutan_tampil',
    ];

    /**
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'total_check' => 'integer',
            'total_ok' => 'integer',
            'total_defect' => 'integer',
            'rasio_defect' => 'decimal:2',
            'urutan_tampil' => 'integer',
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
    public function partTerkait(): BelongsTo
    {
        return $this->belongsTo(QControlPart::class, 'part_id');
    }

    /**
     * @return HasMany<QControlPemeriksaanDefectSlot, $this>
     */
    public function daftarDefectSlot(): HasMany
    {
        return $this->hasMany(QControlPemeriksaanDefectSlot::class, 'pemeriksaan_part_id');
    }
}
