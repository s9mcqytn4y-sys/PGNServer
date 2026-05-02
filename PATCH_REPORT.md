# PATCH_REPORT

## 1. Ringkasan Perubahan

- Mengubah skeleton Laravel 13 standar menjadi fondasi REST API only untuk PGNServer.
- Menambahkan route versioning awal di `/api/v1`.
- Menambahkan endpoint `GET /api/v1/kesehatan` dengan controller, request, resource, service aplikasi, model domain, dan pemeriksa koneksi database.
- Menambahkan envelope respons JSON terpusat melalui `App\Support\Api\ResponApi`.
- Menambahkan kode kesalahan API terpusat melalui `App\Support\Errors\KodeKesalahanApi`.
- Memaksa error route `/api/*` tetap JSON dari `bootstrap/app.php`.
- Menambahkan Laravel Sail dan menyiapkan `compose.yaml` berbasis `postgres:17-alpine`.
- Merapikan `.env.example`, `.gitignore`, `.editorconfig`, dan `composer.json`.
- Menghapus sisa frontend skeleton yang tidak relevan untuk backend API only.
- Menambahkan dokumentasi `README.md` dan `README_SETUP_LOKAL.md`.

## 2. Versi Stack Terdeteksi

- Laravel: `13.7.0`
- PHP: `8.5.4`
- Composer: `2.9.5`
- Laravel Boost: `2.4.6`
- Laravel Sail: `1.58.0`
- Laravel Sanctum: `4.3.1`
- PostgreSQL image: `postgres:17-alpine`

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

Dengan Sail:

```bash
cp .env.example .env
composer install
php artisan key:generate
./vendor/bin/sail up -d
./vendor/bin/sail artisan migrate
./vendor/bin/sail artisan test
curl http://localhost:8000/api/v1/kesehatan
```

Fallback bila wrapper Sail gagal di Windows:

```bash
docker compose up -d
docker compose exec laravel.test php artisan migrate
docker compose exec laravel.test php artisan test
```

## 10. Cara Menjalankan Test

Host:

```bash
php artisan test --compact
vendor/bin/pint --dirty --format agent
```

Container:

```bash
./vendor/bin/sail artisan test
```

## 11. Command Yang Sudah Dieksekusi

Dokumentasi resmi dan inspeksi:

```text
php artisan list --raw
php artisan help install:api
php artisan help sail:install
composer show laravel/sail --all
```

Setup:

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

Verifikasi HTTP lokal:

```text
php artisan serve --host=127.0.0.1 --port=8001
curl -H "Accept: application/json" http://127.0.0.1:8001/api/v1/kesehatan
```

Git dan GitHub:

```text
git init
git status --short --branch
git check-ignore .env .env.example vendor node_modules
gh auth status
gh repo view s9mcqytn4y-sys/PGNServer
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
- `curl http://127.0.0.1:8001/api/v1/kesehatan`: berhasil mengembalikan JSON `503` dengan kode `DATABASE_TIDAK_TERHUBUNG`

## 13. Status Verifikasi Sail / Docker

Percobaan command yang diminta:

```text
./vendor/bin/sail up -d
docker compose up -d
docker compose ps
```

Status:

- `./vendor/bin/sail up -d`: gagal karena environment Windows ini belum memiliki distribusi WSL terpasang.
- `docker compose up -d`: gagal karena Docker Desktop Linux engine tidak aktif, pipe `dockerDesktopLinuxEngine` tidak tersedia.
- Akibatnya command container berikut belum bisa diverifikasi secara nyata di mesin ini:
  - `./vendor/bin/sail artisan route:list`
  - `./vendor/bin/sail artisan test`
  - `./vendor/bin/sail artisan migrate:status`
  - `./vendor/bin/sail artisan config:clear`
  - `./vendor/bin/sail artisan optimize:clear`

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
- `git commit -m "Inisialisasi Laravel REST API PGNServer"`: berhasil
- hash commit awal: `3ff7499`
- `git branch -M main`: berhasil
- `git push -u origin main`: berhasil
- remote aktif: `https://github.com/s9mcqytn4y-sys/PGNServer.git`
- branch tracking: `main -> origin/main`

## 16. Risiko Dan Catatan Lanjutan

- Environment lokal ini memiliki warning PHP global bahwa beberapa extension dimuat ganda: `openssl`, `mbstring`, `pdo_mysql`, `curl`, `fileinfo`. Warning ini tidak berasal dari patch aplikasi, tetapi sebaiknya dibersihkan dari konfigurasi PHP CLI lokal.
- Wrapper Sail di Windows membutuhkan WSL. Repo sudah siap, tetapi host ini belum memenuhi prasyarat tersebut.
- Docker Desktop Linux engine belum aktif, sehingga verifikasi container belum bisa dibuktikan di host ini.
- `node_modules` lama masih tertinggal sebagian di workspace karena ada file yang terkunci oleh Windows, tetapi folder itu sudah di-ignore dan tidak mempengaruhi hasil commit.

## 17. Checklist Definition Of Done

- [x] Laravel 13 project siap REST API only
- [x] Tidak ada Blade/frontend demo yang relevan di entry point aplikasi
- [x] API route `/api/v1/kesehatan` tersedia
- [x] Response health check JSON konsisten
- [x] Error API tidak bocor sebagai HTML untuk route `/api/*`
- [x] `compose.yaml` memakai PostgreSQL 17
- [x] `.env.example` lengkap untuk development lokal
- [x] `.gitignore` aman
- [x] `.editorconfig` dirapikan
- [x] README setup lokal tersedia
- [x] Test health check hijau di host
- [ ] Verifikasi Sail hijau
- [x] Repo berhasil dipush ke `main`

## 18. Catatan Selisih Dari Request

- Remote Git sementara memakai HTTPS agar selaras dengan autentikasi aktif `gh auth status`. Jika diperlukan mutlak SSH, remote bisa diganti ke `git@github.com:s9mcqytn4y-sys/PGNServer.git` setelah kunci SSH pada mesin ini dipastikan siap.
- Verifikasi container belum tuntas karena blocker environment host, bukan karena error kode aplikasi.
