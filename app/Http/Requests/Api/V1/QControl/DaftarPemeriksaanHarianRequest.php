<?php

declare(strict_types=1);

namespace App\Http\Requests\Api\V1\QControl;

use Illuminate\Foundation\Http\FormRequest;

final class DaftarPemeriksaanHarianRequest extends FormRequest
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
            'tanggalProduksi' => ['nullable', 'date'],
            'tanggalMulai' => ['nullable', 'date'],
            'tanggalSelesai' => ['nullable', 'date'],
            'lineProduksiId' => ['nullable', 'exists:qcontrol_line_produksi,id'],
            'limit' => ['nullable', 'integer', 'min:1', 'max:100'],
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
