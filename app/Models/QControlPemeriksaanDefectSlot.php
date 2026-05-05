<?php

declare(strict_types=1);

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

final class QControlPemeriksaanDefectSlot extends Model
{
    use HasUuids;

    public const CREATED_AT = 'dibuat_pada';

    public const UPDATED_AT = 'diperbarui_pada';

    protected $table = 'qcontrol_pemeriksaan_defect_slot';

    public $incrementing = false;

    protected $keyType = 'string';

    /**
     * @var list<string>
     */
    protected $fillable = [
        'pemeriksaan_part_id',
        'relasi_part_defect_id',
        'jenis_defect_id',
        'slot_waktu_id',
        'jumlah_defect',
    ];

    /**
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'jumlah_defect' => 'integer',
        ];
    }

    /**
     * @return BelongsTo<QControlPemeriksaanPart, $this>
     */
    public function pemeriksaanPart(): BelongsTo
    {
        return $this->belongsTo(QControlPemeriksaanPart::class, 'pemeriksaan_part_id');
    }

    /**
     * @return BelongsTo<QControlPartJenisDefect, $this>
     */
    public function relasiPartDefect(): BelongsTo
    {
        return $this->belongsTo(QControlPartJenisDefect::class, 'relasi_part_defect_id');
    }

    /**
     * @return BelongsTo<QControlJenisDefect, $this>
     */
    public function jenisDefect(): BelongsTo
    {
        return $this->belongsTo(QControlJenisDefect::class, 'jenis_defect_id');
    }

    /**
     * @return BelongsTo<QControlSlotWaktu, $this>
     */
    public function slotWaktu(): BelongsTo
    {
        return $this->belongsTo(QControlSlotWaktu::class, 'slot_waktu_id');
    }
}
