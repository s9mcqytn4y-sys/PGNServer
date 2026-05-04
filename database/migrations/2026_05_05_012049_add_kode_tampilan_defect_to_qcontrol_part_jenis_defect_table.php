<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::table('qcontrol_part_jenis_defect', function (Blueprint $table): void {
            $table->string('kode_tampilan_defect')->nullable()->after('jenis_defect_id');
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('qcontrol_part_jenis_defect', function (Blueprint $table): void {
            $table->dropColumn('kode_tampilan_defect');
        });
    }
};
