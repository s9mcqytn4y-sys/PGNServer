<?php

declare(strict_types=1);

namespace App\Application\QControl;

use App\Models\QControlPemeriksaanHarian;
use Illuminate\Database\Eloquent\Collection;

final class MembacaDaftarPemeriksaanHarian
{
    /**
     * @param  array<string, mixed>  $filter
     * @return array{
     *     daftarPemeriksaanHarian: Collection<int, QControlPemeriksaanHarian>,
     *     metadata: array{jumlahData: int, limit: int}
     * }
     */
    public function jalankan(array $filter): array
    {
        $limit = isset($filter['limit']) ? (int) $filter['limit'] : 20;

        $query = QControlPemeriksaanHarian::query()
            ->with('lineProduksi')
            ->orderByDesc('tanggal_produksi')
            ->orderByDesc('diterima_pada')
            ->orderByDesc('dibuat_pada');

        if (isset($filter['tanggalProduksi'])) {
            $query->whereDate('tanggal_produksi', (string) $filter['tanggalProduksi']);
        }

        if (isset($filter['tanggalMulai'])) {
            $query->whereDate('tanggal_produksi', '>=', (string) $filter['tanggalMulai']);
        }

        if (isset($filter['tanggalSelesai'])) {
            $query->whereDate('tanggal_produksi', '<=', (string) $filter['tanggalSelesai']);
        }

        if (isset($filter['lineProduksiId'])) {
            $query->where('line_produksi_id', (string) $filter['lineProduksiId']);
        }

        $daftarPemeriksaanHarian = $query
            ->limit($limit)
            ->get();

        return [
            'daftarPemeriksaanHarian' => $daftarPemeriksaanHarian,
            'metadata' => [
                'jumlahData' => $daftarPemeriksaanHarian->count(),
                'limit' => $limit,
            ],
        ];
    }
}
