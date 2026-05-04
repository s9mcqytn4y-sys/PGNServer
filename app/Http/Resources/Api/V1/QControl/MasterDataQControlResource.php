<?php

declare(strict_types=1);

namespace App\Http\Resources\Api\V1\QControl;

use App\Models\QControlJenisDefect;
use App\Models\QControlKategoriDefect;
use App\Models\QControlLineProduksi;
use App\Models\QControlMaterial;
use App\Models\QControlPart;
use App\Models\QControlPartJenisDefect;
use App\Models\QControlSlotWaktu;
use Illuminate\Database\Eloquent\Collection as KoleksiEloquent;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

final class MasterDataQControlResource extends JsonResource
{
    /**
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        /** @var array{
         *     versiMasterData: string,
         *     lineProduksi: KoleksiEloquent<int, QControlLineProduksi>,
         *     slotWaktu: KoleksiEloquent<int, QControlSlotWaktu>,
         *     material: KoleksiEloquent<int, QControlMaterial>,
         *     part: KoleksiEloquent<int, QControlPart>,
         *     kategoriDefect: KoleksiEloquent<int, QControlKategoriDefect>,
         *     jenisDefect: KoleksiEloquent<int, QControlJenisDefect>,
         *     relasiPartDefect: KoleksiEloquent<int, QControlPartJenisDefect>
         * } $data
         */
        $data = $this->resource;

        return [
            'versiMasterData' => $data['versiMasterData'],
            'lineProduksi' => $data['lineProduksi']->map(fn (QControlLineProduksi $lineProduksi) => [
                'id' => $lineProduksi->id,
                'kodeLine' => $lineProduksi->kode_line,
                'namaLine' => $lineProduksi->nama_line,
                'aktif' => $lineProduksi->aktif,
                'urutanTampil' => $lineProduksi->urutan_tampil,
            ])->values()->all(),
            'slotWaktu' => $data['slotWaktu']->map(fn (QControlSlotWaktu $slotWaktu) => [
                'id' => $slotWaktu->id,
                'kodeSlot' => $slotWaktu->kode_slot,
                'labelSlot' => $slotWaktu->label_slot,
                'jamMulai' => $slotWaktu->jam_mulai?->format('H:i:s'),
                'jamSelesai' => $slotWaktu->jam_selesai?->format('H:i:s'),
                'aktif' => $slotWaktu->aktif,
                'urutanTampil' => $slotWaktu->urutan_tampil,
            ])->values()->all(),
            'material' => $data['material']->map(fn (QControlMaterial $material) => [
                'id' => $material->id,
                'kodeMaterial' => $material->kode_material,
                'namaMaterial' => $material->nama_material,
                'aktif' => $material->aktif,
            ])->values()->all(),
            'part' => $data['part']->map(fn (QControlPart $part) => [
                'id' => $part->id,
                'kodeUnikPart' => $part->kode_unik_part,
                'namaPart' => $part->nama_part,
                'nomorPart' => $part->nomor_part,
                'materialId' => $part->material_id,
                'kodeMaterial' => $part->materialTerkait?->kode_material,
                'namaMaterial' => $part->materialTerkait?->nama_material,
                'kodeProyek' => $part->kode_proyek,
                'jumlahItemPerKanban' => $part->jumlah_item_per_kanban,
                'lineDefaultId' => $part->line_default_id,
                'kodeLineDefault' => $part->lineProduksiDefault?->kode_line,
                'namaLineDefault' => $part->lineProduksiDefault?->nama_line,
                'aktif' => $part->aktif,
                'sumberData' => $part->sumber_data,
            ])->values()->all(),
            'kategoriDefect' => $data['kategoriDefect']->map(fn (QControlKategoriDefect $kategoriDefect) => [
                'id' => $kategoriDefect->id,
                'kodeKategori' => $kategoriDefect->kode_kategori,
                'namaKategori' => $kategoriDefect->nama_kategori,
                'aktif' => $kategoriDefect->aktif,
                'urutanTampil' => $kategoriDefect->urutan_tampil,
            ])->values()->all(),
            'jenisDefect' => $data['jenisDefect']->map(fn (QControlJenisDefect $jenisDefect) => [
                'id' => $jenisDefect->id,
                'kodeDefect' => $jenisDefect->kode_defect,
                'namaDefect' => $jenisDefect->nama_defect,
                'kategoriDefectId' => $jenisDefect->kategori_defect_id,
                'kodeKategori' => $jenisDefect->kategoriDefectTerkait?->kode_kategori,
                'namaKategori' => $jenisDefect->kategoriDefectTerkait?->nama_kategori,
                'aktif' => $jenisDefect->aktif,
            ])->values()->all(),
            'relasiPartDefect' => $data['relasiPartDefect']->map(fn (QControlPartJenisDefect $relasiPartDefect) => [
                'id' => $relasiPartDefect->id,
                'partId' => $relasiPartDefect->part_id,
                'kodeUnikPart' => $relasiPartDefect->partTerkait?->kode_unik_part,
                'jenisDefectId' => $relasiPartDefect->jenis_defect_id,
                'kodeDefect' => $relasiPartDefect->jenisDefectTerkait?->kode_defect,
                'urutanTampil' => $relasiPartDefect->urutan_tampil,
                'aktif' => $relasiPartDefect->aktif,
            ])->values()->all(),
        ];
    }
}
