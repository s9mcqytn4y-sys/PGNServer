<?php

declare(strict_types=1);

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Hash;

uses(RefreshDatabase::class);

test('command membuat HeadQC jika belum ada', function () {
    $emailHeadQC = (string) config('qcontrol.headqc.email');

    expect(User::query()->where('email', $emailHeadQC)->exists())->toBeFalse();

    $this->artisan('qcontrol:pastikan-headqc')
        ->expectsOutputToContain('Bootstrap HeadQC selesai.')
        ->expectsOutputToContain('Email: '.$emailHeadQC)
        ->expectsOutputToContain('Peran: HeadQC')
        ->expectsOutputToContain('Password default diatur ulang: ya')
        ->assertSuccessful();

    $pengguna = User::query()->where('email', $emailHeadQC)->first();

    expect($pengguna)->not->toBeNull();
    expect($pengguna?->peran)->toBe('HeadQC');
    expect(Hash::check((string) config('qcontrol.headqc.password_default'), (string) $pengguna?->password))->toBeTrue();
});

test('command memperbaiki role menjadi HeadQC', function () {
    $pengguna = User::factory()->create([
        'name' => 'Salah Role',
        'email' => (string) config('qcontrol.headqc.email'),
        'password' => Hash::make('kata-sandi-lama'),
        'peran' => 'BukanHeadQC',
    ]);

    $this->artisan('qcontrol:pastikan-headqc')
        ->assertSuccessful();

    $pengguna->refresh();

    expect($pengguna->peran)->toBe('HeadQC');
});

test('command reset password default valid', function () {
    $pengguna = User::factory()->create([
        'name' => 'HeadQC Lama',
        'email' => (string) config('qcontrol.headqc.email'),
        'password' => Hash::make('kata-sandi-lama'),
        'peran' => 'HeadQC',
    ]);

    $this->artisan('qcontrol:pastikan-headqc')
        ->expectsOutputToContain('Password default diatur ulang: ya')
        ->assertSuccessful();

    $pengguna->refresh();

    expect(Hash::check((string) config('qcontrol.headqc.password_default'), $pengguna->password))->toBeTrue();
});

test('command dengan tanpa reset password tidak menimpa password lama', function () {
    $kataSandiLama = 'kata-sandi-lama';

    $pengguna = User::factory()->create([
        'name' => 'HeadQC Lama',
        'email' => (string) config('qcontrol.headqc.email'),
        'password' => Hash::make($kataSandiLama),
        'peran' => 'PeranLama',
    ]);

    $this->artisan('qcontrol:pastikan-headqc --tanpa-reset-password')
        ->expectsOutputToContain('Password default diatur ulang: tidak')
        ->assertSuccessful();

    $pengguna->refresh();

    expect($pengguna->peran)->toBe('HeadQC');
    expect(Hash::check($kataSandiLama, $pengguna->password))->toBeTrue();
    expect(Hash::check((string) config('qcontrol.headqc.password_default'), $pengguna->password))->toBeFalse();
});

test('login API berhasil setelah command dijalankan', function () {
    $this->artisan('qcontrol:pastikan-headqc')
        ->assertSuccessful();

    $this->postJson('/api/v1/login', [
        'email' => config('qcontrol.headqc.email'),
        'password' => config('qcontrol.headqc.password_default'),
    ])
        ->assertSuccessful()
        ->assertJsonPath('data.profil.namaPengguna', config('qcontrol.headqc.nama_pengguna'))
        ->assertJsonPath('data.profil.peran', 'HeadQC');
});
