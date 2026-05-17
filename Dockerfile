# --- Tahap 1: Pembangun (Builder) ---
FROM golang:1.25-alpine AS builder

# Inject env untuk mengaktifkan Go 1.25 JSON v2 eksperimen
ENV GOEXPERIMENT=jsonv2
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Instal dependencies root OS minimal
RUN apk add --no-cache tzdata

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Kompilasi static binary yang aman dan efisien
RUN go build -ldflags="-s -w" -o pgn_api ./cmd/api

# --- Tahap 2: Final Image (Distroless / Minimalist) ---
FROM alpine:latest

# Konfigurasi zona waktu secara lengkap
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /usr/share/zoneinfo/Asia/Jakarta /etc/localtime
RUN echo "Asia/Jakarta" > /etc/timezone
ENV TZ=Asia/Jakarta

# Security Hardening: Menjalankan aplikasi sebagai non-root user
RUN addgroup -S pgnteam && adduser -S pgnuser -G pgnteam

# Buat direktori kerja dan atur kepemilikan
WORKDIR /app
RUN mkdir -p /app/penyimpanan && chown -R pgnuser:pgnteam /app

COPY --from=builder /app/pgn_api .

# Berjalan sebagai non-root
USER pgnuser

# Port eksposur
EXPOSE 8080

# Jalankan server
CMD ["./pgn_api"]
