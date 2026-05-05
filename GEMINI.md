# GEMINI - PGNServer Context

Repo: PGNServer (Backend Utama).
Stack: Laravel 13, PHP 8.3+, PostgreSQL 17, Sanctum, Docker Compose.
Role: HeadQC (Satu-satunya role).
Bahasa: Bahasa Indonesia (Kode, Komentar, Log).
Fase: 2D-R3 Rekonsiliasi PRESS dan SEWING.

## Instruksi Penting:
- Jangan membuat role Admin/Inspector/Leader/QA.
- Jangan buat fitur import Excel otomatis tanpa instruksi fase.
- Bahasa Indonesia wajib untuk semua identitas kode.
- Server adalah source of truth untuk master data.

## Verifikasi:
`docker compose exec laravel.test php artisan test --compact`
`docker compose exec laravel.test vendor/bin/pint --dirty --format agent`

## Output:
Selalu sertakan Patch Report Bahasa Indonesia.
