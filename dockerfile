# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git & dependencies yang diperlukan
RUN apk add --no-cache git

# Copy go.mod dan go.sum terlebih dahulu untuk caching yang lebih efisien
COPY go.mod go.sum ./
RUN go mod tidy

# Copy seluruh kode dan build aplikasi
COPY . .
RUN go build -o main

# Stage 2: Runtime (Final Image)
FROM alpine:latest

WORKDIR /app

# Copy hanya binary hasil build dari stage sebelumnya
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Jalankan aplikasi
# CMD ["./main"]

# Jalankan migrasi sebelum menjalankan aplikasi utama
CMD ["sh", "-c", "./migrate && ./main"]
