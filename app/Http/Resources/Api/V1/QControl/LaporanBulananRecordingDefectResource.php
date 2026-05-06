<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1\QControl;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/**
 * Membentuk payload final read model bulanan yang stabil untuk client QControl.
 */
final class LaporanBulananRecordingDefectResource extends JsonResource
{
    /**
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        /** @var array<string, mixed> $data */
        $data = $this->resource;

        return $data;
    }
}
