package queue

import (
	"aro-shop/db"
	"aro-shop/models"
	"encoding/json"
	"log"
)

// StartWorker menjalankan worker untuk memproses transaksi dari RabbitMQ
func StartWorker() {
	msgs, err := rabbitMQChannel.Consume(
		queueName, // Queue Name
		"",        // Consumer
		true,      // Auto-Ack
		false,     // Exclusive
		false,     // No-local
		false,     // No-wait
		nil,       // Args
	)
	if err != nil {
		log.Fatalf("❌ Gagal mengkonsumsi pesan dari queue: %v", err)
	}

	log.Println("👷 Worker transaksi berjalan...")

	for msg := range msgs {
		var t models.Transaction
		if err := json.Unmarshal(msg.Body, &t); err != nil {
			log.Println("❌ Gagal parsing transaksi:", err)
			continue
		}

		// Simpan transaksi ke database
		if err := db.DB.Create(&t).Error; err != nil {
			log.Println("❌ Gagal menyimpan transaksi:", err)
			continue
		}

		log.Println("✅ Transaksi berhasil disimpan ke database:", t.ID)
	}
}
