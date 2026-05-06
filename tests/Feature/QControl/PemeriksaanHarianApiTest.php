<?php

declare(strict_types=1);

use App\Models\QControlJenisDefect;
use App\Models\QControlLineProduksi;
use App\Models\QControlPart;
use App\Models\QControlPartJenisDefect;
use App\Models\QControlPemeriksaanDefectSlot;
use App\Models\QControlPemeriksaanHarian;
use App\Models\QControlPemeriksaanPart;
use App\Models\QControlSlotWaktu;
use Database\Seeders\DatabaseSeeder;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

uses(RefreshDatabase::class);

beforeEach(function (): void {
    $this->seed(DatabaseSeeder::class);
});

function headerAutentikasiHeadQc(TestCase $pengujian): array
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

function payloadPemeriksaanPress(array $timpa = []): array
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
        'clientDraftId' => 'draft-press-001',
        'tanggalProduksi' => '2026-05-05',
        'lineProduksiId' => $lineProduksi->id,
        'nomorDokumen' => 'FM-QA-025',
        'revisi' => '1',
        'catatan' => 'Pemeriksaan line PRESS',
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

function payloadPemeriksaanSewing(array $timpa = []): array
{
    $lineProduksi = QControlLineProduksi::query()
        ->where('kode_line', 'SEWING')
        ->firstOrFail();

    $part = QControlPart::query()
        ->where('kode_unik_part', 'FSB')
        ->firstOrFail();

    $slotWaktu = QControlSlotWaktu::query()
        ->where('kode_slot', 'SLOT_1300_1530')
        ->firstOrFail();

    $jenisDefect = QControlJenisDefect::query()
        ->where('kode_defect', 'SOBEK')
        ->firstOrFail();

    $relasiPartDefect = QControlPartJenisDefect::query()
        ->where('part_id', $part->id)
        ->where('jenis_defect_id', $jenisDefect->id)
        ->firstOrFail();

    return array_replace_recursive([
        'clientDraftId' => 'draft-sewing-001',
        'tanggalProduksi' => '2026-05-06',
        'lineProduksiId' => $lineProduksi->id,
        'nomorDokumen' => 'FM-QA-025',
        'revisi' => '1',
        'catatan' => 'Pemeriksaan line SEWING',
        'daftarPart' => [
            [
                'partId' => $part->id,
                'totalCheck' => 90,
                'daftarDefect' => [
                    [
                        'relasiPartDefectId' => $relasiPartDefect->id,
                        'slotWaktuId' => $slotWaktu->id,
                        'jumlahDefect' => 3,
                    ],
                ],
            ],
        ],
    ], $timpa);
}

test('endpoint pemeriksaan harian ditolak tanpa autentikasi', function () {
    $this->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress())
        ->assertStatus(401)
        ->assertJsonPath('kesalahan.kode', 'AUTENTIKASI_GAGAL');
});

test('endpoint pemeriksaan harian ditolak tanpa header idempotency', function () {
    $this->withHeaders(headerAutentikasiHeadQc($this))
        ->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress())
        ->assertStatus(422)
        ->assertJsonPath('kesalahan.kode', 'VALIDASI_GAGAL')
        ->assertJsonPath('kesalahan.detail.0.field', 'X-Idempotency-Key');
});

