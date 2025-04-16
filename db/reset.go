package db

import (
	"aro-shop/models"
	"log"
)

func ResetDB() {
	InitDB()

	log.Println("Resetting database...")

	// DROP all tables
	err := DB.Migrator().DropTable(
		&models.TransactionItem{},
		&models.Transaction{},
		&models.Payment{},
		&models.PaymentMethod{},
		&models.Product{},
		&models.Category{},
		&models.Notification{},
		&models.User{},
	)

	if err != nil {
		log.Fatal("Failed to drop tables: ", err)
	}

	log.Println("Tables dropped successfully.")

	// MIGRATE ulang
	Migrate() // Pakai fungsi migrasi yang sudah kamu buat

	log.Println("Database has been reset and migrated successfully!")
}
