# PATCH_REPORT

## 1. Ringkasan Perubahan

- Merapikan fondasi Laravel 13 menjadi backend REST API only untuk PGNServer.
- Menjaga struktur Clean Architecture pragmatis dengan pemisahan `Domain`, `Application`, `Infrastructure`, dan layer HTTP API.
- Menambahkan endpoint `GET /api/v1/kesehatan` dengan service aplikasi, controller, request, resource, dan pemeriksa koneksi database ringan.
- Menstandarkan envelope respons JSON melalui `App\Support\Api\ResponApi`.
- Menstandarkan kode kesalahan API melalui `App\Support\Errors\KodeKesalahanApi`.
- Memaksa route `/api/*` selalu mengembalikan JSON aman melalui konfigurasi exception di `bootstrap/app.php`.
- Membersihkan artefak frontend skeleton Laravel yang tidak relevan untuk proyek API only.
- Mengganti runtime container lokal ke image PHP resmi agar verifikasi Docker lokal stabil di mesin ini.
- Menyiapkan `compose.yaml` dengan PostgreSQL 17, volume persisten, healthcheck database, dan port aplikasi `8000`.
- Merapikan `README.md`, `README_SETUP_LOKAL.md`, `.env.example`, `.gitignore`, dan `.editorconfig`.

## 2. Versi Stack Terdeteksi

- Laravel host: `13.7.0`
- Laravel container: `13.7.0`
- PHP host: `8.5.4`
- PHP container: `8.4.20`
- Composer host: `2.9.5`
- Composer container: `2.9.7`
- Laravel Boost: `2.4.6`
- Laravel Sail: `1.58.0`
- Laravel Sanctum: `4.3.1`
- PostgreSQL image: `postgres:17-alpine`
- Runtime image aplikasi lokal: `php:8.4-cli-bookworm`

## 3. File Yang Dibuat

- `app/Application/Kesehatan/MembacaStatusKesehatanServer.php`
- `app/Domain/Kesehatan/StatusKesehatanServer.php`
- `app/Domain/Kesehatan/StatusKoneksiDatabase.php`
- `app/Http/Controllers/Api/V1/PemeriksaanKesehatanController.php`
- `app/Http/Requests/Api/V1/MembacaStatusKesehatanRequest.php`
- `app/Http/Resources/Api/V1/StatusKesehatanResource.php`
- `app/Infrastructure/Kesehatan/PemeriksaKoneksiDatabase.php`
- `app/Support/Api/ResponApi.php`
- `app/Support/Errors/KodeKesalahanApi.php`
- `docker/laravel/Dockerfile`
- `README_SETUP_LOKAL.md`
- `PATCH_REPORT.md`

## 4. File Yang Diubah

- `.editorconfig`
- `.env.example`
- `.gitignore`
- `README.md`
- `app/Models/User.php`
- `app/Providers/AppServiceProvider.php`
- `bootstrap/app.php`
- `compose.yaml`
- `composer.json`
- `config/app.php`
- `config/database.php`
- `routes/api.php`
- `routes/web.php`
- `tests/Feature/ExampleTest.php`
- `tests/Feature/KesehatanApiTest.php`
- `tests/Pest.php`
- `tests/Unit/ExampleTest.php`
- `tests/Unit/MembacaStatusKesehatanServerTest.php`

## 5. File Yang Dihapus Saat Cleanup

- `.npmrc`
- `package-lock.json`
- `package.json`
- `public/build/assets/app-34mOoJaZ.js`
- `public/build/assets/app-D05hqmaw.css`
- `public/build/manifest.json`
- `resources/css/app.css`
- `resources/js/app.js`
- `resources/views/welcome.blade.php`
- `vite.config.js`

## 6. Struktur Folder Final

```text
app/
  Application/
    Kesehatan/
  Domain/
    Kesehatan/
  Http/
    Controllers/Api/V1/
    Requests/Api/V1/
    Resources/Api/V1/
  Infrastructure/
    Kesehatan/
  Support/
    Api/
    Errors/
docker/
  laravel/
routes/
  api.php
  web.php
tests/
  Feature/
  Unit/
```

## 7. Endpoint Yang Tersedia

- `GET /`
- `GET /api/v1/kesehatan`
- `GET /up`

Catatan:

- `POST /_boost/browser-logs` muncul di environment lokal karena `laravel/boost` adalah dev dependency.
- `GET /sanctum/csrf-cookie` tersedia dari Sanctum.
- Route storage lokal bawaan framework masih terdaftar di environment ini.

## 8. Format JSON Response Final

Respons sukses:

