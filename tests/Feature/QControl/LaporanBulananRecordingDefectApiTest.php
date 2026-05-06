<?php

declare(strict_types=1);

use App\Models\QControlJenisDefect;
use App\Models\QControlLineProduksi;
use App\Models\QControlPart;
use App\Models\QControlPartJenisDefect;
use App\Models\QControlSlotWaktu;
use Database\Seeders\DatabaseSeeder;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

uses(RefreshDatabase::class);

beforeEach(function (): void {
    $this->seed(DatabaseSeeder::class);
});

function headerAutentikasiHeadQcBulanan(TestCase $pengujian): array
{
    $responsMasuk = $pengujian->postJson('/api/v1/login', [
        'email' => (string) config('qcontrol.headqc.email'),
        'password' => (string) config('qcontrol.headqc.password_default'),
    ]);

    $responsMasuk
        ->assertSuccessful()
        ->assertJsonPath('data.profil.peran', 'HeadQC');

    return [
        'Authorization' => 'Bearer '.(string) $responsMasuk->json('data.token'),
        'Accept' => 'application/json',
    ];
}

function payloadBulananPress(array $timpa = []): array
{
    $lineProduksi = QControlLineProduksi::query()
        ->where('kode_line', 'PRESS')
        ->firstOrFail();

    $part = QControlPart::query()
        ->where('kode_unik_part', 'CB9')
        ->firstOrFail();

    $slotWaktu = QControlSlotWaktu::query()
        ->where('kode_slot', 'SLOT_0800_1200')
        ->firstOrFail();

    $jenisDefect = QControlJenisDefect::query()
        ->where('kode_defect', 'TERDAPAT_BENDA_ASING')
        ->firstOrFail();

    $relasiPartDefect = QControlPartJenisDefect::query()
        ->where('part_id', $part->id)
        ->where('jenis_defect_id', $jenisDefect->id)
        ->firstOrFail();

    return array_replace_recursive([
        'clientDraftId' => 'monthly-press-001',
        'tanggalProduksi' => '2026-05-06',
        'lineProduksiId' => $lineProduksi->id,
        'nomorDokumen' => 'FM-QA-025',
        'revisi' => '1',
        'catatan' => 'Data bulanan PRESS',
        'daftarPart' => [
            [
                'partId' => $part->id,
                'totalCheck' => 124,
                'daftarDefect' => [
                    [
                        'relasiPartDefectId' => $relasiPartDefect->id,
                        'slotWaktuId' => $slotWaktu->id,
                        'jumlahDefect' => 2,
                    ],
                ],
            ],
        ],
    ], $timpa);
}

test('monthly read model kosong tetap valid', function () {
    $lineProduksi = QControlLineProduksi::query()
        ->where('kode_line', 'PRESS')
        ->firstOrFail();

    $respons = $this->withHeaders(headerAutentikasiHeadQcBulanan($this))
        ->getJson('/api/v1/qcontrol/laporan-bulanan/recording-defect?bulan=5&tahun=2026&lineProduksiId='.$lineProduksi->id);

    $respons->assertSuccessful()
        ->assertJsonPath('data.bulan', 5)
        ->assertJsonPath('data.tahun', 2026)
        ->assertJsonPath('data.line.kodeLine', 'PRESS')
        ->assertJsonPath('data.daftarPart', [])
        ->assertJsonPath('data.totalBulanan', 0);

    /** @var array<string, mixed> $isiRespons */
    $isiRespons = json_decode($respons->getContent(), true, 512, JSON_THROW_ON_ERROR);
    /** @var array<string, int> $totalHarian */
    $totalHarian = $isiRespons['data']['totalHarian'];

    expect(count($totalHarian))->toBe(31);
    expect(reset($totalHarian))->toBe(0);
    expect(end($totalHarian))->toBe(0);
});

test('monthly read model menghitung total tanggal dan total bulanan dari daily', function () {
    $lineProduksi = QControlLineProduksi::query()
        ->where('kode_line', 'PRESS')
        ->firstOrFail();

    $headerAutentikasi = headerAutentikasiHeadQcBulanan($this);

    $headerPertama = array_merge(
        $headerAutentikasi,
        ['X-Idempotency-Key' => 'monthly-press-submit-001'],
    );
    $headerKedua = array_merge(
        $headerAutentikasi,
        ['X-Idempotency-Key' => 'monthly-press-submit-002'],
    );

    $this->withHeaders($headerPertama)
        ->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadBulananPress())
        ->assertSuccessful();

    $this->withHeaders($headerKedua)
        ->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadBulananPress([
            'clientDraftId' => 'monthly-press-002',
            'tanggalProduksi' => '2026-05-07',
            'daftarPart' => [
                [
                    'totalCheck' => 50,
                    'daftarDefect' => [
                        [
                            'jumlahDefect' => 1,
                        ],
                    ],
                ],
            ],
        ]))
        ->assertSuccessful();

    $respons = $this->withHeaders(headerAutentikasiHeadQcBulanan($this))
        ->getJson('/api/v1/qcontrol/laporan-bulanan/recording-defect?bulan=5&tahun=2026&lineProduksiId='.$lineProduksi->id);

    $respons->assertSuccessful()
        ->assertJsonPath('data.bulan', 5)
        ->assertJsonPath('data.daftarPart.0.kodeUnikPart', 'CB9')
        ->assertJsonPath('data.daftarPart.0.totalCheck', 174)
        ->assertJsonPath('data.daftarPart.0.totalDefect', 3)
        ->assertJsonPath('data.daftarPart.0.totalOk', 171)
        ->assertJsonPath('data.totalBulanan', 3);

    /** @var array<string, mixed> $isiRespons */
    $isiRespons = json_decode($respons->getContent(), true, 512, JSON_THROW_ON_ERROR);
    /** @var array<string, mixed> $data */
    $data = $isiRespons['data'];

    expect($data['daftarPart'][0]['daftarDefect'][0]['kodeTampilanDefect'])->toBe('A');
    expect($data['daftarPart'][0]['daftarDefect'][0]['jumlahPerTanggal']['6'])->toBe(2);
    expect($data['daftarPart'][0]['daftarDefect'][0]['jumlahPerTanggal']['7'])->toBe(1);
    expect($data['daftarPart'][0]['subtotalPerTanggal']['6'])->toBe(2);
    expect($data['daftarPart'][0]['subtotalPerTanggal']['7'])->toBe(1);
    expect($data['totalHarian']['6'])->toBe(2);
    expect($data['totalHarian']['7'])->toBe(1);
});