test('payload valid PRESS berhasil diterima', function () {
    $respons = $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-001'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress());

    $respons
        ->assertSuccessful()
        ->assertJsonPath('pesan', 'Pemeriksaan harian QControl berhasil diterima')
        ->assertJsonPath('data.lineProduksi.kodeLine', 'PRESS')
        ->assertJsonPath('data.totalCheck', 124)
        ->assertJsonPath('data.totalOk', 122)
        ->assertJsonPath('data.totalDefect', 2)
        ->assertJsonPath('data.rasioDefect', 1.61)
        ->assertJsonPath('data.jumlahPart', 1)
        ->assertJsonPath('data.jumlahBarisDefect', 1)
        ->assertJsonPath('data.duplikat', false);

    expect(QControlPemeriksaanHarian::query()->count())->toBe(1);
    expect(QControlPemeriksaanPart::query()->count())->toBe(1);
    expect(QControlPemeriksaanDefectSlot::query()->count())->toBe(1);

    $pemeriksaanHarian = QControlPemeriksaanHarian::query()->firstOrFail();
    $pemeriksaanPart = QControlPemeriksaanPart::query()->firstOrFail();
    $defectSlot = QControlPemeriksaanDefectSlot::query()->firstOrFail();

    expect($pemeriksaanHarian->kode_line_snapshot)->toBe('PRESS');
    expect($pemeriksaanHarian->nama_line_snapshot)->toBe('PRESS');
    expect($pemeriksaanHarian->nomor_dokumen_snapshot)->toBe('FM-QA-025');
    expect($pemeriksaanHarian->revisi_dokumen_snapshot)->toBe('1');
    expect($pemeriksaanHarian->nama_pic_snapshot)->toBe((string) config('qcontrol.headqc.nama_pengguna'));
    expect($pemeriksaanHarian->email_pic_snapshot)->toBe((string) config('qcontrol.headqc.email'));
    expect($pemeriksaanHarian->disiapkan_oleh_snapshot)->toBe((string) config('qcontrol.headqc.nama_pengguna'));
    expect($pemeriksaanPart->kode_unik_part_snapshot)->toBe('CB9');
    expect($pemeriksaanPart->nama_part_snapshot)->toBe('Carpet Console Box');
    expect($pemeriksaanPart->material_id_snapshot)->not->toBeNull();
    expect($pemeriksaanPart->kategori_ng_snapshot)->toBe('NG Material');
    expect($defectSlot->kode_tampilan_defect_snapshot)->toBe('A');
    expect($defectSlot->kode_defect_snapshot)->toBe('TERDAPAT_BENDA_ASING');
    expect($defectSlot->kategori_defect_snapshot)->toBe('NG Material');
    expect($defectSlot->label_slot_snapshot)->toBe('08.00 - 12.00');
});

test('total ok dihitung oleh server', function () {
    $respons = $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-002'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress([
        'daftarPart' => [
            [
                'totalCheck' => 50,
                'daftarDefect' => [
                    [
                        'jumlahDefect' => 7,
                    ],
                ],
            ],
        ],
    ]));

    $respons
        ->assertSuccessful()
        ->assertJsonPath('data.totalCheck', 50)
        ->assertJsonPath('data.totalOk', 43)
        ->assertJsonPath('data.totalDefect', 7);
});

test('total defect dihitung dari detail defect slot', function () {
    $part = QControlPart::query()
        ->where('kode_unik_part', 'CB9')
        ->firstOrFail();

    $slotSatu = QControlSlotWaktu::query()->where('kode_slot', 'SLOT_0800_1200')->firstOrFail();
    $slotDua = QControlSlotWaktu::query()->where('kode_slot', 'SLOT_1300_1530')->firstOrFail();
    $jenisDefectSatu = QControlJenisDefect::query()->where('kode_defect', 'TERDAPAT_BENDA_ASING')->firstOrFail();
    $jenisDefectDua = QControlJenisDefect::query()->where('kode_defect', 'PENYOK')->firstOrFail();
    $relasiSatu = QControlPartJenisDefect::query()->where('part_id', $part->id)->where('jenis_defect_id', $jenisDefectSatu->id)->firstOrFail();
    $relasiDua = QControlPartJenisDefect::query()->where('part_id', $part->id)->where('jenis_defect_id', $jenisDefectDua->id)->firstOrFail();

    $payload = payloadPemeriksaanPress([
        'daftarPart' => [
            [
                'partId' => $part->id,
                'totalCheck' => 30,
                'daftarDefect' => [
                    [
                        'relasiPartDefectId' => $relasiSatu->id,
                        'slotWaktuId' => $slotSatu->id,
                        'jumlahDefect' => 2,
                    ],
                    [
                        'relasiPartDefectId' => $relasiDua->id,
                        'slotWaktuId' => $slotDua->id,
                        'jumlahDefect' => 5,
                    ],
                ],
            ],
        ],
    ]);

    $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-003'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', $payload)
        ->assertSuccessful()
        ->assertJsonPath('data.totalDefect', 7)
        ->assertJsonPath('data.totalOk', 23)
        ->assertJsonPath('data.jumlahBarisDefect', 2);
});

test('total defect lebih besar dari total check ditolak', function () {
    $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-004'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress([
        'daftarPart' => [
            [
                'totalCheck' => 1,
                'daftarDefect' => [
                    [
                        'jumlahDefect' => 2,
                    ],
                ],
            ],
        ],
    ]))
        ->assertStatus(422)
        ->assertJsonPath('kesalahan.kode', 'TOTAL_DEFECT_MELEBIHI_TOTAL_CHECK');
});

