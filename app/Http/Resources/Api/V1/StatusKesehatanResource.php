<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1;

use App\Domain\Kesehatan\StatusKesehatanServer;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

final class StatusKesehatanResource extends JsonResource
{
    /**
     * @return array{
     *     status: string,
     *     namaAplikasi: string,
     *     versiApi: string,
     *     waktuServer: string,
     *     zonaWaktu: string,
     *     koneksiDatabase: array{status: string, driver: string}
     * }
     */
    public function toArray(Request $request): array
    {
        /** @var StatusKesehatanServer $statusKesehatanServer */
        $statusKesehatanServer = $this->resource;

        return [
            'status' => $statusKesehatanServer->status,
            'namaAplikasi' => $statusKesehatanServer->namaAplikasi,
            'versiApi' => $statusKesehatanServer->versiApi,
            'waktuServer' => $statusKesehatanServer->waktuServer,
            'zonaWaktu' => $statusKesehatanServer->zonaWaktu,
            'koneksiDatabase' => $statusKesehatanServer->koneksiDatabase->keArray(),
        ];
    }
}
