<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;

final class QControlLineProduksi extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_line_produksi';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'kode_line',
        'nama_line',
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
     * @return HasMany<QControlPart, $this>
     */
    public function daftarPartDefault(): HasMany
    {
        return $this->hasMany(QControlPart::class, 'line_default_id');
    }
}
