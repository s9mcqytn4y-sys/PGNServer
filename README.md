# PGNServer - Manajemen Ekosistem Terpadu

Sebuah aplikasi backend monolitik berbasis *Go 1.25.x* dan *PostgreSQL 17.x* yang dibangun khusus untuk mengelola operasional dan kualitas tata niaga gas serta komponen material manufaktur. Sistem ini dimotori oleh framework `gin-gonic/gin` sebagai _router_ berkinerja tinggi, berpadu dengan ketangguhan ORM `gorm.io/gorm` guna memastikan integritas data.

## 📦 Ekosistem & Arsitektur

PGNServer mematuhi desain **Arsitektur Monolitik Modular** menggunakan ejaan Bahasa Indonesia yang terstandardisasi guna meningkatkan kohesi domain internal dan menghindari ambiguitas istilah generik.
Subdomain yang ada mencakup:
1. `internal/otentikasi/`: Logika autentikasi berbasis JWT v5, registrasi pengguna baru dan pemulihan (*bcrypt* hash).
2. `internal/manufaktur/`: Entitas `Pemasok`, `Material`, dan skema berelasi struktural `KomposisiMaterialBOM`.
3. `internal/infrastruktur/`: Konfigurasi basis data, migrasi GORM (*AutoMigrate*), dan sistem pencegatan keamanan (`middleware.go`).
4. `pkg/respon/`: Pelaporan respon API ramah-antarmuka tanpa mengekspos rincian basis data mentah.

## 🚀 Instalasi dan Orkestrasi Kontainer

Proyek ini telah diamankan dan dikonfigurasi melalui platform kontainer Docker (*Multi-Stage Build*). Variabel kompilasi eksperimental `GOEXPERIMENT=jsonv2` dari lingkungan Go 1.25 diinjeksikan secara transparan untuk menyokong fungsionalitas JSON ultra-ringan masa depan.

```bash
# Menjalankan ekosistem layanan Postgres 17.x dan PGN_API (pelabuhan 8080)
docker-compose up --build -d
```

## 🔐 Keamanan Sesi

Semua jalur terproteksi (di masa depan) diwajibkan menyertakan tajuk (*header*) `Authorization: Bearer <TOKEN>`. Modul infrastruktur mencegat permintaan untuk mendeteksi `JWT_SECRET` yang dimuat dari `.env`.

## 📚 Open API (Swagger)

Jalur interaksi dan dokumentasi akan di-hosting secara statis, melacak aneka rute API seperti:
- `GET /api/v1/cek_sistem`
- `POST /api/v1/otentikasi/daftar`
- `POST /api/v1/otentikasi/masuk`
- `POST /api/v1/otentikasi/lupa-sandi`

---
*Dikembangkan oleh PGN Quality Assurance Dept. - 2026*
