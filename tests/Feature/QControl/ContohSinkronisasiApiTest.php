<?php

declare(strict_types=1);

use App\Models\IdempotencyKey;
use Database\Seeders\DatabaseSeeder;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

uses(RefreshDatabase::class);

beforeEach(function (): void {
    $this->seed(DatabaseSeeder::class);
});

function headerAutentikasiContohQcontrol(TestCase $pengujian): array
{
    $responsMasuk = $pengujian->postJson('/api/v1/login', [
        'email' => (string) config('qcontrol.headqc.email'),
        'password' => (string) config('qcontrol.headqc.password_default'),
    ]);

    $responsMasuk->assertSuccessful();

    return [
        'Authorization' => 'Bearer '.(string) $responsMasuk->json('data.token'),
        'Accept' => 'application/json',
    ];
}

test('endpoint contoh sinkronisasi qcontrol menerima payload dengan idempotency key', function () {
    $this->withHeaders(array_merge(
        headerAutentikasiContohQcontrol($this),
        ['X-Idempotency-Key' => 'qcontrol:contoh:test-001'],
    ))->postJson(
        '/api/v1/qcontrol/contoh',
        [
            'contoh' => true,
            'sumber' => 'manual',
        ],
    )
        ->assertSuccessful()
        ->assertJson([
            'berhasil' => true,
            'pesan' => 'Payload sinkronisasi QControl diterima',
            'metadata' => null,
            'kesalahan' => null,
        ])
        ->assertJsonPath('data.diterima', true)
        ->assertJsonPath('data.duplikat', false)
        ->assertJsonPath('data.idempotencyKey', 'qcontrol:contoh:test-001')
        ->assertJsonPath('data.endpoint', '/api/v1/qcontrol/contoh')
        ->assertJsonPath('data.mode', 'kontrak_awal');
});

test('endpoint contoh sinkronisasi qcontrol butuh autentikasi HeadQC', function () {
    $this->postJson('/api/v1/qcontrol/contoh', [
        'contoh' => true,
    ], [
        'X-Idempotency-Key' => 'qcontrol:contoh:tanpa-token',
    ])
        ->assertStatus(401)
        ->assertJsonPath('kesalahan.kode', 'AUTENTIKASI_GAGAL');
});

test('endpoint contoh sinkronisasi qcontrol menolak payload tanpa idempotency key', function () {
    $this->withHeaders(headerAutentikasiContohQcontrol($this))
        ->postJson('/api/v1/qcontrol/contoh', [
            'contoh' => true,
        ])
        ->assertStatus(422)
        ->assertJsonPath('kesalahan.kode', 'VALIDASI_GAGAL')
        ->assertJsonPath('kesalahan.detail.0.field', 'X-Idempotency-Key');
});

test('endpoint contoh sinkronisasi qcontrol mengembalikan sukses idempotent untuk duplicate request', function () {
    $headerIdempotency = array_merge(
        headerAutentikasiContohQcontrol($this),
        [
            'X-Idempotency-Key' => 'qcontrol:contoh:dedupe-001',
        ],
    );

    $responsePertama = $this->withHeaders($headerIdempotency)->postJson(
        '/api/v1/qcontrol/contoh',
        [
            'contoh' => true,
            'sumber' => 'fase_2b',
        ],
    );

    $responseKedua = $this->withHeaders($headerIdempotency)->postJson(
        '/api/v1/qcontrol/contoh',
        [
            'contoh' => true,
            'sumber' => 'fase_2b',
        ],
    );

    $responsePertama
        ->assertSuccessful()
        ->assertJsonPath('data.duplikat', false)
        ->assertJsonPath('data.idempotencyKey', 'qcontrol:contoh:dedupe-001');

    $responseKedua
        ->assertSuccessful()
        ->assertJson([
            'berhasil' => true,
            'pesan' => 'Payload sinkronisasi QControl sudah pernah diterima',
        ])
        ->assertJsonPath('data.duplikat', true)
        ->assertJsonPath('data.idempotencyKey', 'qcontrol:contoh:dedupe-001');

    expect(
        IdempotencyKey::query()
            ->where('kunci_idempotency', 'qcontrol:contoh:dedupe-001')
            ->count()
    )->toBe(1);
});
