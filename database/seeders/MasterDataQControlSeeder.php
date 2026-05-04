<?php

declare(strict_types=1);

namespace Database\Seeders;

use App\Models\QControlJenisDefect;
use App\Models\QControlKategoriDefect;
use App\Models\QControlLineProduksi;
use App\Models\QControlMaterial;
use App\Models\QControlPart;
use App\Models\QControlPartJenisDefect;
use App\Models\QControlSlotWaktu;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Seeder;
use Illuminate\Support\Str;

final class MasterDataQControlSeeder extends Seeder
{
    use WithoutModelEvents;

    public function run(): void
    {
        $lineProduksi = [];
        foreach ($this->daftarLineProduksi() as $dataLineProduksi) {
            $lineProduksi[$dataLineProduksi['kode_line']] = $this->simpanModel(
                QControlLineProduksi::class,
                ['kode_line' => $dataLineProduksi['kode_line']],
                $dataLineProduksi,
            );
        }

        foreach ($this->daftarSlotWaktu() as $dataSlotWaktu) {
            $this->simpanModel(
                QControlSlotWaktu::class,
                ['kode_slot' => $dataSlotWaktu['kode_slot']],
                $dataSlotWaktu,
            );
        }

        $material = [];
        foreach ($this->daftarMaterial() as $dataMaterial) {
            $material[$dataMaterial['nama_material']] = $this->simpanModel(
                QControlMaterial::class,
                ['nama_material' => $dataMaterial['nama_material']],
                $dataMaterial,
            );
        }

        $kategoriDefect = [];
        foreach ($this->daftarKategoriDefect() as $dataKategoriDefect) {
            $kategoriDefect[$dataKategoriDefect['kode_kategori']] = $this->simpanModel(
                QControlKategoriDefect::class,
                ['kode_kategori' => $dataKategoriDefect['kode_kategori']],
                $dataKategoriDefect,
            );
        }

        $jenisDefect = [];
        foreach ($this->daftarJenisDefect() as $dataJenisDefect) {
            $kodeKategori = $dataJenisDefect['kode_kategori'];
            $dataSimpan = [
                'kode_defect' => $dataJenisDefect['kode_defect'],
                'nama_defect' => $dataJenisDefect['nama_defect'],
                'kategori_defect_id' => $kategoriDefect[$kodeKategori]->id ?? null,
                'aktif' => true,
            ];

            $jenisDefect[$dataJenisDefect['kode_defect']] = $this->simpanModel(
                QControlJenisDefect::class,
                ['kode_defect' => $dataJenisDefect['kode_defect']],
                $dataSimpan,
            );
        }

        $part = [];
        foreach ($this->daftarPart($material, $lineProduksi) as $dataPart) {
            $part[$dataPart['kode_unik_part']] = $this->simpanModel(
                QControlPart::class,
                ['kode_unik_part' => $dataPart['kode_unik_part']],
                $dataPart,
            );
        }

        foreach ($this->daftarRelasiPartJenisDefect() as $urutan => [$kodePart, $kodeDefect, $kodeTampilan]) {
            $this->simpanModel(
                QControlPartJenisDefect::class,
                [
                    'part_id' => $part[$kodePart]->id,
                    'jenis_defect_id' => $jenisDefect[$kodeDefect]->id,
                ],
                [
                    'kode_tampilan_defect' => $kodeTampilan,
                    'urutan_tampil' => $urutan + 1,
                    'aktif' => true,
                ],
            );
        }
    }

    /**
     * @return list<array{kode_line: string, nama_line: string, aktif: bool, urutan_tampil: int}>
     */
    private function daftarLineProduksi(): array
    {
        return [
            [
                'kode_line' => 'PRESS',
                'nama_line' => 'PRESS',
                'aktif' => true,
                'urutan_tampil' => 1,
            ],
            [
                'kode_line' => 'SEWING',
                'nama_line' => 'SEWING',
                'aktif' => true,
                'urutan_tampil' => 2,
            ],
        ];
    }

