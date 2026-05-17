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

## 3. Intelligence Logic: Auto-Trace NG & TPS Core Hardening
Gemini telah diprogram untuk memahami logika otomatisasi di PostgreSQL 17.x:
- **Input:** Operator menginput NG pada inspection_logs.
- **Process:** Trigger `after_inspection_log_insert` secara cerdas mencari material penyebab melalui tabel `bill_of_materials`.
- **Output:** Jejak cacat material otomatis tercatat di `material_defect_ledger`.
- **TPS Core Hardening:** Sistem secara ketat hanya menerima input shift "NORMAL", memvalidasi format dan ketersediaan Time Slots dari `GET /api/v1/operasi/lembar_periksa/options`, dan menuntut akurasi 100% persamaan `TotalProduksi == KuantitasOK + KuantitasNG` untuk seluruh detail QC yang diunggah.

## 4. 7 QC Tools Query Reference
Gemini dapat langsung memanggil rute API terbuka:
- `GET /api/v1/analitik/metrik_pareto_bulanan` (Menghasilkan proporsi 80/20 dari cacat material berbasis Window Function SQL murni).

Dokumentasi lengkap dan antarmuka uji coba (kontrak API) tersedia secara dinamis di: `http://localhost:8080/swagger/index.html`

## 5. Security & Authentication
Akses API menggunakan standar RESTful modern (Go Gin). Pastikan setiap request menyertakan Bearer Token. Akun default Leader:
- **NIP:** 2211019
- **Password:** admin

Klien QControl mengimplementasikan Ktor HTTP Response Interceptor untuk mendeteksi status `401 Unauthorized` atau kedaluwarsa token secara global, yang secara otomatis memicu alur keluar sesi (`auto-logout`) untuk melindungi integritas sesi pengguna.

Untuk menjamin transmisi aman pada local Outbox, klien memverifikasi liveness server melalui `/api/v1/kesehatan` (atau `/api/v1/health`) serta kesiapan database melalui `/api/v1/readiness` sebelum mengirim batch request.

## 6. Manufacturing & Master Data API Endpoints (GORM CRUD)
Sistem kini mengekspos antarmuka CRUD lengkap yang aman terproteksi JWT untuk pengelolaan rantai pasok manufaktur:
- **Pemasok (Suppliers)**: `POST/GET/PUT/DELETE /api/v1/suppliers`
- **Bahan Baku (Materials)**: `POST/GET/PUT/DELETE /api/v1/materials`
- **Pelanggan (Customers)**: `POST/GET/PUT/DELETE /api/v1/customers`
- **BOM (Bill of Materials)**: `POST/GET/PUT/DELETE /api/v1/boms`

Setiap mutasi data master secara otomatis menghapus cache memori global (`cache.GlobalCache.Clear()`) untuk menyajikan metrik telemetri yang selalu mutakhir pada Landing Page dinamis.

## 7. Dual-Envelope Response Adaptation
Untuk memfasilitasi kelancaran integrasi dengan klien QControl Kotlin Multiplatform (KMP), PGNServer mengimplementasikan struktur amplop JSON ganda (dual-envelope). Respons dikembalikan dengan properti bahasa Indonesia (`sukses`, `pesan`, `data`, `metadata`, `kesalahan`) dan properti legacy (`success`, `message`, `status`, `meta`, `errors`) secara bersamaan. Format ini diuji dengan baik untuk mencegah kegagalan serialisasi di sisi client.

## 8. QControl Visual Canvas Integration & Sync
Klien desktop QControl merender data secara offline-first menggunakan Jetpack Compose Native Canvas:
- **Pareto Chart**: Menampilkan proporsi 80/20 dari defect master dengan dynamic threshold markers dan dynamic cumulative line drawing.
- **Defect Histogram**: Menyajikan visualisasi frekuensi defect slot per jam kerja secara real-time.
- **SQLite Engine**: Seluruh metrik visualisasi tersebut dihitung dan diambil langsung dari database SQLite lokal secara sinkron, yang ditarik dari endpoint analitik saat jaringan tersedia.

## 9. Development Optimization & Git Workflow Hardening
Setiap interaksi dengan sistem wajib memperhatikan panduan optimasi token dan keamanan integrasi:
1. **Token Savings Analytics (RTK)**: Gunakan CLI proxy `rtk` untuk seluruh perintah Git, kompilasi, dan testing guna mereduksi token overhead (hingga 90%). Contoh: `rtk git status`, `rtk go test ./...`.
2. **Environmental Safety Rules**: Seluruh variabel rahasia (`JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, dan `DATABASE_URL`) wajib dimuat secara dinamis via berkas `.env` dan tidak boleh di-hardcode di dalam kode sumber Go.
3. **Mandat Satu Cabang (No-Branching)**: Jangan membuat branch baru. Seluruh progress wajib di-commit dan langsung di-push ke branch `main` repositori `PGNServer`.

