#!/bin/bash
# Script untuk melakukan backup data PostgreSQL secara manual

# Set environment variables
DB_NAME="pgn_db"
DB_USER="admin"
BACKUP_DIR="./docker/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
OUTPUT_FILE="${BACKUP_DIR}/backup_${DB_NAME}_${TIMESTAMP}.sql"

echo "========================================="
echo "🗄️  PGNServer - Database Backup Utility"
echo "========================================="

# Buat direktori backup jika belum ada
mkdir -p "${BACKUP_DIR}"

# Cek apakah container pgn_db sedang berjalan
if docker ps --format '{{.Names}}' | grep -q "^pgn_db$"; then
    echo "🐳 Menggunakan pg_dump dari container Docker pgn_db..."
    # Gunakan PGPASSWORD untuk menghindari prompt password
    docker exec -e PGPASSWORD=admin pgn_db pg_dump -U "${DB_USER}" -d "${DB_NAME}" > "${OUTPUT_FILE}"
    STATUS=$?
else
    echo "💻 Container Docker pgn_db tidak terdeteksi. Mencoba pg_dump lokal..."
    if command -v pg_dump &> /dev/null; then
        export PGPASSWORD="admin"
        pg_dump -h localhost -p 5432 -U "${DB_USER}" -d "${DB_NAME}" > "${OUTPUT_FILE}"
        STATUS=$?
    else
        echo "❌ Error: pg_dump tidak terdeteksi secara lokal maupun di container Docker!"
        exit 1
    fi
fi

if [ $STATUS -eq 0 ]; then
    echo "✅ Backup data berhasil disimpan secara aman!"
    echo "📍 Lokasi file: ${OUTPUT_FILE}"
    echo "📂 Ukuran file: $(du -sh "${OUTPUT_FILE}" | cut -f1)"
else
    echo "❌ Maaf, terjadi kendala saat memproses transaksi backup data."
    # Hapus file backup kosong jika gagal
    rm -f "${OUTPUT_FILE}"
    exit 1
fi
echo "========================================="
