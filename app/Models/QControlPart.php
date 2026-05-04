<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;

final class QControlPart extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_part';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'kode_unik_part',
        'nama_part',
        'nomor_part',
        'material_id',
        'kode_proyek',
        'jumlah_item_per_kanban',
        'line_default_id',
        'aktif',
        'sumber_data',
    ];

    /**
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'jumlah_item_per_kanban' => 'integer',
            'aktif' => 'boolean',
        ];
    }

    /**
     * @return BelongsTo<QControlMaterial, $this>
     */
    public function materialTerkait(): BelongsTo
    {
        return $this->belongsTo(QControlMaterial::class, 'material_id');
    }

    /**
     * @return BelongsTo<QControlLineProduksi, $this>
     */
    public function lineProduksiDefault(): BelongsTo
    {
        return $this->belongsTo(QControlLineProduksi::class, 'line_default_id');
    }

    /**
     * @return HasMany<QControlPartJenisDefect, $this>
     */
    public function daftarRelasiJenisDefect(): HasMany
    {
        return $this->hasMany(QControlPartJenisDefect::class, 'part_id');
    }
}
