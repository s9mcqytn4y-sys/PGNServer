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
- `POST /api/v1/qcontrol/contoh`: Endpoint contoh untuk pengujian integrasi. Butuh token HeadQC.
- `GET /api/v1/qcontrol/master-data`: Mengambil data Master (Part, Jenis Defect, Material) sebagai Source of Truth.
- `POST /api/v1/qcontrol/pemeriksaan-harian`: Menerima dan menyimpan pemeriksaan harian QControl secara idempotent.
- `GET /api/v1/qcontrol/pemeriksaan-harian`: Membaca daftar pemeriksaan harian QControl.
- `GET /api/v1/qcontrol/pemeriksaan-harian/{id}`: Membaca detail pemeriksaan harian QControl dengan prioritas snapshot historis.
- `GET /api/v1/qcontrol/laporan-bulanan/recording-defect`: Membaca agregasi monthly recording defect langsung dari transaksi daily.

## Fase Aktif
**PGNServer Fase 2F-A - Schema Final Daily QC, Validasi Seeder, dan Fondasi Monthly Report.**

Fokus fase ini adalah mengunci schema final daily QC, memvalidasi template master data sesuai form PDF/Excel sumber, dan menyediakan read model monthly yang dibentuk dari agregasi transaksi daily.

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
docker compose exec laravel.test php artisan qcontrol:validasi-master-data
docker compose exec laravel.test php artisan test --compact
curl http://127.0.0.1:8000/api/v1/kesehatan
```

Kredensial development HeadQC:

```text
Email    : headqc@pgn.local
Password : HeadQC@12345
```

Catatan host lokal:

- Untuk browser biasa, `localhost:8000` tetap dapat dipakai.
- Untuk QControl Desktop dan uji manual API, prioritaskan `http://127.0.0.1:8000`.

Contoh kirim pemeriksaan harian dari `cmd.exe`:

```bash
curl.exe -i -X POST "http://127.0.0.1:8000/api/v1/qcontrol/pemeriksaan-harian" ^
  -H "Accept: application/json" ^
  -H "Authorization: Bearer TOKEN_HEADQC" ^
  -H "X-Idempotency-Key: pemeriksaan-harian-press-001" ^
  -H "Content-Type: application/json" ^
  --data-raw "{\"clientDraftId\":\"draft-press-001\",\"tanggalProduksi\":\"2026-05-05\",\"lineProduksiId\":\"UUID_LINE_PRESS\",\"nomorDokumen\":\"FM-QA-025\",\"revisi\":\"1\",\"catatan\":\"Pemeriksaan line PRESS\",\"daftarPart\":[{\"partId\":\"UUID_PART_CB9\",\"totalCheck\":124,\"daftarDefect\":[{\"relasiPartDefectId\":\"UUID_RELASI_CB9_A\",\"slotWaktuId\":\"UUID_SLOT_0800_1200\",\"jumlahDefect\":2}]}]}"
```

Header wajib untuk submit daily:

- `Authorization: Bearer <token>`
- `Accept: application/json`
- `Content-Type: application/json`
- `X-Idempotency-Key: <uuid/string stabil dari client>`

Contoh baca daftar pemeriksaan harian:

```bash
curl.exe -i -X GET "http://127.0.0.1:8000/api/v1/qcontrol/pemeriksaan-harian?tanggalProduksi=2026-05-05&limit=20" ^
  -H "Accept: application/json" ^
  -H "Authorization: Bearer TOKEN_HEADQC"
```

Contoh baca detail pemeriksaan harian:

```bash
curl.exe -i -X GET "http://127.0.0.1:8000/api/v1/qcontrol/pemeriksaan-harian/UUID_PEMERIKSAAN_HARIAN" ^
  -H "Accept: application/json" ^
  -H "Authorization: Bearer TOKEN_HEADQC"
```

Contoh baca monthly recording defect:

```bash
curl.exe -i -X GET "http://127.0.0.1:8000/api/v1/qcontrol/laporan-bulanan/recording-defect?bulan=5&tahun=2026&lineProduksiId=UUID_LINE_PRESS" ^
  -H "Accept: application/json" ^
  -H "Authorization: Bearer TOKEN_HEADQC"
```

Monthly report selalu dibentuk dari agregasi pemeriksaan harian. Jangan membuat input manual untuk data bulanan.

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