    /**
     * @return list<array{kode_slot: string, label_slot: string, jam_mulai: string|null, jam_selesai: string|null, urutan_tampil: int, aktif: bool}>
     */
    private function daftarSlotWaktu(): array
    {
        return [
            [
                'kode_slot' => 'SLOT_0800_1200',
                'label_slot' => '08.00 - 12.00',
                'jam_mulai' => '08:00:00',
                'jam_selesai' => '12:00:00',
                'urutan_tampil' => 1,
                'aktif' => true,
            ],
            [
                'kode_slot' => 'SLOT_1300_1530',
                'label_slot' => '13.00 - 15.30',
                'jam_mulai' => '13:00:00',
                'jam_selesai' => '15:30:00',
                'urutan_tampil' => 2,
                'aktif' => true,
            ],
            [
                'kode_slot' => 'SLOT_1600_1730',
                'label_slot' => '16.00 - 17.30',
                'jam_mulai' => '16:00:00',
                'jam_selesai' => '17:30:00',
                'urutan_tampil' => 3,
                'aktif' => true,
            ],
            [
                'kode_slot' => 'SLOT_1830_SELESAI',
                'label_slot' => '18.30 - Selesai',
                'jam_mulai' => '18:30:00',
                'jam_selesai' => null,
                'urutan_tampil' => 4,
                'aktif' => true,
            ],
        ];
    }

    /**
     * @return list<array{kode_material: string, nama_material: string, aktif: bool}>
     */
    private function daftarMaterial(): array
    {
        return [
            ['kode_material' => 'MAT_CARPET', 'nama_material' => 'CARPET', 'aktif' => true],
            ['kode_material' => 'MAT_SPONGE', 'nama_material' => 'SPONGE', 'aktif' => true],
            ['kode_material' => 'MAT_LEATHER', 'nama_material' => 'LEATHER', 'aktif' => true],
            ['kode_material' => 'MAT_FOAM', 'nama_material' => 'FOAM', 'aktif' => true],
            ['kode_material' => 'MAT_PLASTIC_CLIP', 'nama_material' => 'PLASTIC CLIP', 'aktif' => true],
            ['kode_material' => 'MAT_KAIN_PELAPIS', 'nama_material' => 'KAIN PELAPIS', 'aktif' => true],
            ['kode_material' => 'MAT_KARET_PELINDUNG', 'nama_material' => 'KARET PELINDUNG', 'aktif' => true],
            ['kode_material' => 'MAT_HARDFELT', 'nama_material' => 'HARDFELT', 'aktif' => true],
            ['kode_material' => 'MAT_EPDM', 'nama_material' => 'EPDM', 'aktif' => true],
            ['kode_material' => 'MAT_FUJISEAT', 'nama_material' => 'FUJISEAT', 'aktif' => true],
            ['kode_material' => 'MAT_SILENCER', 'nama_material' => 'SILENCER', 'aktif' => true],
            ['kode_material' => 'MAT_ESTER_CANVAS', 'nama_material' => 'ESTER CANVAS', 'aktif' => true],
            ['kode_material' => 'MAT_QUEENSCORD', 'nama_material' => 'QUEENSCORD', 'aktif' => true],
            ['kode_material' => 'MAT_LAINNYA', 'nama_material' => 'LAINNYA', 'aktif' => true],
        ];
    }

    /**
     * @return list<array{kode_kategori: string, nama_kategori: string, aktif: bool, urutan_tampil: int}>
     */
    private function daftarKategoriDefect(): array
    {
        return [
            ['kode_kategori' => 'MATERIAL', 'nama_kategori' => 'Material', 'aktif' => true, 'urutan_tampil' => 1],
            ['kode_kategori' => 'PROSES_PRESS', 'nama_kategori' => 'Proses Press', 'aktif' => true, 'urutan_tampil' => 2],
            ['kode_kategori' => 'PROSES_SEWING', 'nama_kategori' => 'Proses Sewing', 'aktif' => true, 'urutan_tampil' => 3],
            ['kode_kategori' => 'LAINNYA', 'nama_kategori' => 'Lainnya', 'aktif' => true, 'urutan_tampil' => 4],
        ];
    }

