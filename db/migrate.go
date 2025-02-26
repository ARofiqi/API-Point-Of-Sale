package db

import (
	"aro-shop/models"
	"log"
)

func main() {
	InitDB()

	log.Println("Starting database migration...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.Category{},
	)
	if err != nil {
		log.Fatal("Migration failed : ", err)
	}

	log.Println("Migration completed successfully!")
}
