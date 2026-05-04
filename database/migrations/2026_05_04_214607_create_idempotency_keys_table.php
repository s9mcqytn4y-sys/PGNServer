<?php

declare(strict_types=1);

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('idempotency_keys', function (Blueprint $table) {
            $table->id();
            $table->string('kunci_idempotency')->unique();
            $table->string('metode_http');
            $table->string('endpoint');
            $table->string('hash_payload')->nullable();
            $table->string('status_pemrosesan');
            $table->integer('response_status_http')->nullable();
            $table->json('response_body')->nullable();
            $table->string('sumber_aplikasi')->nullable();
            $table->timestamp('diproses_pada')->nullable();
            $table->timestamp('dibuat_pada')->useCurrent();
            $table->timestamp('diperbarui_pada')->useCurrent();

            $table->index('endpoint');
            $table->index('status_pemrosesan');
            $table->index('dibuat_pada');
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('idempotency_keys');
    }
};
