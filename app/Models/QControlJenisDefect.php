<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;

final class QControlJenisDefect extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_jenis_defect';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'kode_defect',
        'nama_defect',
        'kategori_defect_id',
        'aktif',
    ];

    /**
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'aktif' => 'boolean',
        ];
    }

    /**
     * @return BelongsTo<QControlKategoriDefect, $this>
     */
    public function kategoriDefectTerkait(): BelongsTo
    {
        return $this->belongsTo(QControlKategoriDefect::class, 'kategori_defect_id');
    }

    /**
     * @return HasMany<QControlPartJenisDefect, $this>
     */
    public function daftarRelasiPart(): HasMany
    {
        return $this->hasMany(QControlPartJenisDefect::class, 'jenis_defect_id');
    }

    /**
     * @return HasMany<QControlPemeriksaanDefectSlot, $this>
     */
    public function daftarPemeriksaanDefectSlot(): HasMany
    {
        return $this->hasMany(QControlPemeriksaanDefectSlot::class, 'jenis_defect_id');
    }
}
