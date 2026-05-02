<?php

declare(strict_types=1);

use App\Application\Kesehatan\MembacaStatusKesehatanServer;
use App\Domain\Kesehatan\StatusKoneksiDatabase;
use App\Infrastructure\Kesehatan\PemeriksaKoneksiDatabase;
use Mockery\MockInterface;

test('layanan membaca status kesehatan membentuk payload sehat', function () {
    config()->set('app.name', 'PGNServer');
    config()->set('app.timezone', 'Asia/Jakarta');

    $pemeriksaKoneksiDatabase = Mockery::mock(PemeriksaKoneksiDatabase::class);
    $pemeriksaKoneksiDatabase
        ->shouldReceive('periksa')
        ->once()
        ->andReturn(new StatusKoneksiDatabase(
            status: 'terhubung',
            driver: 'pgsql',
        ));

    $layanan = new MembacaStatusKesehatanServer($pemeriksaKoneksiDatabase);
    $statusKesehatanServer = $layanan->jalankan();

    expect($statusKesehatanServer->status)->toBe('sehat')
        ->and($statusKesehatanServer->namaAplikasi)->toBe('PGNServer')
        ->and($statusKesehatanServer->versiApi)->toBe('v1')
        ->and($statusKesehatanServer->zonaWaktu)->toBe('Asia/Jakarta')
        ->and($statusKesehatanServer->koneksiDatabase->status)->toBe('terhubung')
        ->and($statusKesehatanServer->koneksiDatabase->driver)->toBe('pgsql')
        ->and($statusKesehatanServer->databaseTersedia())->toBeTrue();
});
