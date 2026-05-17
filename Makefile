.PHONY: build run swag test

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
