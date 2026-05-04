<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;

final class QControlKategoriDefect extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_kategori_defect';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'kode_kategori',
        'nama_kategori',
        'aktif',
        'urutan_tampil',
    ];

    /**
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'aktif' => 'boolean',
            'urutan_tampil' => 'integer',
        ];
    }

    /**
     * @return HasMany<QControlJenisDefect, $this>
     */
    public function daftarJenisDefect(): HasMany
    {
        return $this->hasMany(QControlJenisDefect::class, 'kategori_defect_id');
    }
}
