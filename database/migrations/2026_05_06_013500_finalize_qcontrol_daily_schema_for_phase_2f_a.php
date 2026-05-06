<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::table('qcontrol_pemeriksaan_harian', function (Blueprint $table): void {
            $table->string('nomor_dokumen_snapshot')->default('FM-QA-025')->after('nama_line_snapshot');
            $table->string('revisi_dokumen_snapshot')->default('1')->after('nomor_dokumen_snapshot');
            $table->string('nama_pic_snapshot')->nullable()->after('pengguna_headqc_id');
            $table->string('email_pic_snapshot')->nullable()->after('nama_pic_snapshot');
            $table->string('disiapkan_oleh_snapshot')->nullable()->after('catatan');
            $table->string('diperiksa_oleh_snapshot')->nullable()->after('disiapkan_oleh_snapshot');
            $table->string('disetujui_oleh_snapshot')->nullable()->after('diperiksa_oleh_snapshot');
        });

        Schema::table('qcontrol_pemeriksaan_part', function (Blueprint $table): void {
            $table->uuid('material_id_snapshot')->nullable()->after('nama_part_snapshot');
            $table->foreign('material_id_snapshot', 'qcontrol_pemeriksaan_part_material_snapshot_foreign')
                ->references('id')
                ->on('qcontrol_material')
                ->nullOnDelete();
            $table->string('kategori_ng_snapshot')->nullable()->after('nama_material_snapshot');
        });

        Schema::table('qcontrol_pemeriksaan_defect_slot', function (Blueprint $table): void {
            $table->string('kategori_defect_snapshot')->nullable()->after('nama_defect_snapshot');
        });

        DB::table('qcontrol_pemeriksaan_harian')
            ->orderBy('id')
            ->get()
            ->each(function (object $baris): void {
                $pengguna = DB::table('users')
                    ->where('id', $baris->pengguna_headqc_id)
                    ->first();

                DB::table('qcontrol_pemeriksaan_harian')
                    ->where('id', $baris->id)
                    ->update([
                        'nomor_dokumen_snapshot' => $baris->nomor_dokumen ?? 'FM-QA-025',
                        'revisi_dokumen_snapshot' => $baris->revisi ?? '1',
                        'nama_pic_snapshot' => $pengguna?->name,
                        'email_pic_snapshot' => $pengguna?->email,
                        'disiapkan_oleh_snapshot' => $pengguna?->name,
                        'diperiksa_oleh_snapshot' => $pengguna?->name,
                        'disetujui_oleh_snapshot' => $pengguna?->name,
                    ]);
            });

        DB::table('qcontrol_pemeriksaan_part')
            ->orderBy('id')
            ->get()
            ->each(function (object $baris): void {
                $part = DB::table('qcontrol_part')
                    ->where('id', $baris->part_id)
                    ->first();

                DB::table('qcontrol_pemeriksaan_part')
                    ->where('id', $baris->id)
                    ->update([
                        'material_id_snapshot' => $part?->material_id,
                    ]);
            });

        DB::table('qcontrol_pemeriksaan_defect_slot')
            ->orderBy('id')
            ->get()
            ->each(function (object $baris): void {
                $kategoriDefectSnapshot = DB::table('qcontrol_part_jenis_defect')
                    ->leftJoin('qcontrol_jenis_defect', 'qcontrol_jenis_defect.id', '=', 'qcontrol_part_jenis_defect.jenis_defect_id')
                    ->leftJoin('qcontrol_kategori_defect', 'qcontrol_kategori_defect.id', '=', 'qcontrol_jenis_defect.kategori_defect_id')
                    ->where('qcontrol_part_jenis_defect.id', $baris->relasi_part_defect_id)
                    ->value('qcontrol_kategori_defect.nama_kategori');

                DB::table('qcontrol_pemeriksaan_defect_slot')
                    ->where('id', $baris->id)
                    ->update([
                        'kategori_defect_snapshot' => $kategoriDefectSnapshot,
                    ]);
            });

        Schema::table('qcontrol_pemeriksaan_harian', function (Blueprint $table): void {
            $table->dropUnique('qcontrol_pemeriksaan_harian_tanggal_produksi_line_produksi_id_unique');
            $table->dropIndex('qcontrol_pemeriksaan_harian_idempotency_key_index');

            $table->unique('idempotency_key');
            $table->index('tanggal_produksi');
            $table->index('line_produksi_id');
        });

        Schema::table('qcontrol_pemeriksaan_part', function (Blueprint $table): void {
            $table->index('pemeriksaan_harian_id');
            $table->index('part_id');
        });

        Schema::table('qcontrol_pemeriksaan_defect_slot', function (Blueprint $table): void {
            $table->index('pemeriksaan_part_id');
            $table->index('jenis_defect_id');
            $table->index('slot_waktu_id');
        });
    }

    public function down(): void
    {
        Schema::table('qcontrol_pemeriksaan_defect_slot', function (Blueprint $table): void {
            $table->dropIndex(['pemeriksaan_part_id']);
            $table->dropIndex(['jenis_defect_id']);
            $table->dropIndex(['slot_waktu_id']);
            $table->dropColumn('kategori_defect_snapshot');
        });

        Schema::table('qcontrol_pemeriksaan_part', function (Blueprint $table): void {
            $table->dropIndex(['pemeriksaan_harian_id']);
            $table->dropIndex(['part_id']);
            $table->dropForeign('qcontrol_pemeriksaan_part_material_snapshot_foreign');
            $table->dropColumn([
                'material_id_snapshot',
                'kategori_ng_snapshot',
            ]);
        });

        Schema::table('qcontrol_pemeriksaan_harian', function (Blueprint $table): void {
            $table->dropUnique(['idempotency_key']);
            $table->dropIndex(['tanggal_produksi']);
            $table->dropIndex(['line_produksi_id']);
            $table->index('idempotency_key');
            $table->unique(['tanggal_produksi', 'line_produksi_id']);
            $table->dropColumn([
                'nomor_dokumen_snapshot',
                'revisi_dokumen_snapshot',
                'nama_pic_snapshot',
                'email_pic_snapshot',
                'disiapkan_oleh_snapshot',
                'diperiksa_oleh_snapshot',
                'disetujui_oleh_snapshot',
            ]);
        });
    }
};
