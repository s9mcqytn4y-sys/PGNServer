# CLAUDE - PGNServer Context

Identitas: PGNServer (Laravel 13 REST API).
Role: HeadQC (Solo Developer, branch main).
Stack: PHP 8.3+, PostgreSQL 17, Sanctum, Docker.
Bahasa: Bahasa Indonesia Penuh.
Fase Aktif: 2D-R3 (Rekonsiliasi PRESS dan SEWING).

## Batasan:
- Tidak ada role Admin/Inspector/Viewer/QA Manager.
- Tidak ada import Excel otomatis sebelum fase khusus.
- Tidak ada transaksi harian sebelum kontrak disetujui.
- Kode dan log wajib Bahasa Indonesia.

## Command Verifikasi:
`docker compose exec laravel.test vendor/bin/pint --dirty --format agent`
`docker compose exec laravel.test php artisan test --compact`

## Report:
Wajib menggunakan format "PATCH REPORT - PGNServer - AgentOps".
