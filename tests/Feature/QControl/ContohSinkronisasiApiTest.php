<?php

declare(strict_types=1);

use App\Models\IdempotencyKey;
use Illuminate\Foundation\Testing\RefreshDatabase;

uses(RefreshDatabase::class);

test('endpoint contoh sinkronisasi qcontrol menerima payload dengan idempotency key', function () {
    $this->postJson(
        '/api/v1/qcontrol/contoh',
        [
            'contoh' => true,
            'sumber' => 'manual',
        ],
        [
            'X-Idempotency-Key' => 'qcontrol:contoh:test-001',
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

test('endpoint contoh sinkronisasi qcontrol menolak payload tanpa idempotency key', function () {
    $this->postJson('/api/v1/qcontrol/contoh', [
        'contoh' => true,
    ])
        ->assertStatus(422)
        ->assertJson([
            'berhasil' => false,
            'pesan' => 'Header X-Idempotency-Key wajib diisi',
            'metadata' => null,
            'kesalahan' => [
                'kode' => 'VALIDASI_GAGAL',
                'detail' => [
                    [
                        'field' => 'X-Idempotency-Key',
                        'pesan' => 'Header X-Idempotency-Key wajib diisi',
                    ],
                ],
            ],
        ])
        ->assertJsonPath('data.payloadDiterima.contoh', true);
});

test('endpoint contoh sinkronisasi qcontrol mengembalikan sukses idempotent untuk duplicate request', function () {
    $headerIdempotency = [
        'X-Idempotency-Key' => 'qcontrol:contoh:dedupe-001',
    ];

    $responsePertama = $this->postJson(
        '/api/v1/qcontrol/contoh',
        [
            'contoh' => true,
            'sumber' => 'fase_2b',
        ],
        $headerIdempotency,
    );

    $responseKedua = $this->postJson(
        '/api/v1/qcontrol/contoh',
        [
            'contoh' => true,
            'sumber' => 'fase_2b',
        ],
        $headerIdempotency,
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
