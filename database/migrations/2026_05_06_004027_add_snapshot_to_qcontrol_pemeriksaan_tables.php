<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::table('qcontrol_pemeriksaan_harian', function (Blueprint $table): void {
            $table->string('kode_line_snapshot')->nullable()->after('line_produksi_id');
            $table->string('nama_line_snapshot')->nullable()->after('kode_line_snapshot');
        });

        Schema::table('qcontrol_pemeriksaan_part', function (Blueprint $table): void {
            $table->string('kode_unik_part_snapshot')->nullable()->after('part_id');
            $table->string('nomor_part_snapshot')->nullable()->after('kode_unik_part_snapshot');
            $table->string('nama_part_snapshot')->nullable()->after('nomor_part_snapshot');
            $table->string('nama_material_snapshot')->nullable()->after('nama_part_snapshot');
        });

        Schema::table('qcontrol_pemeriksaan_defect_slot', function (Blueprint $table): void {
            $table->string('kode_tampilan_defect_snapshot')->nullable()->after('relasi_part_defect_id');
            $table->string('kode_defect_snapshot')->nullable()->after('kode_tampilan_defect_snapshot');
            $table->string('nama_defect_snapshot')->nullable()->after('kode_defect_snapshot');
            $table->string('kode_slot_snapshot')->nullable()->after('slot_waktu_id');
            $table->string('label_slot_snapshot')->nullable()->after('kode_slot_snapshot');
            $table->string('jam_mulai_snapshot')->nullable()->after('label_slot_snapshot');
            $table->string('jam_selesai_snapshot')->nullable()->after('jam_mulai_snapshot');
        });
    }

    public function down(): void
    {
        Schema::table('qcontrol_pemeriksaan_defect_slot', function (Blueprint $table): void {
            $table->dropColumn([
                'kode_tampilan_defect_snapshot',
                'kode_defect_snapshot',
                'nama_defect_snapshot',
                'kode_slot_snapshot',
                'label_slot_snapshot',
                'jam_mulai_snapshot',
                'jam_selesai_snapshot',
            ]);
        });

        Schema::table('qcontrol_pemeriksaan_part', function (Blueprint $table): void {
            $table->dropColumn([
                'kode_unik_part_snapshot',
                'nomor_part_snapshot',
                'nama_part_snapshot',
                'nama_material_snapshot',
            ]);
        });

        Schema::table('qcontrol_pemeriksaan_harian', function (Blueprint $table): void {
            $table->dropColumn([
                'kode_line_snapshot',
                'nama_line_snapshot',
            ]);
        });
    }
};
