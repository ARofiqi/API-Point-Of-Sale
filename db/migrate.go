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
		&models.PaymentMethod{},
		&models.Payment{},
		&models.TransactionItem{},
		&models.Transaction{},
		&models.Category{},
		&models.Notification{},
	)

	// Menambahkan index dengan B-Tree di PostgreSQL
	DB.Exec("CREATE INDEX idx_product_category_id ON products USING btree (category_id)")
	DB.Exec("CREATE INDEX idx_product_name ON products USING btree (name)")
	DB.Exec("CREATE INDEX idx_product_price ON products USING btree (price)")

	// Index pada Email (sudah unik, tapi tetap bisa eksplisit)
	DB.Exec("CREATE UNIQUE INDEX idx_users_email ON users USING btree (email)")
	// Index pada Role untuk query filter lebih cepat
	DB.Exec("CREATE INDEX idx_users_role ON users USING btree (role)")
	// Index pada Name untuk pencarian
	DB.Exec("CREATE INDEX idx_users_name ON users USING btree (name)")

	// Index pada TransactionID untuk mempercepat join transaksi
	DB.Exec("CREATE INDEX idx_transaction_items_transaction_id ON transaction_items USING btree (transaction_id)")
	// Index pada ProductID untuk pencarian transaksi berdasarkan produk
	DB.Exec("CREATE INDEX idx_transaction_items_product_id ON transaction_items USING btree (product_id)")
	// Index pada Date untuk pencarian transaksi berdasarkan tanggal
	DB.Exec("CREATE INDEX idx_transactions_date ON transactions USING btree (date)")

	if err != nil {
		log.Fatal("Migration failed : ", err)
	}

	log.Println("Migration completed successfully!")
}
