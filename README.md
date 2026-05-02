# PGNServer

Backend monolith Laravel 13 untuk RESTful API yang disiapkan untuk dikonsumsi oleh Compose Multiplatform dan Kotlin 2.x. Repo ini sengaja dimulai sebagai API-only: tanpa Blade UI, tanpa Vite frontend, tanpa Livewire, dan tanpa fitur bisnis.

## Fondasi Teknologi

- Laravel 13
- PHP 8.5
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

- `GET /api/v1/kesehatan`

Contoh respons sukses:

```json
{
  "berhasil": true,
  "pesan": "Server berjalan normal",
  "data": {
    "status": "sehat",
    "namaAplikasi": "PGNServer",
    "versiApi": "v1",
    "waktuServer": "2026-05-03T10:15:30+07:00",
    "zonaWaktu": "Asia/Jakarta",
    "koneksiDatabase": {
      "status": "terhubung",
      "driver": "pgsql"
    }
  },
  "metadata": null,
  "kesalahan": null
}
```

Contoh respons saat database belum tersedia:

```json
{
  "berhasil": false,
  "pesan": "Server berjalan, tetapi koneksi database belum tersedia",
  "data": {
    "status": "terganggu",
    "namaAplikasi": "PGNServer",
    "versiApi": "v1",
    "waktuServer": "2026-05-03T10:15:30+07:00",
    "zonaWaktu": "Asia/Jakarta",
    "koneksiDatabase": {
      "status": "tidakTerhubung",
      "driver": "pgsql"
    }
  },
  "metadata": null,
  "kesalahan": {
    "kode": "DATABASE_TIDAK_TERHUBUNG",
    "detail": []
  }
}
```

## Kontrak Kotlin

```kotlin
@Serializable
data class ResponApi<T>(
    val berhasil: Boolean,
    val pesan: String,
    val data: T? = null,
    val metadata: JsonObject? = null,
    val kesalahan: KesalahanApi? = null
)

@Serializable
data class KesalahanApi(
    val kode: String,
    val detail: List<DetailKesalahanApi> = emptyList()
)

@Serializable
data class DetailKesalahanApi(
    val field: String? = null,
    val pesan: String
)
```

## Menjalankan Lokal

Ringkasan cepat:

```bash
cp .env.example .env
composer install
php artisan key:generate
./vendor/bin/sail up -d
./vendor/bin/sail artisan migrate
./vendor/bin/sail artisan test
curl http://localhost:8000/api/v1/kesehatan
```

Catatan penting untuk Windows:

- Wrapper `./vendor/bin/sail` membutuhkan WSL terpasang.
- Jika WSL belum tersedia, gunakan `docker compose` langsung setelah Docker Desktop Linux engine aktif.
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
