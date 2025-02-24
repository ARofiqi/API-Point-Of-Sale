#!/bin/bash

# Load .env secara manual
if [ -f .env ]; then
    export $(grep -E '^[A-Z_]+=' .env | xargs)
fi

# Path ke database migrations
MIGRATE="migrate -database \"mysql://$DB_USER:$DB_PASS@tcp($DB_HOST:$DB_PORT)/$DB_NAME\" -path db/migrations"

# Fungsi untuk menjalankan perintah migrate
migrate_up() {
    eval "$MIGRATE up"
}

migrate_down() {
    eval "$MIGRATE down 1"
}

migrate_reset() {
    eval "$MIGRATE down"
    eval "$MIGRATE up"
}

run() {
    go run main.go
}

# Menangani argumen perintah
case "$1" in
    up)
        migrate_up
        ;;
    down)
        migrate_down
        ;;
    reset)
        migrate_reset
        ;;
    run)
        run
        ;;
    *)
        echo "Usage: $0 {up|down|reset|run}"
        exit 1
        ;;
esac
