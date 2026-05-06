<?php

declare(strict_types=1);

namespace App\Http\Requests\Api\V1\QControl;

use App\Models\User;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Validator;

final class SimpanPemeriksaanHarianRequest extends FormRequest
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
            'clientDraftId' => ['nullable', 'string', 'max:100'],
            'tanggalProduksi' => ['required', 'date'],
            'lineProduksiId' => ['required', 'exists:qcontrol_line_produksi,id'],
            'nomorDokumen' => ['nullable', 'string', 'max:50'],
            'revisi' => ['nullable', 'string', 'max:20'],
            'catatan' => ['nullable', 'string'],
            'daftarPart' => ['required', 'array', 'min:1'],
            'daftarPart.*.partId' => ['required', 'exists:qcontrol_part,id'],
            'daftarPart.*.totalCheck' => ['required', 'integer', 'min:0'],
            'daftarPart.*.daftarDefect' => ['sometimes', 'array'],
            'daftarPart.*.daftarDefect.*.relasiPartDefectId' => ['required', 'exists:qcontrol_part_jenis_defect,id'],
            'daftarPart.*.daftarDefect.*.slotWaktuId' => ['required', 'exists:qcontrol_slot_waktu,id'],
            'daftarPart.*.daftarDefect.*.jumlahDefect' => ['required', 'integer', 'min:0'],
        ];
    }

    /**
     * @return array<int, \Closure(Validator): void>
     */
    public function after(): array
    {
        return [
            function (Validator $validator): void {
                if ($this->idempotencyKey() !== null) {
                    return;
                }

                $validator->errors()->add(
                    'X-Idempotency-Key',
                    'Header X-Idempotency-Key wajib diisi',
                );
            },
        ];
    }

    public function idempotencyKey(): ?string
    {
        $idempotencyKey = $this->header('X-Idempotency-Key');

        if (! is_string($idempotencyKey)) {
            return null;
        }

        $idempotencyKeyDipangkas = trim($idempotencyKey);

        return $idempotencyKeyDipangkas === '' ? null : $idempotencyKeyDipangkas;
    }

    /**
     * @return array<string, mixed>
     */
    public function payloadTervalidasi(): array
    {
        /** @var array<string, mixed> $payload */
        $payload = $this->validated();

        return $payload;
    }

    public function hashPayload(): string
    {
        $payloadJson = json_encode(
            $this->payloadTervalidasi(),
            JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES,
        );

        if (! is_string($payloadJson)) {
            $payloadJson = '{}';
        }

        return hash('sha256', $payloadJson);
    }

    public function penggunaHeadQC(): User
    {
        /** @var User $pengguna */
        $pengguna = $this->user();

        return $pengguna;
    }
}