    /**
     * @return list<array{kode_defect: string, nama_defect: string, kode_kategori: string}>
     */
    private function daftarJenisDefect(): array
    {
        return [
            // PRESS
            ['kode_defect' => 'LAMINASI_BERKERUT', 'nama_defect' => 'Laminasi Berkerut', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'LAMINASI_BOLONG', 'nama_defect' => 'Laminasi Bolong', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'LAMINASI_TIDAK_MATANG', 'nama_defect' => 'Laminasi Tidak Matang', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'TERDAPAT_BENDA_ASING', 'nama_defect' => 'Terdapat Benda Asing', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'BAHAN_TIPIS', 'nama_defect' => 'Bahan Tipis', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'POTONGAN_BERLEBIH', 'nama_defect' => 'Potongan Berlebih', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'LAMINASI_TERSOBEK', 'nama_defect' => 'Laminasi Tersobek', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'DIMENSI_TIDAK_SESUAI', 'nama_defect' => 'Dimensi Tidak Sesuai', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'PENYOK', 'nama_defect' => 'Penyok', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'KARPET_BERJAMUR', 'nama_defect' => 'Karpet Berjamur', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'KARPET_TIPIS', 'nama_defect' => 'Karpet Tipis', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'SOBEK', 'nama_defect' => 'Sobek', 'kode_kategori' => 'PROSES_PRESS'],

            // SEWING
            ['kode_defect' => 'BRUDUL', 'nama_defect' => 'Brudul', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'SPUNBOND_TIDAK_MEREKAT', 'nama_defect' => 'Spunbond Tidak Merekat', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'SPUNBOND_TERLIPAT', 'nama_defect' => 'Spunbond Terlipat', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'SPUNBOND_HARDEN', 'nama_defect' => 'Spunbond Harden', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'SPUNBOND_KOTOR', 'nama_defect' => 'Spunbond Kotor', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'SPUNBOND_TERPOTONG', 'nama_defect' => 'Spunbond Terpotong', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'LAMINATING_TIDAK_MATANG', 'nama_defect' => 'Laminating Tidak Matang', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'LAMINATING_BOLONG', 'nama_defect' => 'Laminating Bolong', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'TERBALIK', 'nama_defect' => 'Terbalik', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'OVERCUTTING', 'nama_defect' => 'Overcutting', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'SEWING_MIRING', 'nama_defect' => 'Sewing Miring', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'MARGIN_SEWING', 'nama_defect' => 'Margin Sewing', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'MARGIN_OUT_DIMENSI', 'nama_defect' => 'Margin Out Dimensi', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'BACKSTITCH_KURANG_DARI_15MM', 'nama_defect' => 'Backstitch Kurang dari 15mm', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'DENT', 'nama_defect' => 'Dent', 'kode_kategori' => 'PROSES_PRESS'],
        ];
    }

