# ♊ Gemini AI & MCP Integration Guide

Panduan konfigurasi Gemini AI sebagai asisten teknis dan analitik untuk proyek **PGNServer**.

## 1. Environment Configuration
Berdasarkan file `.env` di direktori `C:\Software\PGNServer`:
- **DB_NAME:** `pgn_db`
- **DB_USER:** `admin`
- **DB_PASSWORD:** `admin`
- **Host:** `localhost` (Port: 5432)

## 2. Running Gemini as MCP Server
Gunakan perintah berikut di terminal (PowerShell/CMD) untuk menghubungkan Gemini dengan database Anda:

```bash
# Menjalankan Postgres MCP Server
npx -y @modelcontextprotocol/server-postgres "postgresql://admin:admin@localhost:5432/pgn_db"
```

## 3. Intelligence Logic: Auto-Trace NG
Gemini telah diprogram untuk memahami logika otomatisasi di PostgreSQL 17.x:
- **Input:** Operator menginput NG pada inspection_logs.
- **Process:** Trigger `after_inspection_log_insert` secara cerdas mencari material penyebab melalui tabel `bill_of_materials`.
- **Output:** Jejak cacat material otomatis tercatat di `material_defect_ledger`.

## 4. 7 QC Tools Query Reference
Gemini dapat langsung memanggil rute API terbuka:
- `GET /api/v1/analitik/metrik_pareto_bulanan` (Menghasilkan proporsi 80/20 dari cacat material berbasis Window Function SQL murni).

Dokumentasi lengkap dan antarmuka uji coba (kontrak API) tersedia secara dinamis di: `http://localhost:8080/swagger/index.html`

## 5. Security & Authentication
Akses API menggunakan standar RESTful modern (Go Gin). Pastikan setiap request menyertakan Bearer Token. Akun default Leader:
- **NIP:** 2211019
- **Password:** admin

Environment Path: `C:\Software\PGNServer` | Docker Engine: Active
