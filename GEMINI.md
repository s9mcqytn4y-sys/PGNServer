# GEMINI - PGNServer Context

Aplikasi: PGNServer (Backend Utama).
Stack: Laravel 13, PostgreSQL 17, Sanctum.
Role: HeadQC.
Bahasa: Bahasa Indonesia.

**SUMBER INSTRUKSI**:
Lihat **AGENTS.md** untuk detail arsitektur, endpoint, dan batasan role. Server adalah source of truth untuk master data.

## Verifikasi:
`docker compose exec laravel.test php artisan test --compact`
