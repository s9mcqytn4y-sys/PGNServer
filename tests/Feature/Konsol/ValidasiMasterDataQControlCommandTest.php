<?php

declare(strict_types=1);

use App\Application\QControl\MemvalidasiMasterDataQControl;
use App\Models\QControlSlotWaktu;
use Database\Seeders\DatabaseSeeder;
use Illuminate\Foundation\Testing\RefreshDatabase;

uses(RefreshDatabase::class);

beforeEach(function (): void {
    $this->seed(DatabaseSeeder::class);
});

test('validator master data qcontrol lulus setelah seeder dijalankan', function () {
    $hasilValidasi = app(MemvalidasiMasterDataQControl::class)->jalankan();

    expect($hasilValidasi['valid'])->toBeTrue();
    expect($hasilValidasi['temuan'])->toBe([]);
});

test('command qcontrol validasi master data lulus', function () {
    $this->artisan('qcontrol:validasi-master-data')
        ->expectsOutputToContain('Validasi master data QControl lulus tanpa temuan.')
        ->assertSuccessful();
});

test('command qcontrol validasi master data gagal bila slot tidak sesuai', function () {
    $slot = QControlSlotWaktu::query()
        ->where('kode_slot', 'SLOT_1300_1530')
        ->firstOrFail();

    $slot->forceFill([
        'label_slot' => '13.30 - 16.00',
    ])->save();

    $this->artisan('qcontrol:validasi-master-data')
        ->expectsOutputToContain('Validasi master data QControl gagal.')
        ->assertFailed();
});
