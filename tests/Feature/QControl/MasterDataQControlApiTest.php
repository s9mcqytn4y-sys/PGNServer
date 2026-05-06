<?php

declare(strict_types=1);

use App\Models\User;
use Database\Seeders\DatabaseSeeder;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Hash;
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
        ->assertJsonPath('data.versiMasterData', '2026.05.2F-A');
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
        ->assertJsonCount(18, 'data.part')
        ->assertJsonPath('metadata.jumlahPart', 18);
});

test('response master data memiliki jenis defect', function () {
    $token = tokenAutentikasiHeadQc($this);

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/qcontrol/master-data')
        ->assertSuccessful()
        ->assertJsonCount(34, 'data.jenisDefect')
        ->assertJsonPath('metadata.jumlahJenisDefect', 34);
});

test('response master data memiliki relasi part defect', function () {
    $token = tokenAutentikasiHeadQc($this);

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/qcontrol/master-data')
        ->assertSuccessful()
        ->assertJsonCount(181, 'data.relasiPartDefect')
        ->assertJsonPath('metadata.jumlahRelasiPartDefect', 181);
});

test('master data test kriteria lengkap sesuai dokumen target', function () {
    $token = tokenAutentikasiHeadQc($this);

    $respon = $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/qcontrol/master-data');

    $respon->assertSuccessful();

    // 3. response memiliki line PRESS dan SEWING.
    $respon->assertJsonFragment(['kodeLine' => 'PRESS']);
    $respon->assertJsonFragment(['kodeLine' => 'SEWING']);

    // 4. response memiliki part SEWING Felt Seat Back.
    // 5. response memiliki part SEWING Cover FR Seat Hinge.
    // 6. response memiliki part Protector RR Seat Back RH/LH.
    $respon->assertJsonFragment(['kodeUnikPart' => 'FSB', 'namaPart' => 'Felt Seat Back']);
    $respon->assertJsonFragment(['kodeUnikPart' => 'CFRSH', 'namaPart' => 'Cover FR Seat Hinge']);
    $respon->assertJsonFragment(['kodeUnikPart' => 'PRSB_RH_070', 'namaPart' => 'Protector RR Seat Back RH']);
    $respon->assertJsonFragment(['kodeUnikPart' => 'PRSB_LH_080', 'namaPart' => 'Protector RR Seat Back LH']);

    // 7. response memiliki jenis defect SPUNBOND_TIDAK_MEREKAT.
    // 8. response memiliki jenis defect BACKSTITCH_KURANG_DARI_15MM.
    $respon->assertJsonFragment(['kodeDefect' => 'SPUNBOND_TIDAK_MEREKAT']);
    $respon->assertJsonFragment(['kodeDefect' => 'BACKSTITCH_KURANG_DARI_15MM']);

    // 9. response relasiPartDefect memiliki field kodeTampilanDefect.
    $relasiFSB_A = collect($respon->json('data.relasiPartDefect'))
        ->where('kodeUnikPart', 'FSB')
        ->where('kodeTampilanDefect', 'A')
        ->first();
    expect($relasiFSB_A)->not->toBeNull();
    expect($relasiFSB_A['kodeTampilanDefect'])->toBe('A');
    expect($relasiFSB_A['kodeDefect'])->toBe('SOBEK');

    // 10. relasi Felt Seat Back memiliki kode tampilan A sampai L.
    $relasiFSB = collect($respon->json('data.relasiPartDefect'))
        ->where('kodeUnikPart', 'FSB');
    expect($relasiFSB->count())->toBe(12);
    expect($relasiFSB->pluck('kodeTampilanDefect')->sort()->values()->all())->toBe(['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L']);

    // 11. relasi Protector memiliki kode tampilan A sampai M.
    $relasiPRSB = collect($respon->json('data.relasiPartDefect'))
        ->where('kodeUnikPart', 'PRSB_RH_070');
    expect($relasiPRSB->count())->toBe(13);
    expect($relasiPRSB->pluck('kodeTampilanDefect')->sort()->values()->all())->toBe(['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M']);
});

test('pengguna non HeadQC ditolak saat mengambil master data', function () {
    $pengguna = User::factory()->create([
        'email' => 'operator-master@pgn.local',
        'password' => Hash::make('password'),
        'peran' => 'Operator',
    ]);

    $token = $pengguna->createToken('qcontrol-desktop')->plainTextToken;

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/qcontrol/master-data')
        ->assertStatus(403)
        ->assertJsonPath('kesalahan.kode', 'AKSES_DITOLAK');
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
    $token = tokenAutentikasiHeadQc($this);

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->postJson(
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
