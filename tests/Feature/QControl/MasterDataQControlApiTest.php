<?php

declare(strict_types=1);

use App\Models\User;
use Database\Seeders\DatabaseSeeder;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

uses(RefreshDatabase::class);

beforeEach(function (): void {
    $this->seed(DatabaseSeeder::class);
});

function tokenAutentikasiHeadQc(TestCase $pengujian): string
{
    $responsMasuk = $pengujian->postJson('/api/v1/login', [
        'email' => 'headqc@pgn.local',
        'password' => 'HeadQC@12345',
    ]);

    $responsMasuk
        ->assertSuccessful()
        ->assertJsonPath('data.profil.peran', 'HeadQC');

    return (string) $responsMasuk->json('data.token');
}

test('endpoint master data gagal tanpa token', function () {
    $this->getJson('/api/v1/qcontrol/master-data')
        ->assertStatus(401)
        ->assertJsonPath('kesalahan.kode', 'AUTENTIKASI_GAGAL');
});

test('HeadQC login lalu bisa mengambil master data', function () {
    $token = tokenAutentikasiHeadQc($this);

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/qcontrol/master-data')
        ->assertSuccessful()
        ->assertJson([
            'berhasil' => true,
            'pesan' => 'Master data QControl berhasil dimuat',
            'kesalahan' => null,
        ])
        ->assertJsonPath('data.versiMasterData', '2026.05.2D');
});

test('response master data memiliki line produksi', function () {
    $token = tokenAutentikasiHeadQc($this);

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/qcontrol/master-data')
        ->assertSuccessful()
        ->assertJsonCount(2, 'data.lineProduksi')
        ->assertJsonPath('data.lineProduksi.0.kodeLine', 'PRESS')
        ->assertJsonPath('data.lineProduksi.1.kodeLine', 'SEWING');
});

test('response master data memiliki part', function () {
    $token = tokenAutentikasiHeadQc($this);

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/qcontrol/master-data')
        ->assertSuccessful()
        ->assertJsonCount(10, 'data.part')
        ->assertJsonPath('metadata.jumlahPart', 10);
});

test('response master data memiliki jenis defect', function () {
    $token = tokenAutentikasiHeadQc($this);

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/qcontrol/master-data')
        ->assertSuccessful()
        ->assertJsonCount(23, 'data.jenisDefect')
        ->assertJsonPath('metadata.jumlahJenisDefect', 23);
});

test('response master data memiliki relasi part defect', function () {
    $token = tokenAutentikasiHeadQc($this);

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/qcontrol/master-data')
        ->assertSuccessful()
        ->assertJsonCount(80, 'data.relasiPartDefect')
        ->assertJsonPath('metadata.jumlahRelasiPartDefect', 80);
});

test('tidak ada peran selain HeadQC yang dibuat', function () {
    $daftarPeran = User::query()
        ->select('peran')
        ->distinct()
        ->pluck('peran')
        ->all();

    expect($daftarPeran)->toBe(['HeadQC']);
});

test('endpoint kesehatan tetap tersedia', function () {
    $this->getJson('/api/v1/kesehatan')
        ->assertOk();
});

test('endpoint qcontrol contoh tetap tersedia', function () {
    $this->postJson(
        '/api/v1/qcontrol/contoh',
        [
            'contoh' => true,
            'sumber' => 'uji-master-data',
        ],
        [
            'X-Idempotency-Key' => 'qcontrol:master-data:test-001',
        ],
    )->assertOk();
});
