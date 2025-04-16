package queue

import (
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func StartTransactionWorker() {
	for {
		if err := ensureChannel(); err != nil {
			log.Println("⚠️ Worker transaksi menunggu koneksi RabbitMQ pulih...")
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := rabbitMQChannel.Consume(
			transactionQueue, // Queue Name
			"",               // Consumer
			true,             // Auto-Ack
			false,            // Exclusive
			false,            // No-local
			false,            // No-wait
			nil,              // Args
		)
		if err != nil {
			log.Printf("❌ Gagal mengkonsumsi pesan dari queue transaksi: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("👷 Worker transaksi berjalan...")
		
		for msg := range msgs {
			var transaction models.Transaction
			if err := json.Unmarshal(msg.Body, &transaction); err != nil {
				log.Println("❌ Gagal parsing transaksi:", err)
				continue
			}

			if err := db.DB.Create(&transaction).Error; err != nil {
				log.Println("❌ Gagal menyimpan transaksi:", err)
				continue
			}

			log.Println("✅ Transaksi berhasil disimpan ke database:", transaction.ID)
		}
	}
}

func StartNotificationWorker() {
	for {
		if err := ensureChannel(); err != nil {
			log.Println("⚠️ Worker notifikasi menunggu koneksi RabbitMQ pulih...")
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := rabbitMQChannel.Consume(
			notificationQueue,
			"",
			true,  // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			log.Printf("❌ Gagal mengkonsumsi pesan dari queue notifikasi: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("👷 Worker notifikasi berjalan...")

		for msg := range msgs {
			message := string(msg.Body)
			fmt.Printf("📩 New Notification: %s\n", message)

			err := CreateNotification(message)
			if err != nil {
				log.Printf("❌ Gagal menyimpan notifikasi: %v", err)
			}
		}
	}
}

func CreateNotification(message string) error {
	notification := models.Notification{
		Message:   message,
		IsRead:    false,
		CreatedAt: utils.GetCurrentTime(),
	}

	if err := db.DB.Create(&notification).Error; err != nil {
		return err
	}
	return nil
}

func StartWorker() {
	go StartTransactionWorker()
	go StartNotificationWorker()

	log.Println("🚀 Semua worker berjalan...")
}
