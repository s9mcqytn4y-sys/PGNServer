<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;

final class QControlSlotWaktu extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_slot_waktu';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'kode_slot',
        'label_slot',
        'jam_mulai',
        'jam_selesai',
        'urutan_tampil',
        'aktif',
    ];

    /**
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'jam_mulai' => 'datetime:H:i:s',
            'jam_selesai' => 'datetime:H:i:s',
            'urutan_tampil' => 'integer',
            'aktif' => 'boolean',
        ];
    }
}
