# PGNServer - Sistem Manufaktur & Kualitas Terpadu

PGNServer adalah sistem *backend* monolitik modular berbasis Go 1.25.x yang dirancang khusus untuk ekosistem tata niaga manufaktur dan inspeksi kualitas (QC). Sistem ini memanfaatkan pangkalan data PostgreSQL 17.x, kerangka kerja Gin (untuk *routing* RESTful), dan GORM (sebagai ORM), dengan fokus pada performa, integritas data (ACID), dan skalabilitas.

## 🏗 Arsitektur Sistem

Proyek ini menerapkan konsep **Modular Monolith** dengan pendekatan berorientasi pada domain (*Domain-Driven Design*), di mana setiap entitas bisnis (manufaktur, kualitas, otentikasi) dipisahkan ke dalam modul mandiri di bawah direktori `internal/`. Nomenklatur dan variabel menggunakan terminologi **Bahasa Indonesia** sesuai dengan standar konvensi operasional perusahaan.

### Struktur Utama:

*   **`cmd/api/`**: Titik masuk aplikasi (`main.go`). Berisi inisialisasi koneksi pangkalan data, pengaturan *middleware*, konfigurasi *router*, dan pendaftaran *handler*.
*   **`internal/otentikasi/`**: Modul manajemen akses pengguna, pendaftaran (registrasi), proses masuk (login), dan lupa sandi menggunakan **JWT v5** serta enkripsi *bcrypt*.
*   **`internal/manufaktur/`**: Modul tata niaga dan *Master Data*. Mengelola entitas **Pemasok**, **Material**, dan **Komposisi Material BOM** (Bill of Materials) dengan aturan Kunci Asing (*Foreign Keys*) yang kohesif.
*   **`internal/kualitas/`**: Modul fungsionalitas inspeksi *Quality Control* (QC). Menangani operasi perekaman data lembar periksa (*checksheet*), penguraian hierarki BOM, dan pelacakan cacat material (*defect ledger*) menggunakan fungsi SQL tingkat lanjut seperti `JSON_TABLE`.
*   **`internal/infrastruktur/`**: Modul konfigurasi teknis lintas domain (seperti pengaturan koneksi, penanganan automigrasi skema, penjaga sesi JWT).
*   **`pkg/respon/`**: Pustaka standardisasi lapisan penyajian (*presentation layer*) yang membungkus respon *HTTP API* JSON ke antarmuka aplikasi.
*   **`pkg/pencatatan_log/`**: Pustaka standardisasi pembuatan log sistem.

## 🚀 Fitur Utama

1.  **Eksekusi Transaksi Atomik Terpadu**: Di modul Kualitas, penyimpanan ribuan baris data *Detail Inspeksi* dilakukan tanpa jebakan performa "N+1" melalui eksploitasi skema *transmisi tunggal* (*batch query* berbasis `JSON_TABLE` di PostgreSQL 17.x).
2.  **Pelacakan Pohon Cacat Otomatis (Auto-Trace NG)**: Sistem mendeteksi ketika suatu lini produksi mengalami cacat (NG - *No Good*), membongkar pohon BOM untuk menemukan material akar permasalahan, lalu mencatatnya secara otonom ke dalam Buku Besar Cacat.
3.  **Keamanan JWT & Hash Kata Sandi**: Jalur masuk dilindungi menggunakan mekanisme *Bearer Token*, dengan penyandian basis sandi yang aman.
4.  **Swagger / OpenAPI Ready**: Mendukung *generasi otomatis* dokumentasi REST API melalui integrasi `swag`.

## ⚙️ Prasyarat & Instalasi

*   **Go** v1.25.x atau lebih baru.
*   **PostgreSQL** v17.x (dapat dijalankan via Docker Compose).
*   Berkas lingkungan (`.env`) berisi:
    ```env
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=admin
    DB_PASSWORD=admin
    DB_NAME=pgn_db
    APP_PORT=8080
    JWT_SECRET=super-rahasia-jangan-disebar
    ```

### Menjalankan Sistem Secara Lokal

1.  Jalankan modul pangkalan data:
    ```bash
    docker-compose up -d
    ```
2.  Kompilasi dan jalankan server Go:
    ```bash
    go run cmd/api/main.go
    ```
3.  Server PGNServer akan berjalan pada *port* yang telah didefinisikan (secara default `8080`).

## 📚 Panduan Agen AI (Agentic Workflows)

Repositori ini menyertakan panduan kontekstual khusus untuk asisten AI (seperti Google Gemini atau Claude) guna memahami peran mereka di dalam ekosistem:

*   **`AGENTS.MD`**: Merumuskan profil (Persona) agen yang berperan sebagai "Leader QC & QA Analyst" dan "Database Architect".
*   **`GEMINI.MD`**: Panduan integrasi Gemini AI menggunakan Model Context Protocol (MCP) untuk terhubung secara langsung ke pangkalan data `pgn_db` demi tujuan eksekusi kueri analitik (7 QC Tools).

---
*PGNServer Backend Development - Hak Cipta Dilindungi*
