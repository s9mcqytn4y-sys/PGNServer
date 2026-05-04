<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1\Autentikasi;

use App\Models\User;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

final class AutentikasiResource extends JsonResource
{
    /**
     * @return array{
     *     token: string,
     *     profil: array{namaPengguna: string, peran: string}
     * }
     */
    public function toArray(Request $request): array
    {
        /** @var array{token: string, pengguna: User} $data */
        $data = $this->resource;

        return [
            'token' => $data['token'],
            'profil' => [
                'namaPengguna' => $data['pengguna']->name,
                'peran' => (string) $data['pengguna']->peran,
            ],
        ];
    }
}
