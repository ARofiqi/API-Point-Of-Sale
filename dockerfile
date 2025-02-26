# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git curl wget

# Install Reflex untuk hot reload
RUN go install github.com/cespare/reflex@latest

# Copy go.mod dan go.sum untuk dependency caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy seluruh kode
COPY . .

# Stage 2: Runtime
FROM golang:1.23-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy kode dari stage builder
COPY --from=builder /app /app
COPY --from=builder /go/bin/reflex /usr/local/bin/reflex

# Expose port
EXPOSE 8080

# Jalankan aplikasi dengan hot reload menggunakan Reflex
CMD ["reflex", "-r", "\\.go$", "--", "go", "run", "."]

