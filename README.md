# PGNServer - Manufacturing Traceability & Advanced Analytics Engine 🚀

[![CI/CD Status](https://github.com/bonecom-group/pgn-server/actions/workflows/ci.yml/badge.svg)](https://github.com/bonecom-group/pgn-server/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25.x-blue.svg)](https://go.dev)
[![Postgres Version](https://img.shields.io/badge/PostgreSQL-17.x-blue.svg)](https://www.postgresql.org)
[![TPS Standard](https://img.shields.io/badge/Toyota_Production_System-AQL_Compliant-green.svg)](#)

Introduce **PGNServer**, *literally* the ultimate *modular monolith backend* yang dirancang *highly robust* untuk mengawal tata niaga manufaktur dan inspeksi kualitas (*Quality Control*) berstandar Toyota Production System (TPS). Di-build dengan **Go 1.25.x** (with `GOEXPERIMENT=jsonv2` enabled) dan **PostgreSQL 17.x**, sistem ini memberikan *performance benchmark* kelas wahid dengan *strict transaction integrity* demi *zero-defect tolerance*.

---

## 🏗️ Monolithic Modular Architecture

Aplikasi ini mengadopsi pola *Clean Architecture* dan *Domain-Driven Design (DDD)* untuk mengisolasi fungsionalitas bisnis secara modular, *which is* sangat memudahkan tim untuk *scaling up* fitur tanpa menimbulkan *dependency hell*.

```mermaid
graph TD
    subgraph API Gateway
        A[cmd/api/main.go] --> B[infrastruktur.Middleware]
    end
    subgraph Modular Domain
        B --> C[otentikasi.Module]
        B --> D[manufaktur.Module]
        B --> E[kualitas.Module]
        B --> F[analitik.Module]
        B --> G[media.Module]
    end
    subgraph Data Infrastructure
        C & D & E & F & G --> H[(PostgreSQL 17.x)]
    end
```

### Module Breakdown:
*   **`cmd/api/`**: The absolute entrypoint. Mengatur inisialisasi koneksi DB, dependency injection, routing, serta *graceful logging*.
*   **`internal/otentikasi/`**: Keamanan tingkat tinggi menggunakan **JWT v5** dan enkripsi *bcrypt*. Menyediakan proteksi *Role-Based Access Control (RBAC)* ketat khusus role **Leader QC**.
*   **`internal/manufaktur/`**: Mengelola *Master Data* utama seperti `customers`, `suppliers` (Pemasok), `MATERIAL` (Material & Finished Good), dan `bill_of_materials` (BOM).
*   **`internal/kualitas/`**: Inti operasional kontrol kualitas (*Quality Control*). Merekam checksheet secara atomik dengan *zero-loop* memanfaatkan fungsi native PostgreSQL `jsonb_to_recordset`.
*   **`internal/analitik/`**: Engine analitik premium 7 QC Tools. Menghitung Pareto Defect secara dinamis menggunakan *PostgreSQL 17 Window Functions* serta melacak akar masalah cacat material (*BOM Tracing*).
*   **`internal/media/`**: Manajemen berkas gambar dengan perlindungan *Dual-Default Fallback* yang *literally* aman dari aksi penghapusan tidak disengaja.

---

## 🚀 Key Value Propositions (Feature Unggulan)

### 📊 1. PostgreSQL 17 Window-Powered Pareto Analytics
Kami meniadakan looping array lambat di level aplikasi Go. Kalkulasi rasio cacat kumulatif dilakukan langsung di pangkalan data melalui kueri analitik canggih menggunakan *SQL Window Function* `SUM(rasio_cacat) OVER (ORDER BY ...)`:
> [!TIP]
> Hal ini *literally* menghemat memori alokasi (0 allocations in Go) dan menjamin pengurutan Pareto 80/20 yang *lightning-fast* bahkan ketika mengolah jutaan baris data secara realtime.

### 🌳 2. Auto-Trace NG & Recursive BOM Tracing
Ketika inspektur menemukan produk NG (*No Good*), sistem secara rekursif menelusuri pohon Bill of Materials (BOM) untuk mendeteksi material pembentuk dan pemasok (*supplier*) asalnya.
*   **Circular Reference Shield**: Dilengkapi *visited map validation* untuk mendeteksi *circular dependency* (A → B → A) secara otomatis dan memutus putaran tak terbatas sebelum memicu *stack overflow*.
*   **Toyota-Style Internal Defect Resolver**: Jika cacat tergolong `PROCESS`, sistem secara pintar memetakannya sebagai *Internal Process Defect* disuplai oleh *Internal Production Line*.

### ⚡ 3. Single-Atomic Transmisi Lembar Periksa (Zero-Loop Batch Write)
Merekam ribuan baris checksheet cacat sekaligus dalam satu kali jalan menggunakan PostgreSQL native **`jsonb_to_recordset`**:
```sql
INSERT INTO detail_inspeksis (lembar_periksa_id, unik_part_id, kode_cacat, waktu_pergeseran, rasio_cacat, rasio_total_ok)
SELECT ?, "unikPartId", "kodeCacat", "waktuPergeseran", "rasioCacat", "rasioTotalOK"
FROM jsonb_to_recordset(?::jsonb) AS x(...)
```
> [!NOTE]
> Menghilangkan masalah laten N+1 queries. Seluruh detail inspeksi dimasukkan secara atomik dalam satu transaksi pangkalan data, *which is* dilindungi oleh `tx.Rollback()` jika terjadi kegagalan sekecil apa pun.

### 🛡️ 4. Dual-Default Media Protection
Modul media dilengkapi pertahanan ganda:
1.  **Dynamic PNG Builder**: Secara otonom mendeteksi dan menciptakan berkas 1x1 piksel dummy `part.png` dan `avatar.png` pada saat *startup* jika berkas asli absen di server.
2.  **Deletion Guard**: Mencegah penghapusan berkas cadangan default oleh siapa pun melalui API, sehingga sistem terbebas dari ancaman *blank image crash*.

---

## ⚙️ Quick Start & Production Setup

### 1. Kebutuhan Environment (`.env`)
Salin atau buat berkas `.env` di direktori utama proyek dengan konfigurasi berikut:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin
DB_NAME=pgn_db
APP_PORT=8080
JWT_SECRET=pgn-rahasia-korporat-2026
GIN_MODE=debug
```

### 2. Kompilasi & Menjalankan Kontainer Docker (Production-Ready)
Kami menyediakan multi-stage `Dockerfile` berbasis *Alpine Linux* yang di-hardening secara ketat dengan *Non-Root User* (`pgnuser`) untuk meminimalkan *attack surface*.

Volume database kini dipetakan secara fisik ke `./docker/pgdata` untuk menjamin ketahanan persisten jika kontainer direstart atau dimatikan tanpa sengaja via `docker compose down -v`.

Jalankan seluruh infrastruktur (Database Postgres 17 & API Server) hanya dengan satu perintah:
```bash
docker-compose up -d --build
```

### 3. Eksekusi Pengujian Otomatis (Testing Pyramid)
Untuk menjamin integritas kode saat *refactoring*, jalankan seluruh suite pengujian integrasi (termasuk validasi database, window function, circular tracing, dan RBAC):
```bash
go test -v -race ./...
```

---

## 🗄️ Disaster Recovery & Manual Database Backup

Demi mencegah kehilangan data historis inspeksi kualitas akibat kesalahan operasional lokal, kami menyertakan utilitas backup otomatis. 

Jalankan backup data manual secara berkala dengan perintah:
```bash
# Berikan izin eksekusi jika berjalan di Linux/macOS
chmod +x scripts/backup_db.sh

# Eksekusi skrip backup
./scripts/backup_db.sh
```
> [!NOTE]
> Skrip ini secara otomatis mendeteksi apakah database berjalan di dalam container Docker `pgn_db` atau lokal host, lalu mengekstrak skema dan data transaksi ke dalam direktori aman `./docker/backups/` dengan penamaan terstruktur berbasis timestamp. File dump SQL ini terlindung di dalam `.gitignore`.

---

## 🌐 Akses Root UI & Konsumsi API (Siap di-Consume)

Server API PGNServer kini secara resmi dinyatakan **"Siap di-Consume"** untuk integrasi penuh dengan tim Frontend dan Client Apps (seperti QControl Desktop Client).

### 🖥️ 1. Root Landing Page Dashboard
Akses landing page interaktif bawaan server (yang dipaketkan secara mandiri via `go:embed`) langsung di browser Anda:
*   **Root Dashboard URL**: [http://localhost:8080/](http://localhost:8080/)
*   **Fitur**: Menampilkan metrik kesehatan sistem secara realtime, telemetri operasional, status database, serta panduan interaksi API.

### 📚 2. REST API Kontrak & Swagger Docs
Semua kontrak endpoint terdokumentasi secara interaktif via Swagger UI. Anda bisa melakukan testing secara dinamis dengan Bearer Token JWT.
*   **Swagger URL**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
*   **Bahan Uji Cepat**: Gunakan berkas [test.http](file:///c:/Software/PGNServer/test.http) yang kompatibel dengan VSCode REST Client untuk melakukan simulasi request registrasi, login, pengiriman checksheet, hingga pelacakan akar masalah BOM.

---

## 🔒 SOP Pengembangan & Kebijakan Rilis Beta

Demi menjaga stabilitas sistem selama fase rilis *Beta Production*, seluruh kontributor wajib mematuhi aturan berikut secara ketat:

| SOP Category | Kebijakan & Aturan Teknis | Dampak & Konsekuensi |
| :--- | :--- | :--- |
| **No-Branching** | Semua commit & push wajib langsung ke cabang `main`. Dilarang membuat cabang fitur baru secara terpisah. | Continuous Integration yang linier, bebas konflik penggabungan (*merge conflict*). |
| **Nomenklatur** | Variabel, skema model, dan endpoint wajib memakai Bahasa Indonesia sesuai Pedoman Umum Ejaan Bahasa Indonesia (PUEBI). | Kemudahan telusur (*traceability*) secara manajerial operasional internal. |
| **Security First** | Endpoint penulisan / mutasi data wajib melampirkan valid token JWT dengan klaim role khusus `Leader QC` (NIP default: `2211019`). | Mencegah penyalahgunaan data kualitas oleh pihak non-otoritas. |
| **Silent Git** | Jejak biner kompilasi (`*.exe`), session profiling (`.vscode/`), dan temporary file wajib diisolasi di `.gitignore`. | Menjaga kebersihan repositori awan dari berkas-berkas sampah. |

---
*PGN Quality Assurance & Database Architecture Dept. - Hak Cipta Dilindungi © 2026*
