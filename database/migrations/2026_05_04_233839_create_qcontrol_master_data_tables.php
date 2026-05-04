<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('qcontrol_line_produksi', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->string('kode_line')->unique();
            $table->string('nama_line');
            $table->boolean('aktif')->default(true);
            $table->integer('urutan_tampil')->default(0);
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();
        });

        Schema::create('qcontrol_slot_waktu', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->string('kode_slot')->unique();
            $table->string('label_slot');
            $table->time('jam_mulai')->nullable();
            $table->time('jam_selesai')->nullable();
            $table->integer('urutan_tampil')->default(0);
            $table->boolean('aktif')->default(true);
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();
        });

        Schema::create('qcontrol_material', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->string('kode_material')->nullable();
            $table->string('nama_material')->unique();
            $table->boolean('aktif')->default(true);
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();
        });

        Schema::create('qcontrol_kategori_defect', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->string('kode_kategori')->unique();
            $table->string('nama_kategori');
            $table->boolean('aktif')->default(true);
            $table->integer('urutan_tampil')->default(0);
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();
        });

        Schema::create('qcontrol_jenis_defect', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->string('kode_defect')->unique();
            $table->string('nama_defect');
            $table->foreignUuid('kategori_defect_id')
                ->nullable()
                ->constrained('qcontrol_kategori_defect')
                ->nullOnDelete();
            $table->boolean('aktif')->default(true);
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();
        });

        Schema::create('qcontrol_part', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->string('kode_unik_part')->unique();
            $table->string('nama_part');
            $table->string('nomor_part')->nullable();
            $table->foreignUuid('material_id')
                ->nullable()
                ->constrained('qcontrol_material')
                ->nullOnDelete();
            $table->string('kode_proyek')->nullable();
            $table->integer('jumlah_item_per_kanban')->nullable();
            $table->foreignUuid('line_default_id')
                ->nullable()
                ->constrained('qcontrol_line_produksi')
                ->nullOnDelete();
            $table->boolean('aktif')->default(true);
            $table->string('sumber_data')->nullable();
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();
        });

        Schema::create('qcontrol_part_jenis_defect', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->foreignUuid('part_id')
                ->constrained('qcontrol_part')
                ->cascadeOnDelete();
            $table->foreignUuid('jenis_defect_id')
                ->constrained('qcontrol_jenis_defect')
                ->cascadeOnDelete();
            $table->integer('urutan_tampil')->default(0);
            $table->boolean('aktif')->default(true);
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();

            $table->unique(['part_id', 'jenis_defect_id']);
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('qcontrol_part_jenis_defect');
        Schema::dropIfExists('qcontrol_part');
        Schema::dropIfExists('qcontrol_jenis_defect');
        Schema::dropIfExists('qcontrol_kategori_defect');
        Schema::dropIfExists('qcontrol_material');
        Schema::dropIfExists('qcontrol_slot_waktu');
        Schema::dropIfExists('qcontrol_line_produksi');
    }
};
