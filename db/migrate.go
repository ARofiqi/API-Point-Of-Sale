package db

import (
	"aro-shop/models"
	"log"
)

func Migrate() {
	InitDB()

	log.Println("Starting database migration...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.TransactionItem{},
		&models.Transaction{},
		&models.Category{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatal("Migration failed : ", err)
	}

	log.Println("Migration completed successfully!")
}