```json
{
  "berhasil": true,
  "pesan": "Server berjalan normal",
  "data": {},
  "metadata": null,
  "kesalahan": null
}
```

Respons gagal:

```json
{
  "berhasil": false,
  "pesan": "Permintaan tidak dapat diproses",
  "data": null,
  "metadata": null,
  "kesalahan": {
    "kode": "VALIDASI_GAGAL",
    "detail": [
      {
        "field": "namaField",
        "pesan": "Pesan validasi dalam Bahasa Indonesia"
      }
    ]
  }
}
```

## 9. Cara Menjalankan Lokal

Jalur standar sesuai permintaan:

```bash
cp .env.example .env
composer install
php artisan key:generate
./vendor/bin/sail up -d
./vendor/bin/sail artisan migrate
./vendor/bin/sail artisan test
curl http://localhost:8000/api/v1/kesehatan
```

Fallback yang terverifikasi di mesin ini:

```bash
docker compose up -d --force-recreate
docker compose exec -T laravel.test php artisan migrate
docker compose exec -T laravel.test php artisan test --compact
docker compose exec -T laravel.test php artisan route:list
curl http://127.0.0.1:8000/api/v1/kesehatan
```

## 10. Cara Menjalankan Test

Host:

```bash
php artisan test --compact
vendor/bin/pint --dirty --format agent
```

Container:

```bash
docker compose exec -T laravel.test php artisan test --compact
docker compose exec -T laravel.test vendor/bin/pint --dirty --format agent
```

## 11. Command Yang Dipakai

Dokumentasi dan inspeksi:

```text
php artisan list --raw
php artisan help install:api
php artisan help sail:install
composer show laravel/sail --all
```

Setup awal:

```text
composer require laravel/sail --dev
php artisan install:api --without-migration-prompt --no-interaction
composer dump-autoload --no-scripts
php artisan sail:install --with=pgsql --php=8.5 --no-interaction
```

Verifikasi host:

```text
composer update --lock --no-scripts
composer dump-autoload --no-scripts
php artisan --version
php artisan route:list
php artisan config:clear
php artisan test --compact
vendor/bin/pint --dirty --format agent
php artisan optimize:clear
```

Verifikasi Docker Compose:

```text
docker --version
docker context show
docker info --format "{{.ServerVersion}}"
wsl.exe --status
docker compose up -d --force-recreate
docker compose ps -a
docker compose exec -T laravel.test php artisan route:list
docker compose exec -T laravel.test php artisan migrate --force
docker compose exec -T laravel.test php artisan migrate:status --database=pgsql
docker compose exec -T laravel.test php artisan test --compact
docker compose exec -T laravel.test php artisan config:clear
docker compose exec -T laravel.test php artisan optimize:clear
docker compose exec -T laravel.test vendor/bin/pint --dirty --format agent
docker compose exec -T laravel.test php artisan about
curl.exe -i -H "Accept: application/json" http://127.0.0.1:8000/api/v1/kesehatan
```

Git dan GitHub:

```text
git init
git status --short --branch
git check-ignore .env .env.example vendor node_modules
gh auth status
gh repo view s9mcqytn4y-sys/PGNServer
gh repo set-default s9mcqytn4y-sys/PGNServer
gh repo edit s9mcqytn4y-sys/PGNServer --description "Laravel RESTful API backend untuk PGNServer dengan PostgreSQL dan Docker local environment"
git add .
git commit -m "Inisialisasi Laravel REST API PGNServer"
git branch -M main
git push -u origin main
```

## 12. Hasil Nyata Command Penting

- `php artisan --version`: berhasil, hasil `Laravel Framework 13.7.0`
- `php artisan route:list`: berhasil, route `GET /api/v1/kesehatan` terdaftar
- `php artisan config:clear`: berhasil
- `php artisan test --compact`: berhasil, `5` test lulus, `44` assertion
- `vendor/bin/pint --dirty --format agent`: berhasil
- `php artisan optimize:clear`: berhasil
- `docker compose ps -a`: berhasil, `laravel.test` dan `pgsql` berstatus `Up`, database `healthy`
- `docker compose exec -T laravel.test php artisan route:list`: berhasil, `GET /api/v1/kesehatan` terdaftar
- `docker compose exec -T laravel.test php artisan migrate:status --database=pgsql`: berhasil, tiga migrasi awal berstatus `Ran`
- `docker compose exec -T laravel.test php artisan test --compact`: berhasil, `5` test lulus, `44` assertion
- `docker compose exec -T laravel.test php artisan config:clear`: berhasil
- `docker compose exec -T laravel.test php artisan optimize:clear`: berhasil
- `docker compose exec -T laravel.test vendor/bin/pint --dirty --format agent`: berhasil
- `curl.exe -i -H "Accept: application/json" http://127.0.0.1:8000/api/v1/kesehatan`: berhasil, HTTP `200`, `Content-Type: application/json`, database `terhubung`

