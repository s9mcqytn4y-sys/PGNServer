<?php

declare(strict_types=1);

use App\Models\User;
use Database\Seeders\HeadQCPenggunaSeeder;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Laravel\Sanctum\PersonalAccessToken;

uses(RefreshDatabase::class);

beforeEach(function (): void {
    $this->seed(HeadQCPenggunaSeeder::class);
});

test('HeadQC bisa login dengan email dan password yang benar', function () {
    $this->postJson('/api/v1/login', [
        'email' => 'headqc@pgn.local',
        'password' => 'HeadQC@12345',
    ])
        ->assertSuccessful()
        ->assertJson([
            'berhasil' => true,
            'pesan' => 'Sesi berhasil dibuat',
            'metadata' => null,
            'kesalahan' => null,
        ])
        ->assertJsonPath('data.profil.namaPengguna', 'HeadQC')
        ->assertJsonPath('data.profil.peran', 'HeadQC');
});

test('response login HeadQC memiliki token', function () {
    $respons = $this->postJson('/api/v1/login', [
        'email' => 'headqc@pgn.local',
        'password' => 'HeadQC@12345',
    ]);

    $respons
        ->assertSuccessful()
        ->assertJsonPath('data.profil.peran', 'HeadQC');

    expect($respons->json('data.token'))
        ->toBeString()
        ->not->toBe('');
});

test('login gagal dengan password yang salah', function () {
    $this->postJson('/api/v1/login', [
        'email' => 'headqc@pgn.local',
        'password' => 'salah-total',
    ])
        ->assertStatus(401)
        ->assertJson([
            'berhasil' => false,
            'pesan' => 'Email atau password tidak sesuai',
            'metadata' => null,
            'kesalahan' => [
                'kode' => 'AUTENTIKASI_GAGAL',
                'detail' => [],
            ],
        ]);
});

test('profil saya gagal tanpa token', function () {
    $this->getJson('/api/v1/profil-saya')
        ->assertStatus(401)
        ->assertJson([
            'berhasil' => false,
            'pesan' => 'Autentikasi diperlukan untuk mengakses endpoint ini',
            'metadata' => null,
            'kesalahan' => [
                'kode' => 'AUTENTIKASI_GAGAL',
                'detail' => [],
            ],
        ]);
});

test('profil saya berhasil dengan token', function () {
    $pengguna = User::query()->where('email', 'headqc@pgn.local')->firstOrFail();
    $token = $pengguna->createToken('qcontrol-desktop')->plainTextToken;

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/profil-saya')
        ->assertSuccessful()
        ->assertJson([
            'berhasil' => true,
            'pesan' => 'Profil pengguna berhasil dimuat',
            'metadata' => null,
            'kesalahan' => null,
        ])
        ->assertJsonPath('data.namaPengguna', 'HeadQC')
        ->assertJsonPath('data.email', 'headqc@pgn.local')
        ->assertJsonPath('data.peran', 'HeadQC');
});

test('logout berhasil dengan token', function () {
    $pengguna = User::query()->where('email', 'headqc@pgn.local')->firstOrFail();
    $token = $pengguna->createToken('qcontrol-desktop')->plainTextToken;

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->postJson('/api/v1/logout')
        ->assertSuccessful()
        ->assertJson([
            'berhasil' => true,
            'pesan' => 'Sesi berhasil ditutup',
            'metadata' => null,
            'kesalahan' => null,
        ]);
});

test('setelah logout token tidak bisa dipakai lagi untuk profil saya', function () {
    $pengguna = User::query()->where('email', 'headqc@pgn.local')->firstOrFail();
    $token = $pengguna->createToken('qcontrol-desktop')->plainTextToken;

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->postJson('/api/v1/logout')
        ->assertSuccessful();

    expect(PersonalAccessToken::query()->count())->toBe(0);

    $this->withHeader('Authorization', 'Bearer '.$token)
        ->getJson('/api/v1/profil-saya')
        ->assertStatus(401)
        ->assertJsonPath('kesalahan.kode', 'AUTENTIKASI_GAGAL');
});
