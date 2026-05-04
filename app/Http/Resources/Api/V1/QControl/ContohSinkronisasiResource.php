<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1\QControl;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

final class ContohSinkronisasiResource extends JsonResource
{
    /**
     * @return array{
     *     diterima: bool,
     *     duplikat: bool,
     *     idempotencyKey: string,
     *     endpoint: string,
     *     mode: string
     * }
     */
    public function toArray(Request $request): array
    {
        /** @var array{
         *     diterima: bool,
         *     duplikat: bool,
         *     idempotencyKey: string,
         *     endpoint: string,
         *     mode: string
         * } $data
         */
        $data = $this->resource;

        return [
            'diterima' => $data['diterima'],
            'duplikat' => $data['duplikat'],
            'idempotencyKey' => $data['idempotencyKey'],
            'endpoint' => $data['endpoint'],
            'mode' => $data['mode'],
        ];
    }
}
