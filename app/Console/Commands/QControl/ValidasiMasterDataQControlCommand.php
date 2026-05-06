<?php

declare(strict_types=1);

namespace App\Console\Commands\QControl;

use App\Application\QControl\MemvalidasiMasterDataQControl;
use Illuminate\Console\Command;

/**
 * Menjalankan validasi otomatis master data QControl terhadap aturan fase aktif.
 */
final class ValidasiMasterDataQControlCommand extends Command
{
    protected $signature = 'qcontrol:validasi-master-data';

    protected $description = 'Memvalidasi kelengkapan, konsistensi, dan template master data QControl.';

    public function __construct(
        private readonly MemvalidasiMasterDataQControl $memvalidasiMasterDataQControl,
    ) {
        parent::__construct();
    }

    public function handle(): int
    {
        $hasilValidasi = $this->memvalidasiMasterDataQControl->jalankan();

        $this->info('Memulai validasi master data QControl...');
        $this->line('Line aktif: '.$hasilValidasi['ringkasan']['jumlahLineProduksiAktif']);
        $this->line('Slot aktif: '.$hasilValidasi['ringkasan']['jumlahSlotWaktuAktif']);
        $this->line('Part aktif: '.$hasilValidasi['ringkasan']['jumlahPartAktif']);
        $this->line('Jenis defect aktif: '.$hasilValidasi['ringkasan']['jumlahJenisDefectAktif']);
        $this->line('Relasi part-defect aktif: '.$hasilValidasi['ringkasan']['jumlahRelasiPartDefectAktif']);

        if ($hasilValidasi['valid']) {
            $this->info('Validasi master data QControl lulus tanpa temuan.');

            return self::SUCCESS;
        }

        $this->error('Validasi master data QControl gagal. Temuan:');

        foreach ($hasilValidasi['temuan'] as $nomor => $temuan) {
            $this->line(($nomor + 1).'. '.$temuan);
        }

        return self::FAILURE;
    }
}