test('relasi defect yang bukan milik part ditolak', function () {
    $partSewing = QControlPart::query()
        ->where('kode_unik_part', 'FSB')
        ->firstOrFail();

    $jenisDefectSewing = QControlJenisDefect::query()
        ->where('kode_defect', 'SOBEK')
        ->firstOrFail();

    $relasiSewing = QControlPartJenisDefect::query()
        ->where('part_id', $partSewing->id)
        ->where('jenis_defect_id', $jenisDefectSewing->id)
        ->firstOrFail();

    $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-005'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress([
        'daftarPart' => [
            [
                'daftarDefect' => [
                    [
                        'relasiPartDefectId' => $relasiSewing->id,
                    ],
                ],
            ],
        ],
    ]))
        ->assertStatus(422)
        ->assertJsonPath('kesalahan.kode', 'TEMPLATE_DEFECT_TIDAK_VALID');
});

test('part beda line ditolak', function () {
    $partSewing = QControlPart::query()
        ->where('kode_unik_part', 'FSB')
        ->firstOrFail();

    $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-005a'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress([
        'daftarPart' => [
            [
                'partId' => $partSewing->id,
                'totalCheck' => 10,
                'daftarDefect' => [],
            ],
        ],
    ]))
        ->assertStatus(422)
        ->assertJsonPath('kesalahan.kode', 'TEMPLATE_DEFECT_TIDAK_VALID');
});

test('slot waktu tidak aktif ditolak', function () {
    $slotWaktu = QControlSlotWaktu::query()
        ->where('kode_slot', 'SLOT_0800_1200')
        ->firstOrFail();

    $slotWaktu->forceFill(['aktif' => false])->save();

    $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-006'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress())
        ->assertStatus(422)
        ->assertJsonPath('kesalahan.kode', 'TEMPLATE_DEFECT_TIDAK_VALID');
});

test('duplicate part dalam payload ditolak', function () {
    $payload = payloadPemeriksaanPress();
    $payload['daftarPart'][] = $payload['daftarPart'][0];

    $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-007'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', $payload)
        ->assertStatus(422)
        ->assertJsonPath('kesalahan.kode', 'VALIDASI_GAGAL');
});

test('duplicate defect slot dalam part ditolak', function () {
    $payload = payloadPemeriksaanPress();
    $payload['daftarPart'][0]['daftarDefect'][] = $payload['daftarPart'][0]['daftarDefect'][0];

    $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-008'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', $payload)
        ->assertStatus(422)
        ->assertJsonPath('kesalahan.kode', 'VALIDASI_GAGAL');
});

test('retry idempotency key dengan payload sama berhasil dan tidak insert ulang', function () {
    $header = array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-009'],
    );
    $payload = payloadPemeriksaanPress();

    $this->withHeaders($header)
        ->postJson('/api/v1/qcontrol/pemeriksaan-harian', $payload)
        ->assertSuccessful()
        ->assertJsonPath('data.duplikat', false);

    $this->withHeaders($header)
        ->postJson('/api/v1/qcontrol/pemeriksaan-harian', $payload)
        ->assertSuccessful()
        ->assertJsonPath('data.duplikat', true);

    expect(QControlPemeriksaanHarian::query()->count())->toBe(1);
    expect(QControlPemeriksaanPart::query()->count())->toBe(1);
    expect(QControlPemeriksaanDefectSlot::query()->count())->toBe(1);
});

test('retry idempotency key dengan payload berbeda ditolak', function () {
    $header = array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-010'],
    );

    $this->withHeaders($header)
        ->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress())
        ->assertSuccessful();

    $this->withHeaders($header)
        ->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress([
            'daftarPart' => [
                [
                    'totalCheck' => 200,
                ],
            ],
        ]))
        ->assertStatus(409)
        ->assertJsonPath('kesalahan.kode', 'KONFLIK_IDEMPOTENCY');
});

test('tanggal dan line yang sama dengan idempotency berbeda tetap bisa disimpan', function () {
    $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-011'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress())
        ->assertSuccessful();

    $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-press-012'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPress([
        'daftarPart' => [
            [
                'totalCheck' => 999,
            ],
        ],
    ]))
        ->assertSuccessful()
        ->assertJsonPath('data.totalCheck', 999)
        ->assertJsonPath('data.duplikat', false);

    expect(QControlPemeriksaanHarian::query()->count())->toBe(2);
});

test('payload SEWING valid berhasil diterima', function () {
    $this->withHeaders(array_merge(
        headerAutentikasiHeadQc($this),
        ['X-Idempotency-Key' => 'pemeriksaan-sewing-001'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanSewing())
        ->assertSuccessful()
        ->assertJsonPath('data.lineProduksi.kodeLine', 'SEWING')
        ->assertJsonPath('data.totalCheck', 90)
        ->assertJsonPath('data.totalOk', 87)
        ->assertJsonPath('data.totalDefect', 3);
});
