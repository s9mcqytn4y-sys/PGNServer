<?php

declare(strict_types=1);

use App\Domain\Kesehatan\StatusKoneksiDatabase;
use App\Infrastructure\Kesehatan\PemeriksaKoneksiDatabase;
use Mockery\MockInterface;

test('endpoint kesehatan mengembalikan status berhasil saat basis data terhubung', function () {
    $this->mock(PemeriksaKoneksiDatabase::class, function (MockInterface $mock): void {
        $mock->shouldReceive('periksa')
            ->once()
            ->andReturn(new StatusKoneksiDatabase(
                status: 'terhubung',
                driver: 'pgsql',
            ));
    });

    $this->getJson('/api/v1/kesehatan')
        ->assertSuccessful()
        ->assertJson([
            'berhasil' => true,
            'pesan' => 'Server berjalan normal',
            'metadata' => null,
            'kesalahan' => null,
        ])
        ->assertJsonPath('data.status', 'sehat')
        ->assertJsonPath('data.namaAplikasi', 'PGNServer')
        ->assertJsonPath('data.versiApi', 'v1')
        ->assertJsonPath('data.zonaWaktu', 'Asia/Jakarta')
        ->assertJsonPath('data.koneksiDatabase.status', 'terhubung')
        ->assertJsonPath('data.koneksiDatabase.driver', 'pgsql')
        ->assertJsonStructure([
            'berhasil',
            'pesan',
            'data' => [
                'status',
                'namaAplikasi',
                'versiApi',
                'waktuServer',
                'zonaWaktu',
                'koneksiDatabase' => [
                    'status',
                    'driver',
                ],
            ],
            'metadata',
            'kesalahan',
        ]);
});

test('endpoint kesehatan mengembalikan status terganggu saat basis data tidak terhubung', function () {
    $this->mock(PemeriksaKoneksiDatabase::class, function (MockInterface $mock): void {
        $mock->shouldReceive('periksa')
            ->once()
            ->andReturn(new StatusKoneksiDatabase(
                status: 'tidakTerhubung',
                driver: 'pgsql',
            ));
    });

    $this->getJson('/api/v1/kesehatan')
        ->assertStatus(503)
        ->assertJson([
            'berhasil' => false,
            'pesan' => 'Server berjalan, tetapi koneksi database belum tersedia',
            'metadata' => null,
            'kesalahan' => [
                'kode' => 'DATABASE_TIDAK_TERHUBUNG',
                'detail' => [],
            ],
        ])
        ->assertJsonPath('data.status', 'terganggu')
        ->assertJsonPath('data.namaAplikasi', 'PGNServer')
        ->assertJsonPath('data.versiApi', 'v1')
        ->assertJsonPath('data.zonaWaktu', 'Asia/Jakarta')
        ->assertJsonPath('data.koneksiDatabase.status', 'tidakTerhubung')
        ->assertJsonPath('data.koneksiDatabase.driver', 'pgsql');
});