## 13. Status Verifikasi Sail / Docker

Status saat ini:

- `docker compose` berhasil dipakai end-to-end untuk menjalankan aplikasi dan PostgreSQL 17.
- Service `pgsql` memakai image `postgres:17-alpine` dengan healthcheck aktif.
- Port aplikasi diekspos ke host pada `http://127.0.0.1:8000`.
- PostgreSQL tidak dipublish ke host. Akses database untuk aplikasi tetap melalui service internal `pgsql`.

Catatan kompatibilitas:

- Wrapper `./vendor/bin/sail` di host Windows ini masih gagal karena WSL default shell mengarah ke `docker-desktop` yang tidak menyediakan `/bin/bash`.
- Karena itu, jalur operasional lokal yang benar-benar terverifikasi pada mesin ini adalah `docker compose`, bukan wrapper Sail.
- Ini bukan blocker untuk kode aplikasi, tetapi tetap perlu dicatat agar setup tim tidak membingungkan.

## 14. Status GH CLI

- `gh auth status`: berhasil
- akun aktif: `s9mcqytn4y-sys`
- protocol git dari GH CLI: `https`
- repo target `s9mcqytn4y-sys/PGNServer` dapat dibaca
- deskripsi repo berhasil diperbarui
- `gh repo set-default s9mcqytn4y-sys/PGNServer`: berhasil

## 15. Status Git Dan Push

- `git init`: berhasil
- `git check-ignore .env .env.example vendor node_modules`: berhasil, `.env`, `vendor`, dan `node_modules` ter-ignore, `.env.example` tidak
- commit awal dan update lanjutan sudah berhasil dipush ke `main`
- remote aktif: `https://github.com/s9mcqytn4y-sys/PGNServer.git`
- branch tracking: `main -> origin/main`

## 16. Risiko Dan Catatan Lanjutan

- PHP CLI host masih memunculkan warning extension ganda seperti `openssl`, `mbstring`, `pdo_mysql`, `curl`, dan `fileinfo`. Ini berasal dari konfigurasi PHP lokal host, bukan dari repo.
- Wrapper Sail di Windows belum usable di mesin ini selama WSL default shell tidak menyediakan `/bin/bash`.
- Runtime container lokal saat ini memakai image PHP resmi `8.4` agar build stabil. Sementara itu, tooling host masih berada di PHP `8.5`. Perbedaan minor ini aman untuk fondasi sekarang, tetapi sebaiknya diseragamkan saat image PHP resmi atau base runtime yang diinginkan sudah siap.
- Route dev dari package seperti Boost, Sanctum, dan storage masih muncul pada environment lokal. Tidak ada route bisnis yang tidak disengaja, tetapi ini tetap perlu dipahami saat membaca `route:list`.

## 17. Checklist Definition Of Done

- [x] Laravel 13 project siap REST API only
- [x] Tidak ada Blade atau frontend demo yang relevan di entry point aplikasi
- [x] API route `/api/v1/kesehatan` tersedia
- [x] Response health check JSON konsisten
- [x] Error API tidak bocor sebagai HTML untuk route `/api/*`
- [x] PostgreSQL 17 berjalan via Docker Compose
- [x] `.env.example` lengkap untuk development lokal
- [x] `.gitignore` aman
- [x] `.editorconfig` dirapikan
- [x] README setup lokal tersedia
- [x] Test health check hijau
- [x] Verifikasi container lokal hijau melalui `docker compose`
- [x] Repo berhasil dipush ke `main`
- [x] `PATCH_REPORT.md` lengkap

## 18. Catatan Selisih Dari Request

- Dokumentasi resmi dan struktur Laravel tetap diikuti, tetapi verifikasi container praktis pada host ini memakai `docker compose` langsung karena wrapper Sail gagal akibat konfigurasi WSL host.
- Port PostgreSQL sengaja tidak dipublish ke host untuk menghindari konflik lokal dan menjaga setup tetap minimal. Aplikasi tetap terkoneksi normal melalui network internal Compose.
- Request menyebut PostgreSQL 17 via Sail. Implementasi akhirnya tetap memakai artefak Sail untuk SQL init testing database, tetapi orkestrasi harian yang terverifikasi menggunakan `docker compose`.
