<?php

declare(strict_types=1);

use App\Models\QControlJenisDefect;
use App\Models\QControlLineProduksi;
use App\Models\QControlPart;
use App\Models\QControlPartJenisDefect;
use App\Models\QControlPemeriksaanHarian;
use App\Models\QControlSlotWaktu;
use App\Models\User;
use Database\Seeders\DatabaseSeeder;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Carbon;
use Illuminate\Support\Str;
use Tests\TestCase;

uses(RefreshDatabase::class);

beforeEach(function (): void {
    $this->seed(DatabaseSeeder::class);
});

function headerAutentikasiHeadQcBaca(TestCase $pengujian): array
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

function payloadPemeriksaanPressUntukBaca(array $timpa = []): array
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
        'clientDraftId' => 'baca-press-001',
        'tanggalProduksi' => '2026-05-08',
        'lineProduksiId' => $lineProduksi->id,
        'nomorDokumen' => 'FM-QA-025',
        'revisi' => '1',
        'catatan' => 'Data baca PRESS',
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

function simpanPemeriksaanPressUntukBaca(TestCase $pengujian, array $timpaPayload = [], string $kunciIdempotency = 'baca-pemeriksaan-001'): QControlPemeriksaanHarian
{
    $pengujian->withHeaders(array_merge(
        headerAutentikasiHeadQcBaca($pengujian),
        ['X-Idempotency-Key' => $kunciIdempotency],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPressUntukBaca($timpaPayload))
        ->assertSuccessful();

    return QControlPemeriksaanHarian::query()->latest('dibuat_pada')->firstOrFail();
}

test('endpoint list pemeriksaan harian ditolak tanpa autentikasi', function () {
    $this->getJson('/api/v1/qcontrol/pemeriksaan-harian')
        ->assertStatus(401)
        ->assertJsonPath('kesalahan.kode', 'AUTENTIKASI_GAGAL');
});

test('endpoint detail pemeriksaan harian ditolak tanpa autentikasi', function () {
    $lineProduksi = QControlLineProduksi::query()
        ->where('kode_line', 'PRESS')
        ->firstOrFail();
    $pengguna = User::query()->where('email', (string) config('qcontrol.headqc.email'))->firstOrFail();

    $pemeriksaanHarian = QControlPemeriksaanHarian::query()->create([
        'id' => (string) Str::uuid(),
        'tanggal_produksi' => '2026-05-08',
        'line_produksi_id' => $lineProduksi->id,
        'kode_line_snapshot' => $lineProduksi->kode_line,
        'nama_line_snapshot' => $lineProduksi->nama_line,
        'nomor_dokumen' => 'FM-QA-025',
        'revisi' => '1',
        'pengguna_headqc_id' => $pengguna->id,
        'status' => 'DITERIMA',
        'total_check' => 0,
        'total_ok' => 0,
        'total_defect' => 0,
        'rasio_defect' => 0,
        'diterima_pada' => Carbon::now(),
    ]);

    $this->getJson('/api/v1/qcontrol/pemeriksaan-harian/'.$pemeriksaanHarian->id)
        ->assertStatus(401)
        ->assertJsonPath('kesalahan.kode', 'AUTENTIKASI_GAGAL');
});

test('list kosong berhasil dimuat', function () {
    $this->withHeaders(headerAutentikasiHeadQcBaca($this))
        ->getJson('/api/v1/qcontrol/pemeriksaan-harian')
        ->assertSuccessful()
        ->assertJsonPath('data.daftarPemeriksaanHarian', [])
        ->assertJsonPath('metadata.jumlahData', 0);
});

test('setelah submit valid list menampilkan transaksi', function () {
    $pemeriksaanHarian = simpanPemeriksaanPressUntukBaca($this);

    $this->withHeaders(headerAutentikasiHeadQcBaca($this))
        ->getJson('/api/v1/qcontrol/pemeriksaan-harian')
        ->assertSuccessful()
        ->assertJsonPath('metadata.jumlahData', 1)
        ->assertJsonPath('data.daftarPemeriksaanHarian.0.id', $pemeriksaanHarian->id)
        ->assertJsonPath('data.daftarPemeriksaanHarian.0.lineProduksi.kodeLine', 'PRESS');
});

test('filter tanggalProduksi bekerja', function () {
    simpanPemeriksaanPressUntukBaca($this, [
        'clientDraftId' => 'baca-press-002',
        'tanggalProduksi' => '2026-05-09',
    ], 'baca-pemeriksaan-002');

    simpanPemeriksaanPressUntukBaca($this, [
        'clientDraftId' => 'baca-press-003',
        'tanggalProduksi' => '2026-05-10',
    ], 'baca-pemeriksaan-003');

    $this->withHeaders(headerAutentikasiHeadQcBaca($this))
        ->getJson('/api/v1/qcontrol/pemeriksaan-harian?tanggalProduksi=2026-05-09')
        ->assertSuccessful()
        ->assertJsonPath('metadata.jumlahData', 1)
        ->assertJsonPath('data.daftarPemeriksaanHarian.0.tanggalProduksi', '2026-05-09');
});

test('filter lineProduksiId bekerja', function () {
    $linePress = QControlLineProduksi::query()->where('kode_line', 'PRESS')->firstOrFail();
    $lineSewing = QControlLineProduksi::query()->where('kode_line', 'SEWING')->firstOrFail();

    simpanPemeriksaanPressUntukBaca($this, [
        'clientDraftId' => 'baca-press-004',
        'tanggalProduksi' => '2026-05-11',
    ], 'baca-pemeriksaan-004');

    $partSewing = QControlPart::query()->where('kode_unik_part', 'FSB')->firstOrFail();
    $slotWaktu = QControlSlotWaktu::query()->where('kode_slot', 'SLOT_1300_1530')->firstOrFail();
    $jenisDefect = QControlJenisDefect::query()->where('kode_defect', 'SOBEK')->firstOrFail();
    $relasiPartDefect = QControlPartJenisDefect::query()
        ->where('part_id', $partSewing->id)
        ->where('jenis_defect_id', $jenisDefect->id)
        ->firstOrFail();

    $this->withHeaders(array_merge(
        headerAutentikasiHeadQcBaca($this),
        ['X-Idempotency-Key' => 'baca-pemeriksaan-005'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', [
        'clientDraftId' => 'baca-sewing-001',
        'tanggalProduksi' => '2026-05-12',
        'lineProduksiId' => $lineSewing->id,
        'nomorDokumen' => 'FM-QA-025',
        'revisi' => '1',
        'daftarPart' => [
            [
                'partId' => $partSewing->id,
                'totalCheck' => 50,
                'daftarDefect' => [
                    [
                        'relasiPartDefectId' => $relasiPartDefect->id,
                        'slotWaktuId' => $slotWaktu->id,
                        'jumlahDefect' => 1,
                    ],
                ],
            ],
        ],
    ])->assertSuccessful();

    $this->withHeaders(headerAutentikasiHeadQcBaca($this))
        ->getJson('/api/v1/qcontrol/pemeriksaan-harian?lineProduksiId='.$linePress->id)
        ->assertSuccessful()
        ->assertJsonPath('metadata.jumlahData', 1)
        ->assertJsonPath('data.daftarPemeriksaanHarian.0.lineProduksi.kodeLine', 'PRESS');
});

test('detail menampilkan part dan defect slot', function () {
    $pemeriksaanHarian = simpanPemeriksaanPressUntukBaca($this);

    $this->withHeaders(headerAutentikasiHeadQcBaca($this))
        ->getJson('/api/v1/qcontrol/pemeriksaan-harian/'.$pemeriksaanHarian->id)
        ->assertSuccessful()
        ->assertJsonPath('data.lineProduksi.kodeLine', 'PRESS')
        ->assertJsonPath('data.daftarPart.0.kodeUnikPartSnapshot', 'CB9')
        ->assertJsonPath('data.daftarPart.0.daftarDefectSlot.0.kodeDefectSnapshot', 'TERDAPAT_BENDA_ASING')
        ->assertJsonPath('data.daftarPart.0.daftarDefectSlot.0.labelSlotSnapshot', '08.00 - 12.00');
});

test('detail memakai snapshot historis', function () {
    $pemeriksaanHarian = simpanPemeriksaanPressUntukBaca($this);

    $part = QControlPart::query()->where('kode_unik_part', 'CB9')->firstOrFail();
    $jenisDefect = QControlJenisDefect::query()->where('kode_defect', 'TERDAPAT_BENDA_ASING')->firstOrFail();
    $slotWaktu = QControlSlotWaktu::query()->where('kode_slot', 'SLOT_0800_1200')->firstOrFail();
    $line = QControlLineProduksi::query()->where('kode_line', 'PRESS')->firstOrFail();

    $part->forceFill([
        'nama_part' => 'NAMA BARU PART',
    ])->save();
    $jenisDefect->forceFill([
        'nama_defect' => 'NAMA BARU DEFECT',
        'kode_defect' => 'KODE_BARU_DEFECT',
    ])->save();
    $slotWaktu->forceFill([
        'label_slot' => 'LABEL SLOT BARU',
        'kode_slot' => 'KODE_SLOT_BARU',
    ])->save();
    $line->forceFill([
        'nama_line' => 'NAMA LINE BARU',
        'kode_line' => 'PRESS_BARU',
    ])->save();

    $this->withHeaders(headerAutentikasiHeadQcBaca($this))
        ->getJson('/api/v1/qcontrol/pemeriksaan-harian/'.$pemeriksaanHarian->id)
        ->assertSuccessful()
        ->assertJsonPath('data.lineProduksi.kodeLine', 'PRESS')
        ->assertJsonPath('data.lineProduksi.namaLine', 'PRESS')
        ->assertJsonPath('data.daftarPart.0.namaPartSnapshot', 'Carpet Console Box')
        ->assertJsonPath('data.daftarPart.0.daftarDefectSlot.0.kodeDefectSnapshot', 'TERDAPAT_BENDA_ASING')
        ->assertJsonPath('data.daftarPart.0.daftarDefectSlot.0.namaDefectSnapshot', 'Terdapat Benda Asing')
        ->assertJsonPath('data.daftarPart.0.daftarDefectSlot.0.labelSlotSnapshot', '08.00 - 12.00');
});

test('detail id tidak ditemukan mengembalikan envelope 404 standar', function () {
    $this->withHeaders(headerAutentikasiHeadQcBaca($this))
        ->getJson('/api/v1/qcontrol/pemeriksaan-harian/0f173e79-7f2e-4bb1-9664-82a0d95a4d2a')
        ->assertStatus(404)
        ->assertJsonPath('kesalahan.kode', 'DATA_TIDAK_DITEMUKAN');
});

test('limit maksimal divalidasi', function () {
    $this->withHeaders(headerAutentikasiHeadQcBaca($this))
        ->getJson('/api/v1/qcontrol/pemeriksaan-harian?limit=101')
        ->assertStatus(422)
        ->assertJsonPath('kesalahan.kode', 'VALIDASI_GAGAL');
});

test('endpoint post lama tetap lulus', function () {
    $this->withHeaders(array_merge(
        headerAutentikasiHeadQcBaca($this),
        ['X-Idempotency-Key' => 'baca-pemeriksaan-006'],
    ))->postJson('/api/v1/qcontrol/pemeriksaan-harian', payloadPemeriksaanPressUntukBaca([
        'clientDraftId' => 'baca-press-006',
        'tanggalProduksi' => '2026-05-13',
    ]))
        ->assertSuccessful()
        ->assertJsonPath('data.duplikat', false)
        ->assertJsonPath('data.totalOk', 122);
});
