FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main

EXPOSE 8080

CMD ["./main"]

# Jalankan migrasi sebelum menjalankan aplikasi utama
# CMD ["sh", "-c", "./migrate && ./main"]
