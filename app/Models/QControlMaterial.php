<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;

final class QControlMaterial extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_material';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'kode_material',
        'nama_material',
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
     * @return HasMany<QControlPart, $this>
     */
    public function daftarPart(): HasMany
    {
        return $this->hasMany(QControlPart::class, 'material_id');
    }
}