    /**
     * @param  array<string, QControlMaterial>  $material
     * @param  array<string, QControlLineProduksi>  $lineProduksi
     * @return list<array<string, string|int|bool|null>>
     */
    private function daftarPart(array $material, array $lineProduksi): array
    {
        return [
            [
                'kode_unik_part' => 'CR6',
                'nama_part' => 'Carpet RR Seat No. 2 RH',
                'nomor_part' => '72996-X7H00',
                'material_id' => $material['CARPET']->id,
                'kode_proyek' => '560B',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['PRESS']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Press.xlsx',
            ],
            [
                'kode_unik_part' => 'CL7',
                'nama_part' => 'Carpet RR Seat No. 2 LH',
                'nomor_part' => '72997-X7H00',
                'material_id' => $material['CARPET']->id,
                'kode_proyek' => '560B',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['PRESS']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Press.xlsx',
            ],
            [
                'kode_unik_part' => 'BT136',
                'nama_part' => 'Felt FR Back RH',
                'nomor_part' => '11101-A1211',
                'material_id' => $material['HARDFELT']->id,
                'kode_proyek' => '560B',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['PRESS']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Press.xlsx',
            ],
            [
                'kode_unik_part' => 'BT137',
                'nama_part' => 'Felt FR Back RH',
                'nomor_part' => '11102-A1211',
                'material_id' => $material['HARDFELT']->id,
                'kode_proyek' => '560B',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['PRESS']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Press.xlsx',
            ],
            [
                'kode_unik_part' => 'BT144',
                'nama_part' => 'Felt FR Back LH',
                'nomor_part' => '12101-A1211',
                'material_id' => $material['HARDFELT']->id,
                'kode_proyek' => '560B',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['PRESS']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Press.xlsx',
            ],
            [
                'kode_unik_part' => 'BM7',
                'nama_part' => 'Pad RR Seat Back RH',
                'nomor_part' => '71651-BZ020',
                'material_id' => $material['HARDFELT']->id,
                'kode_proyek' => 'D25',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['PRESS']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Press.xlsx',
            ],
            [
                'kode_unik_part' => 'BM8',
                'nama_part' => 'Pad RR Seat Back LH',
                'nomor_part' => '71652-BZ020',
                'material_id' => $material['HARDFELT']->id,
                'kode_proyek' => 'D25',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['PRESS']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Press.xlsx',
            ],
            [
                'kode_unik_part' => 'CB9',
                'nama_part' => 'Carpet Console Box',
                'nomor_part' => '58815-KK010',
                'material_id' => $material['CARPET']->id,
                'kode_proyek' => '650/660A',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['PRESS']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Press.xlsx',
            ],
            [
                'kode_unik_part' => 'FJ0',
                'nama_part' => 'Pad Seaten RH',
                'nomor_part' => '71075-F1V01',
                'material_id' => $material['FUJISEAT']->id,
                'kode_proyek' => 'D14N',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['PRESS']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Press.xlsx',
            ],
            [
                'kode_unik_part' => 'FJ1',
                'nama_part' => 'Pad Seaten LH',
                'nomor_part' => '71075-F1V02',
                'material_id' => $material['FUJISEAT']->id,
                'kode_proyek' => 'D14N',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['PRESS']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Press.xlsx',
            ],
            [
                'kode_unik_part' => 'FSB',
                'nama_part' => 'Felt Seat Back',
                'nomor_part' => '79977-BZO20',
                'material_id' => $material['LAINNYA']->id,
                'kode_proyek' => null,
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['SEWING']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Sewing.xlsx',
            ],
            [
                'kode_unik_part' => 'CFRSH',
                'nama_part' => 'Cover FR Seat Hinge',
                'nomor_part' => '71831-BZ150',
                'material_id' => $material['LAINNYA']->id,
                'kode_proyek' => null,
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['SEWING']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Sewing.xlsx',
            ],
            [
                'kode_unik_part' => 'PRSB_RH_070',
                'nama_part' => 'Protector RR Seat Back RH',
                'nomor_part' => '71695-VT070',
                'material_id' => $material['LAINNYA']->id,
                'kode_proyek' => null,
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['SEWING']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Sewing.xlsx',
            ],
            [
                'kode_unik_part' => 'PRSB_LH_080',
                'nama_part' => 'Protector RR Seat Back LH',
                'nomor_part' => '71695-VT080',
                'material_id' => $material['LAINNYA']->id,
                'kode_proyek' => null,
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['SEWING']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Sewing.xlsx',
            ],
            [
                'kode_unik_part' => 'PRSB_RH_090',
                'nama_part' => 'Protector RR Seat Back RH',
                'nomor_part' => '71695-VT090',
                'material_id' => $material['LAINNYA']->id,
                'kode_proyek' => null,
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['SEWING']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Sewing.xlsx',
            ],
            [
                'kode_unik_part' => 'PRSB_LH_100',
                'nama_part' => 'Protector RR Seat Back LH',
                'nomor_part' => '71695-VT100',
                'material_id' => $material['LAINNYA']->id,
                'kode_proyek' => null,
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['SEWING']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Sewing.xlsx',
            ],
            [
                'kode_unik_part' => 'PRSB_RH_110',
                'nama_part' => 'Protector RR Seat Back RH',
                'nomor_part' => '71695-VT110',
                'material_id' => $material['LAINNYA']->id,
                'kode_proyek' => null,
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['SEWING']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Sewing.xlsx',
            ],
            [
                'kode_unik_part' => 'PRSB_LH_120',
                'nama_part' => 'Protector RR Seat Back LH',
                'nomor_part' => '71695-VT120',
                'material_id' => $material['LAINNYA']->id,
                'kode_proyek' => null,
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['SEWING']->id,
                'aktif' => true,
                'sumber_data' => 'Daily NG Sewing.xlsx',
            ],
        ];
    }

