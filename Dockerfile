# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o pgnserver ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/pgnserver .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs

# Create storage directory
RUN mkdir -p penyimpanan

# Expose port
EXPOSE 8080

# Command to run
CMD ["./pgnserver"]
