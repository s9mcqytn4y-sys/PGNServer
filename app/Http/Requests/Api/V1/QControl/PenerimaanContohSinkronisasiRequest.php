<?php

declare(strict_types=1);

namespace App\Http\Requests\Api\V1\QControl;

use Illuminate\Foundation\Http\FormRequest;

final class PenerimaanContohSinkronisasiRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    /**
     * @return array<string, array<int, string>|string>
     */
    public function rules(): array
    {
        return [];
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
    public function payloadDiterima(): array
    {
        $payloadDiterima = $this->json()->all();

        return is_array($payloadDiterima) ? $payloadDiterima : [];
    }

    public function hashPayload(): string
    {
        $payloadJson = json_encode(
            $this->payloadDiterima(),
            JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES,
        );

        if (! is_string($payloadJson)) {
            $payloadJson = '[]';
        }

        return hash('sha256', $payloadJson);
    }
}
