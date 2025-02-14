# Gunakan base image Golang
FROM golang:1.23.2 AS builder

# Set working directory dalam container
WORKDIR /app

# Copy semua file ke dalam container
COPY . .

# Unduh dependencies
RUN go mod tidy

# Compile aplikasi
RUN go build -o main .

# Install migrate CLI
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Gunakan stage runtime agar lebih ringan
FROM golang:1.23.2

# Set working directory dalam container
WORKDIR /app

# Copy binary aplikasi dan migrate CLI dari builder
COPY --from=builder /app/main /app/main
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Expose port sesuai aplikasi
EXPOSE 8080

# Jalankan migrasi sebelum API dimulai
CMD migrate -database "mysql://root:rootpassword@tcp(db:3306)/mydb" -path db/migrations up && ./main
