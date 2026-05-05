# Setup Lokal PGNServer

Dokumen ini fokus pada setup development lokal untuk backend REST API PGNServer.

## Fase Aktif

**PGNServer Fase 2E-B - Kontrak dan Penyimpanan Pemeriksaan Harian QControl**

## Prasyarat

- PHP 8.5
- Composer 2.x
- Docker Desktop
- Runtime container lokal memakai PHP 8.4
- Untuk Windows dengan Sail wrapper:
  - WSL terpasang
  - minimal satu distribusi Linux aktif

## Variabel Penting

`.env.example` sudah disiapkan dengan nilai development:

```dotenv
APP_NAME=PGNServer
APP_ENV=local
APP_DEBUG=true
APP_URL=http://localhost:8000
APP_TIMEZONE=Asia/Jakarta
APP_PORT=8000
WWWUSER=1000
WWWGROUP=1000

DB_CONNECTION=pgsql
DB_HOST=pgsql
DB_PORT=5432
DB_DATABASE=pgn_server
DB_USERNAME=pgn_server
DB_PASSWORD=password_lokal_dev
QCONTROL_HEADQC_EMAIL=headqc@pgn.local
QCONTROL_HEADQC_PASSWORD=HeadQC@12345
```

## Jalur Standar Dengan Sail

```bash
cp .env.example .env
composer install
php artisan key:generate
./vendor/bin/sail up -d
./vendor/bin/sail artisan migrate
./vendor/bin/sail artisan test
./vendor/bin/sail artisan route:list
./vendor/bin/sail artisan config:clear
./vendor/bin/sail artisan optimize:clear
./vendor/bin/sail artisan migrate:status
curl http://localhost:8000/api/v1/kesehatan
```

## Fallback Bila Wrapper Sail Tidak Bisa Dipakai di Windows

Jika `./vendor/bin/sail` gagal karena WSL belum ada, aktifkan Docker Desktop Linux engine lalu gunakan:

```bash
docker compose up -d
docker compose exec laravel.test php artisan migrate
docker compose exec laravel.test php artisan qcontrol:pastikan-headqc
docker compose exec laravel.test php artisan db:seed --class=MasterDataQControlSeeder
docker compose exec laravel.test php artisan test --compact
docker compose exec laravel.test php artisan route:list --path=api
docker compose exec laravel.test php artisan migrate:status
```

Workflow ini penting karena volume PostgreSQL bersifat persisten. Saat `docker compose up -d`, container database bisa tetap memakai data lama dan bootstrap HeadQC tidak dijalankan otomatis.

## Uji Login HeadQC

Setelah bootstrap runtime lokal selesai, uji login dari `cmd.exe`:

```bash
curl.exe -i -X POST "http://127.0.0.1:8000/api/v1/login" ^
  -H "Accept: application/json" ^
  -H "Content-Type: application/json" ^
  --data-raw "{\"email\":\"headqc@pgn.local\",\"password\":\"HeadQC@12345\"}"
```

## Uji Endpoint Pemeriksaan Harian

Setelah login berhasil dan token HeadQC tersedia, uji endpoint pemeriksaan harian dari `cmd.exe`:

```bash
curl.exe -i -X POST "http://127.0.0.1:8000/api/v1/qcontrol/pemeriksaan-harian" ^
  -H "Accept: application/json" ^
  -H "Authorization: Bearer TOKEN_HEADQC" ^
  -H "X-Idempotency-Key: pemeriksaan-harian-press-001" ^
  -H "Content-Type: application/json" ^
  --data-raw "{\"clientDraftId\":\"draft-press-001\",\"tanggalProduksi\":\"2026-05-05\",\"lineProduksiId\":\"UUID_LINE_PRESS\",\"nomorDokumen\":\"FM-QA-025\",\"revisi\":\"1\",\"catatan\":\"Pemeriksaan line PRESS\",\"daftarPart\":[{\"partId\":\"UUID_PART_CB9\",\"totalCheck\":124,\"daftarDefect\":[{\"relasiPartDefectId\":\"UUID_RELASI_CB9_A\",\"slotWaktuId\":\"UUID_SLOT_0800_1200\",\"jumlahDefect\":2}]}]}"
```

Nilai `UUID_LINE_PRESS`, `UUID_PART_CB9`, `UUID_RELASI_CB9_A`, dan `UUID_SLOT_0800_1200` diambil dari endpoint `GET /api/v1/qcontrol/master-data`.

## Endpoint Verifikasi

Health check:

```bash
curl http://localhost:8000/api/v1/kesehatan
```

Jika database belum aktif, respons `503` tetap JSON dan itu normal untuk fase setup awal.

## Catatan Docker

- Aplikasi diekspos di port `8000`
- PostgreSQL tidak diekspos ke host. Gunakan service `pgsql` melalui network internal Compose.
- Service database bernama `pgsql`
- Image database: `postgres:17-alpine`
- Volume database persisten: `pgnserver-pgsql-data`

## Troubleshooting Singkat

### `./vendor/bin/sail` gagal di Windows

Penyebab umum:

- WSL belum diinstal
- distribusi Linux belum diinstal
- default shell WSL tidak menyediakan `/bin/bash`

Periksa:

```bash
wsl.exe --list --online
```

Lalu instal distribusi yang diinginkan:

```bash
wsl.exe --install Ubuntu
```

### `docker compose up -d` gagal konek daemon

Pastikan Docker Desktop Linux engine aktif. Error yang umum adalah pipe `dockerDesktopLinuxEngine` tidak ditemukan.

### Health check mengembalikan `DATABASE_TIDAK_TERHUBUNG`

Periksa:

- container `pgsql` sudah hidup
- kredensial di `.env` sesuai
- migrasi sudah dijalankan
- jalankan `docker compose exec laravel.test php artisan qcontrol:pastikan-headqc` jika login HeadQC lokal gagal
