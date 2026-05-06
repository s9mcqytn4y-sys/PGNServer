<?php

declare(strict_types=1);

namespace App\Http\Requests\Api\V1\QControl;

use Illuminate\Foundation\Http\FormRequest;

/**
 * Memvalidasi query read model bulanan recording defect QControl.
 */
final class LaporanBulananRecordingDefectRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    /**
     * @return array<string, array<int, string>>
     */
    public function rules(): array
    {
        return [
            'bulan' => ['required', 'integer', 'between:1,12'],
            'tahun' => ['required', 'integer', 'min:2000'],
            'lineProduksiId' => ['required', 'exists:qcontrol_line_produksi,id'],
            'partId' => ['nullable', 'exists:qcontrol_part,id'],
            'materialId' => ['nullable', 'exists:qcontrol_material,id'],
            'jenisDefectId' => ['nullable', 'exists:qcontrol_jenis_defect,id'],
        ];
    }

    /**
     * @return array<string, mixed>
     */
    public function filterTervalidasi(): array
    {
        /** @var array<string, mixed> $filter */
        $filter = $this->validated();

        return $filter;
    }
}
