<?php

declare(strict_types=1);

namespace App\Http\Requests\Api\V1\Autentikasi;

use Illuminate\Foundation\Http\FormRequest;

final class MasukSesiRequest extends FormRequest
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
            'email' => ['required', 'email'],
            'password' => ['required', 'string'],
        ];
    }

    public function emailMasuk(): string
    {
        return (string) $this->string('email');
    }

    public function kataSandi(): string
    {
        return (string) $this->string('password');
    }
}
