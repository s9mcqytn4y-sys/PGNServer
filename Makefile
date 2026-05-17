.PHONY: build run swag test docker-up docker-down release-beta

# Sinkronisasi pembaruan Swagger
swag:
	swag init -g cmd/api/main.go

# Kompilasi kode biner
build:
	go build -o tmp/server.exe cmd/api/main.go

# Menjalankan server lokal
run:
	go run cmd/api/main.go

# Menjalankan tes
test:
	go test ./...

# Docker Compose Up
docker-up:
	docker compose up --build -d

# Docker Compose Down
docker-down:
	docker compose down

# Rilis Beta ke GitHub
release-beta:
	gh release create v1.0.0-beta --title "v1.0.0-beta" --notes "BETA Release PGNServer Backend Modular Monolith with secure transactions, GORM rollback protection, live dashboard telemetry, and worker pool concurrency control."
