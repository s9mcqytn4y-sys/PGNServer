<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1\Autentikasi;

use App\Models\User;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

final class ProfilPenggunaResource extends JsonResource
{
    /**
     * @return array{namaPengguna: string, email: string, peran: string}
     */
    public function toArray(Request $request): array
    {
        /** @var User $pengguna */
        $pengguna = $this->resource;

        return [
            'namaPengguna' => $pengguna->name,
            'email' => $pengguna->email,
            'peran' => (string) $pengguna->peran,
        ];
    }
}
