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
    $emailHeadQC = (string) config('qcontrol.headqc.email');
    $passwordHeadQC = (string) config('qcontrol.headqc.password_default');
    $namaPenggunaHeadQC = (string) config('qcontrol.headqc.nama_pengguna');
    $peranHeadQC = (string) config('qcontrol.headqc.peran');

    $this->postJson('/api/v1/login', [
        'email' => $emailHeadQC,
        'password' => $passwordHeadQC,
    ])
        ->assertSuccessful()
        ->assertJson([
            'berhasil' => true,
            'pesan' => 'Sesi berhasil dibuat',
            'metadata' => null,
            'kesalahan' => null,
        ])
        ->assertJsonPath('data.profil.namaPengguna', $namaPenggunaHeadQC)
        ->assertJsonPath('data.profil.peran', $peranHeadQC);
});

test('response login HeadQC memiliki token', function () {
    $emailHeadQC = (string) config('qcontrol.headqc.email');
    $passwordHeadQC = (string) config('qcontrol.headqc.password_default');

    $respons = $this->postJson('/api/v1/login', [
        'email' => $emailHeadQC,
        'password' => $passwordHeadQC,
    ]);

    $respons
        ->assertSuccessful()
        ->assertJsonPath('data.profil.peran', config('qcontrol.headqc.peran'));

    expect($respons->json('data.token'))
        ->toBeString()
        ->not->toBe('');
});

test('login HeadQC berhasil dengan credential bawaan dari konfigurasi runtime', function () {
    $this->postJson('/api/v1/login', [
        'email' => config('qcontrol.headqc.email'),
        'password' => config('qcontrol.headqc.password_default'),
    ])
        ->assertSuccessful()
        ->assertJsonPath('data.profil.namaPengguna', config('qcontrol.headqc.nama_pengguna'))
        ->assertJsonPath('data.profil.peran', config('qcontrol.headqc.peran'));
});

test('login gagal dengan password yang salah', function () {
    $this->postJson('/api/v1/login', [
        'email' => config('qcontrol.headqc.email'),
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

test('login pengguna non HeadQC ditolak', function () {
    User::factory()->create([
        'email' => 'operator@pgn.local',
        'password' => 'password',
        'peran' => 'Operator',
    ]);

    $this->postJson('/api/v1/login', [
        'email' => 'operator@pgn.local',
        'password' => 'password',
    ])
        ->assertStatus(403)
        ->assertJsonPath('kesalahan.kode', 'AKSES_DITOLAK');
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
    $pengguna = User::query()->where('email', config('qcontrol.headqc.email'))->firstOrFail();
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
        ->assertJsonPath('data.namaPengguna', config('qcontrol.headqc.nama_pengguna'))
        ->assertJsonPath('data.email', config('qcontrol.headqc.email'))
        ->assertJsonPath('data.peran', config('qcontrol.headqc.peran'));
});

test('logout berhasil dengan token', function () {
    $pengguna = User::query()->where('email', config('qcontrol.headqc.email'))->firstOrFail();
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
    $pengguna = User::query()->where('email', config('qcontrol.headqc.email'))->firstOrFail();
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
