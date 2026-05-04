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

        foreach ($this->daftarRelasiPartJenisDefect() as $urutan => [$kodePart, $kodeDefect]) {
            $this->simpanModel(
                QControlPartJenisDefect::class,
                [
                    'part_id' => $part[$kodePart]->id,
                    'jenis_defect_id' => $jenisDefect[$kodeDefect]->id,
                ],
                [
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
            ['kode_defect' => 'DENT', 'nama_defect' => 'DENT', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'GALER', 'nama_defect' => 'GALER', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'CARPET_TIPIS', 'nama_defect' => 'CARPET TIPIS', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'CARPET_BERJAMUR', 'nama_defect' => 'CARPET BERJAMUR', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'BELANG', 'nama_defect' => 'BELANG', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'HOLE_TA', 'nama_defect' => 'HOLE TA', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'OVERCUTTING', 'nama_defect' => 'OVERCUTTING', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'SOBEK', 'nama_defect' => 'SOBEK', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'BRUDUL', 'nama_defect' => 'BRUDUL', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'TERBALIK', 'nama_defect' => 'TERBALIK', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'KOTOR', 'nama_defect' => 'KOTOR', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'TERDAPAT_BENDA_ASING', 'nama_defect' => 'TERDAPAT BENDA ASING', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'LAMINATING_BERKERUT', 'nama_defect' => 'LAMINATING BERKERUT', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'LAMINATING_BOLONG', 'nama_defect' => 'LAMINATING BOLONG', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'LAMINATING_TIDAK_MATANG', 'nama_defect' => 'LAMINATING TIDAK MATANG', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'LAMINATING_TERSOBEK', 'nama_defect' => 'LAMINATING TERSOBEK', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'MATERIAL_TIPIS', 'nama_defect' => 'MATERIAL TIPIS', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'DIMENSI_OUT_STANDARD', 'nama_defect' => 'DIMENSI OUT STANDARD', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'UKURAN_TIDAK_SESUAI', 'nama_defect' => 'UKURAN TIDAK SESUAI', 'kode_kategori' => 'MATERIAL'],
            ['kode_defect' => 'SEWING_MIRING', 'nama_defect' => 'SEWING MIRING', 'kode_kategori' => 'PROSES_SEWING'],
            ['kode_defect' => 'SPUNBOND_TIDAK_MEREKAT', 'nama_defect' => 'SPUNBOND TIDAK MEREKAT', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'SPUNBOND_TERLIPAT', 'nama_defect' => 'SPUNBOND TERLIPAT', 'kode_kategori' => 'PROSES_PRESS'],
            ['kode_defect' => 'POTONGAN_TIDAK_RATA', 'nama_defect' => 'POTONGAN TIDAK RATA', 'kode_kategori' => 'PROSES_PRESS'],
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
                'sumber_data' => 'audit_drive_curated',
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
                'sumber_data' => 'audit_drive_curated',
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
                'sumber_data' => 'audit_drive_curated',
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
                'sumber_data' => 'audit_drive_curated',
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
                'sumber_data' => 'audit_drive_curated',
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
                'sumber_data' => 'audit_drive_curated',
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
                'sumber_data' => 'audit_drive_curated',
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
                'sumber_data' => 'audit_drive_curated',
            ],
            [
                'kode_unik_part' => 'FJ0',
                'nama_part' => 'Pad Seaten RH',
                'nomor_part' => '71075-F1V01',
                'material_id' => $material['FUJISEAT']->id,
                'kode_proyek' => 'D14N',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['SEWING']->id,
                'aktif' => true,
                'sumber_data' => 'audit_drive_curated',
            ],
            [
                'kode_unik_part' => 'FJ1',
                'nama_part' => 'Pad Seaten LH',
                'nomor_part' => '71075-F1V02',
                'material_id' => $material['FUJISEAT']->id,
                'kode_proyek' => 'D14N',
                'jumlah_item_per_kanban' => null,
                'line_default_id' => $lineProduksi['SEWING']->id,
                'aktif' => true,
                'sumber_data' => 'audit_drive_curated',
            ],
        ];
    }

    /**
     * @return list<array{0: string, 1: string}>
     */
    private function daftarRelasiPartJenisDefect(): array
    {
        $daftarRelasi = [];

        foreach (['CR6', 'CL7', 'CB9'] as $kodePart) {
            foreach ([
                'DENT',
                'GALER',
                'CARPET_TIPIS',
                'CARPET_BERJAMUR',
                'BELANG',
                'HOLE_TA',
                'OVERCUTTING',
                'SOBEK',
            ] as $kodeDefect) {
                $daftarRelasi[] = [$kodePart, $kodeDefect];
            }
        }

        foreach (['BT136', 'BT137', 'BT144', 'BM7', 'BM8', 'FJ0', 'FJ1'] as $kodePart) {
            foreach ([
                'LAMINATING_BERKERUT',
                'LAMINATING_BOLONG',
                'LAMINATING_TIDAK_MATANG',
                'TERDAPAT_BENDA_ASING',
                'MATERIAL_TIPIS',
                'OVERCUTTING',
                'LAMINATING_TERSOBEK',
                'DIMENSI_OUT_STANDARD',
            ] as $kodeDefect) {
                $daftarRelasi[] = [$kodePart, $kodeDefect];
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
