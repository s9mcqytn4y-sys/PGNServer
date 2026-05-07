<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('qcontrol_pemeriksaan_produksi_tanpa_ng', function (Blueprint $table): void {
            $table->uuid('id')->primary();
            $table->foreignUuid('pemeriksaan_harian_id')
                ->constrained('qcontrol_pemeriksaan_harian')
                ->cascadeOnDelete();

            // Relasi ke part (nullable jika part dihapus di masa depan tapi snapshot tetap ada)
            $table->foreignUuid('part_id')
                ->nullable()
                ->constrained('qcontrol_part')
                ->nullOnDelete();

            // Snapshot data part saat input
            $table->string('uniq_no_part');
            $table->string('nomor_part_snapshot');
            $table->string('nama_part_snapshot');

            $table->integer('total_produksi')->default(0);
            $table->text('catatan')->nullable();

            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();

            // Constraint: satu part hanya muncul sekali dalam section tanpa NG per pemeriksaan harian
            $table->unique(['pemeriksaan_harian_id', 'uniq_no_part'], 'idx_pemeriksaan_harian_part_tanpa_ng');
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('qcontrol_pemeriksaan_produksi_tanpa_ng');
    }
};
