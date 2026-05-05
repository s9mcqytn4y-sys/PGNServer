# CODEX - PGNServer Context

Identitas: PGNServer (Backend REST API Laravel 13).
Stack: PHP 8.3+, PostgreSQL 17, Sanctum, Docker.
Role: HeadQC (Tunggal).
Bahasa: Bahasa Indonesia penuh untuk semua elemen kode.
Fase: 2D-R3 (Rekonsiliasi PRESS & SEWING).

## Aturan Utama:
1. Semua penamaan (variabel, fungsi, file) wajib Bahasa Indonesia.
2. HeadQC adalah satu-satunya role.
3. Master data di server adalah source of truth.
4. Verifikasi dengan:
   - `docker compose exec laravel.test vendor/bin/pint --dirty --format agent`
   - `docker compose exec laravel.test php artisan test --compact`

## Patch Report:
Gunakan format "PATCH REPORT - PGNServer - AgentOps" di setiap akhir tugas.
