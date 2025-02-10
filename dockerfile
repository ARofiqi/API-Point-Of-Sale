# Gunakan base image Golang
FROM golang:1.23.2

# Set working directory dalam container
WORKDIR /app

# Copy semua file ke dalam container
COPY . .

# Unduh dependencies
RUN go mod tidy

# Compile aplikasi
RUN go build -o main .

# Expose port sesuai aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]