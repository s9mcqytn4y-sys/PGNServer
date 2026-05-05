<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('qcontrol_pemeriksaan_harian', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->date('tanggal_produksi');
            $table->foreignUuid('line_produksi_id')
                ->constrained('qcontrol_line_produksi')
                ->restrictOnDelete();
            $table->string('nomor_dokumen')->nullable()->default('FM-QA-025');
            $table->string('revisi')->nullable()->default('1');
            $table->foreignId('pengguna_headqc_id')
                ->constrained('users')
                ->restrictOnDelete();
            $table->string('client_draft_id')->nullable();
            $table->string('idempotency_key')->nullable()->index();
            $table->string('hash_payload')->nullable();
            $table->string('status')->default('DITERIMA');
            $table->integer('total_check')->default(0);
            $table->integer('total_ok')->default(0);
            $table->integer('total_defect')->default(0);
            $table->decimal('rasio_defect', 8, 2)->default(0);
            $table->text('catatan')->nullable();
            $table->timestamp('diterima_pada')->nullable();
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();

            $table->unique(['tanggal_produksi', 'line_produksi_id']);
        });

        Schema::create('qcontrol_pemeriksaan_part', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->foreignUuid('pemeriksaan_harian_id')
                ->constrained('qcontrol_pemeriksaan_harian')
                ->cascadeOnDelete();
            $table->foreignUuid('part_id')
                ->constrained('qcontrol_part')
                ->restrictOnDelete();
            $table->integer('total_check')->default(0);
            $table->integer('total_ok')->default(0);
            $table->integer('total_defect')->default(0);
            $table->decimal('rasio_defect', 8, 2)->default(0);
            $table->integer('urutan_tampil')->default(0);
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();

            $table->unique(['pemeriksaan_harian_id', 'part_id']);
        });

        Schema::create('qcontrol_pemeriksaan_defect_slot', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->foreignUuid('pemeriksaan_part_id')
                ->constrained('qcontrol_pemeriksaan_part')
                ->cascadeOnDelete();
            $table->foreignUuid('relasi_part_defect_id')
                ->constrained('qcontrol_part_jenis_defect')
                ->restrictOnDelete();
            $table->foreignUuid('jenis_defect_id')
                ->constrained('qcontrol_jenis_defect')
                ->restrictOnDelete();
            $table->foreignUuid('slot_waktu_id')
                ->constrained('qcontrol_slot_waktu')
                ->restrictOnDelete();
            $table->integer('jumlah_defect')->default(0);
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();

            $table->unique(['pemeriksaan_part_id', 'relasi_part_defect_id', 'slot_waktu_id']);
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('qcontrol_pemeriksaan_defect_slot');
        Schema::dropIfExists('qcontrol_pemeriksaan_part');
        Schema::dropIfExists('qcontrol_pemeriksaan_harian');
    }
};
