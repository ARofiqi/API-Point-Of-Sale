package queue

import (
	"aro-shop/config"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

var (
	rabbitMQConn      *amqp.Connection
	rabbitMQChannel   *amqp.Channel
	transactionQueue  = "transaction_queue"
	notificationQueue = "notification_queue"
	cfg               = config.LoadConfig()
	mu                sync.Mutex
	closeOnce         sync.Once
	notifyCloseChan   chan *amqp.Error
)

func InitRabbitMQ() {
	var err error

	rabbitMQConn, err = amqp.Dial(cfg.RABBITMQUrl)
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke RabbitMQ: %v", err)
	}

	rabbitMQChannel, err = rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("❌ Gagal membuat channel: %v", err)
	}

	notifyCloseChan = rabbitMQChannel.NotifyClose(make(chan *amqp.Error))

	_, err = rabbitMQChannel.QueueDeclare(
		transactionQueue,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("❌ Gagal mendeklarasikan transactionQueue: %v", err)
	}

	_, err = rabbitMQChannel.QueueDeclare(
		notificationQueue,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("❌ Gagal mendeklarasikan notificationQueue: %v", err)
	}

	log.Println("🚀 RabbitMQ siap digunakan!")
}

func ensureChannel() error {
	mu.Lock()
	defer mu.Unlock()

	select {
	case <-notifyCloseChan:
		log.Println("⚠️ Channel terdeteksi tertutup, membuat ulang channel...")
		var err error
		rabbitMQChannel, err = rabbitMQConn.Channel()
		if err != nil {
			log.Printf("❌ Gagal membuat ulang channel: %v", err)
			return err
		}
		notifyCloseChan = rabbitMQChannel.NotifyClose(make(chan *amqp.Error))
	default:
	}

	return nil
}

func PublishTransaction(message []byte) error {
	if err := ensureChannel(); err != nil {
		return err
	}

	err := rabbitMQChannel.Publish(
		"",               // Exchange
		transactionQueue, // Routing Key
		false,            // Mandatory
		false,            // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Printf("❌ Gagal mengirim transaksi ke queue: %v", err)
		return err
	}

	log.Println("📨 Transaksi berhasil dikirim ke queue")
	return nil
}

func PublishNotification(message string) error {
	if err := ensureChannel(); err != nil {
		return err
	}

	err := rabbitMQChannel.Publish(
		"",                // Exchange
		notificationQueue, // Routing Key
		false,             // Mandatory
		false,             // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Printf("❌ Gagal mengirim notifikasi ke queue: %v", err)
		return err
	}

	log.Println("📨 Notifikasi berhasil dikirim ke queue")
	return nil
}

func CloseRabbitMQ() {
	closeOnce.Do(func() {
		if rabbitMQChannel != nil {
			rabbitMQChannel.Close()
		}
		if rabbitMQConn != nil {
			rabbitMQConn.Close()
		}
		log.Println("✅ RabbitMQ connection closed")
	})
}
