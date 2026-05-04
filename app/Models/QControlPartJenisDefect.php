<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

final class QControlPartJenisDefect extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_part_jenis_defect';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'part_id',
        'jenis_defect_id',
        'kode_tampilan_defect',
        'urutan_tampil',
        'aktif',
    ];

    /**
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'urutan_tampil' => 'integer',
            'aktif' => 'boolean',
        ];
    }

    /**
     * @return BelongsTo<QControlPart, $this>
     */
    public function partTerkait(): BelongsTo
    {
        return $this->belongsTo(QControlPart::class, 'part_id');
    }

    /**
     * @return BelongsTo<QControlJenisDefect, $this>
     */
    public function jenisDefectTerkait(): BelongsTo
    {
        return $this->belongsTo(QControlJenisDefect::class, 'jenis_defect_id');
    }
}
