# AGENTS - PGNServer

Instruksi terpusat untuk agen AI agar memahami konteks proyek PGNServer secara ringkas dan hemat token. File ini adalah **Source of Truth** untuk semua instruksi agent.

## 1. Identitas Repo
- **Nama**: PGNServer
- **Tipe**: Backend Monolith / RESTful API Server.
- **Tujuan**: Backend utama untuk software desktop QControl dan aplikasi internal lainnya.
- **Bahasa Proyek**: **Bahasa Indonesia penuh** untuk penamaan file, fungsi, variabel, komentar, dan log (kecuali istilah teknis standar).

## 2. Stack Utama
- **Framework**: Laravel 13
- **Bahasa**: PHP 8.3+
- **Database**: PostgreSQL 17
- **Auth**: Laravel Sanctum (Token-based)
- **Environment**: Docker Compose (Laravel Sail compatible)
- **Testing**: Pest PHP
- **Linting**: Laravel Pint

## 3. Arsitektur & Panduan Laravel
- **Pola**: Clean Architecture pragmatis.
- **Folder Utama**:
  - `app/Domain`: Logic bisnis murni, entity, value object.
  - `app/Application`: Service layer, orchestrator.
  - `app/Infrastructure`: Implementasi database, client API eksternal.
  - `app/Http/Controllers/Api/V1`: Layer entry point API.
  - `app/Http/Resources/Api/V1`: Transformasi JSON response.
- **Panduan Khusus**: Untuk instruksi Laravel Boost dan konvensi framework yang mendalam, lihat [docs/LARAVEL_BOOST_GUIDELINES.md](file:///c:/Software/PGNServer/docs/LARAVEL_BOOST_GUIDELINES.md).

## 4. Aturan Role & Auth
- **Role Tunggal**: **HeadQC**.
- **Larangan**: Jangan membuat role baru (Admin, QC Inspector, QC Leader, Viewer, atau QA Manager).
- **Auth**: HeadQC login via `/api/v1/login` untuk mendapatkan token Sanctum.

## 5. Endpoint Penting
- `GET /api/v1/kesehatan`: Cek status server & database.
- `POST /api/v1/login`: Autentikasi HeadQC.
- `POST /api/v1/logout`: Revoke token.
- `GET /api/v1/profil-saya`: Data user aktif.
- `POST /api/v1/qcontrol/contoh`: Endpoint dummy/contoh integrasi.
- `GET /api/v1/qcontrol/master-data`: Source of truth untuk Part, Defect, dan Material.

## 6. Batasan & Kebijakan
- **Bahasa**: Wajib Bahasa Indonesia untuk semua kode baru.
- **Data**: Server adalah **Source of Truth**. QControl Desktop hanya melakukan cache.
- **Fase Saat Ini**: **2D-R3** (Rekonsiliasi PRESS dan SEWING).
- **Larangan Keras**:
  - Jangan mengimpor Excel secara otomatis tanpa instruksi fase khusus.
  - Jangan mengizinkan transaksi harian sebelum kontrak/template disetujui.
  - Jangan menambah role di luar HeadQC.

## 7. Command Verifikasi
Gunakan command berikut di lingkungan Docker:
```bash
docker compose exec laravel.test vendor/bin/pint --dirty --format agent
docker compose exec laravel.test php artisan test --compact
docker compose exec laravel.test php artisan route:list --path=api
```

## 8. Format Patch Report Wajib
Setiap perubahan wajib dilaporkan dengan format:
```text
PATCH REPORT - PGNServer - AgentOps-R1

1. Ringkasan keputusan teknis
2. File dibuat
3. File diubah
4. Koreksi penting yang dilakukan
5. Command yang dijalankan
6. Hasil verifikasi
7. Risiko tersisa
8. Rekomendasi fase berikutnya
```
