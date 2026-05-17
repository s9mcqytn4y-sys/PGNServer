# Tahap Kompilasi (Build Stage)
FROM golang:1.25-alpine AS pembangun

# Setel variabel lingkungan
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOEXPERIMENT=jsonv2

WORKDIR /app

# Salin modul dan unduh dependensi
COPY go.mod go.sum ./
RUN go mod download

# Salin kode sumber
COPY . .

# Kompilasi aplikasi dengan optimasi ukuran
RUN go build -ldflags="-s -w" -o pgn_api ./cmd/api

# Tahap Produksi (Final Stage)
FROM alpine:latest

WORKDIR /app

# Tambahkan sertifikat SSL dan zona waktu
RUN apk --no-cache add ca-certificates tzdata

# Salin biner hasil kompilasi
COPY --from=pembangun /app/pgn_api .
COPY --from=pembangun /app/.env .
# COPY --from=pembangun /app/docs ./docs # abaikan jika belum ada

# Buat direktori penyimpanan
RUN mkdir -p penyimpanan

EXPOSE 8080

CMD ["./pgn_api"]
