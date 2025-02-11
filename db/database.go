package db

import (
	"aro-shop/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var cfg = config.LoadConfig()

	dbUser := cfg.DBUser
	dbPass := cfg.DBPass
	dbHost := cfg.DBHost
	dbPort := cfg.DBPort
	dbName := cfg.DBName

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database is not reachable:", err)
	}

	log.Println("Connected to database successfully")
}