    /**
     * @return list<array{0: string, 1: string, 2: string}>
     */
    private function daftarRelasiPartJenisDefect(): array
    {
        $daftarRelasi = [];

        // PRESS CB9
        $defectsCB9 = [
            'A' => 'TERDAPAT_BENDA_ASING',
            'B' => 'PENYOK',
            'C' => 'KARPET_BERJAMUR',
            'D' => 'KARPET_TIPIS',
            'E' => 'POTONGAN_BERLEBIH',
            'F' => 'SOBEK',
            'G' => 'DIMENSI_TIDAK_SESUAI',
        ];
        foreach (['CB9'] as $kodePart) {
            foreach ($defectsCB9 as $kodeT => $kodeD) {
                $daftarRelasi[] = [$kodePart, $kodeD, $kodeT];
            }
        }

        // PRESS BT136, BT137, BT144, BM7, BM8, FJ0, FJ1
        $defectsPressLain = [
            'A' => 'LAMINASI_BERKERUT',
            'B' => 'LAMINASI_BOLONG',
            'C' => 'LAMINASI_TIDAK_MATANG',
            'D' => 'TERDAPAT_BENDA_ASING',
            'E' => 'BAHAN_TIPIS',
            'F' => 'POTONGAN_BERLEBIH',
            'G' => 'LAMINASI_TERSOBEK',
            'H' => 'DIMENSI_TIDAK_SESUAI',
        ];
        foreach (['BT136', 'BT137', 'BT144', 'BM7', 'BM8', 'FJ0', 'FJ1'] as $kodePart) {
            foreach ($defectsPressLain as $kodeT => $kodeD) {
                $daftarRelasi[] = [$kodePart, $kodeD, $kodeT];
            }
        }

        // SEWING Felt Seat Back
        $defectsFSB = [
            'A' => 'SOBEK',
            'B' => 'BRUDUL',
            'C' => 'SPUNBOND_TIDAK_MEREKAT',
            'D' => 'SPUNBOND_TERLIPAT',
            'E' => 'LAMINATING_TIDAK_MATANG',
            'F' => 'LAMINATING_BOLONG',
            'G' => 'SPUNBOND_TERPOTONG',
            'H' => 'TERBALIK',
            'I' => 'OVERCUTTING',
            'J' => 'SEWING_MIRING',
            'K' => 'MARGIN_OUT_DIMENSI',
            'L' => 'BACKSTITCH_KURANG_DARI_15MM',
        ];
        foreach (['FSB'] as $kodePart) {
            foreach ($defectsFSB as $kodeT => $kodeD) {
                $daftarRelasi[] = [$kodePart, $kodeD, $kodeT];
            }
        }

        // SEWING Cover FR Seat Hinge
        $defectsCFRSH = [
            'A' => 'SOBEK',
            'B' => 'BRUDUL',
            'C' => 'TERDAPAT_BENDA_ASING',
            'D' => 'KARPET_BERJAMUR',
            'E' => 'TERBALIK',
            'F' => 'OVERCUTTING',
            'G' => 'SEWING_MIRING',
        ];
        foreach (['CFRSH'] as $kodePart) {
            foreach ($defectsCFRSH as $kodeT => $kodeD) {
                $daftarRelasi[] = [$kodePart, $kodeD, $kodeT];
            }
        }

        // SEWING Protector RR Seat Back RH/LH semua part number
        $defectsProtector = [
            'A' => 'SOBEK',
            'B' => 'BRUDUL',
            'C' => 'SPUNBOND_TIDAK_MEREKAT',
            'D' => 'SPUNBOND_TERLIPAT',
            'E' => 'SPUNBOND_HARDEN',
            'F' => 'TERDAPAT_BENDA_ASING',
            'G' => 'SPUNBOND_KOTOR',
            'H' => 'LAMINATING_BOLONG',
            'I' => 'SPUNBOND_TERPOTONG',
            'J' => 'TERBALIK',
            'K' => 'OVERCUTTING',
            'L' => 'SEWING_MIRING',
            'M' => 'MARGIN_SEWING',
        ];
        $partProtectors = [
            'PRSB_RH_070',
            'PRSB_LH_080',
            'PRSB_RH_090',
            'PRSB_LH_100',
            'PRSB_RH_110',
            'PRSB_LH_120',
        ];
        foreach ($partProtectors as $kodePart) {
            foreach ($defectsProtector as $kodeT => $kodeD) {
                $daftarRelasi[] = [$kodePart, $kodeD, $kodeT];
            }
        }

        return $daftarRelasi;
    }

    /**
     * @param  class-string<Model>  $kelasModel
     * @param  array<string, mixed>  $kondisi
     * @param  array<string, mixed>  $atribut
     */
    private function simpanModel(string $kelasModel, array $kondisi, array $atribut): Model
    {
        /** @var Model $model */
        $model = $kelasModel::query()->firstOrNew($kondisi);

        if (! $model->exists) {
            $model->setAttribute('id', (string) Str::uuid());
        }

        $model->fill(array_merge($kondisi, $atribut));
        $model->save();

        return $model;
    }
}
