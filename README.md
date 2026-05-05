# PGNServer

Backend monolith Laravel 13 untuk RESTful API yang disiapkan untuk dikonsumsi oleh Compose Multiplatform dan Kotlin 2.x. Repo ini sengaja dimulai sebagai API-only: tanpa Blade UI, tanpa Vite frontend, tanpa Livewire, dan tanpa fitur bisnis.

## Fondasi Teknologi

- Laravel 13
- PHP 8.5 untuk tooling host
- PHP 8.4 untuk runtime container lokal
- PostgreSQL 17
- Laravel Sanctum untuk fondasi API auth
- Laravel Sail untuk local container workflow
- Pest untuk testing
- Laravel Pint untuk formatting

## Tujuan Arsitektur

Struktur aplikasi mengikuti Clean Architecture pragmatis ala Laravel:

- `app/Domain`
- `app/Application`
- `app/Infrastructure`
- `app/Http/Controllers/Api/V1`
- `app/Http/Requests/Api/V1`
- `app/Http/Resources/Api/V1`
- `app/Support/Api`
- `app/Support/Errors`
- `routes/api.php`
- `tests/Feature`
- `tests/Unit`

Controller hanya menerima request, memanggil service aplikasi, lalu mengembalikan response JSON. Detail teknis database tidak bocor ke layer HTTP.

## Endpoint Saat Ini

### Umum
- `GET /api/v1/kesehatan`: Pemeriksaan kesehatan server dan database.

### Autentikasi (HeadQC Only)
- `POST /api/v1/login`: Masuk ke sistem dan mendapatkan token Sanctum.
- `POST /api/v1/logout`: Keluar dan mencabut token aktif.
- `GET /api/v1/profil-saya`: Mengambil data profil HeadQC yang sedang login.

### QControl Integrasi
- `POST /api/v1/qcontrol/contoh`: Endpoint contoh untuk pengujian integrasi.
- `GET /api/v1/qcontrol/master-data`: Mengambil data Master (Part, Jenis Defect, Material) sebagai Source of Truth.

## Fase Aktif
**PGNServer Fase 2E-A - Hardening Bootstrap HeadQC dan Runtime Lokal.**

Fokus fase ini adalah memastikan bootstrap pengguna HeadQC lokal tetap stabil walaupun container PostgreSQL memakai volume lama dan seeder tidak otomatis dijalankan saat `docker compose up`.

## Aturan Pengembangan
- **Role**: Hanya ada role **HeadQC**. Jangan membuat role lain.
- **Bahasa**: Wajib menggunakan **Bahasa Indonesia** untuk kode (variabel, fungsi, komentar) dan log.
- **Data**: Server adalah sumber kebenaran data (*Source of Truth*).

## Menjalankan Lokal
... (sisanya tetap)

Ringkasan cepat:

```bash
cp .env.example .env
composer install
php artisan key:generate
docker compose up -d
docker compose exec laravel.test php artisan migrate
docker compose exec laravel.test php artisan qcontrol:pastikan-headqc
docker compose exec laravel.test php artisan db:seed --class=MasterDataQControlSeeder
docker compose exec laravel.test php artisan test --compact
curl http://localhost:8000/api/v1/kesehatan
```

Catatan penting untuk Windows:

- Wrapper `./vendor/bin/sail` membutuhkan WSL terpasang.
- Pada host ini wrapper `./vendor/bin/sail` masih gagal karena default shell WSL mengarah ke `docker-desktop` yang tidak menyediakan `/bin/bash`.
- Gunakan `docker compose` langsung setelah Docker Desktop Linux engine aktif.
- PostgreSQL tidak dipublish ke host; koneksi aplikasi tetap memakai service internal `pgsql`.
- Dokumen detail ada di `README_SETUP_LOKAL.md`.

## Verifikasi Host

Command dasar yang relevan untuk host:

```bash
php artisan --version
php artisan route:list
php artisan config:clear
php artisan test --compact
vendor/bin/pint --dirty --format agent
```

## Catatan API

- Semua endpoint API berada di bawah prefix `/api/v1`.
- Response selalu JSON object.
- Error API untuk `/api/*` dipaksa tetap JSON melalui exception rendering terpusat di `bootstrap/app.php`.
- Pesan user-facing konsisten dalam Bahasa Indonesia.

## File Terkait Setup

- `README_SETUP_LOKAL.md`
- `PATCH_REPORT.md`
