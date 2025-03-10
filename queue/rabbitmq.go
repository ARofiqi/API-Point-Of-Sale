package queue

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

var rabbitMQConn *amqp.Connection
var rabbitMQChannel *amqp.Channel
var queueName = "transaction_queue"

func InitRabbitMQ() {
	var err error

	// Koneksi ke RabbitMQ
	rabbitMQConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke RabbitMQ: %v", err)
	}

	// Buat channel
	rabbitMQChannel, err = rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("❌ Gagal membuat channel: %v", err)
	}

	// Deklarasi antrian
	_, err = rabbitMQChannel.QueueDeclare(
		queueName,
		true,  // Durable
		false, // Auto-delete
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("❌ Gagal mendeklarasikan queue: %v", err)
	}

	log.Println("✅ RabbitMQ siap digunakan!")
}

// PublishTransaction mengirim transaksi ke queue
func PublishTransaction(message []byte) error {
	err := rabbitMQChannel.Publish(
		"",        // Exchange
		queueName, // Routing Key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}

	log.Println("📨 Transaksi berhasil dikirim ke queue")
	return nil
}

// CloseRabbitMQ menutup koneksi RabbitMQ
func CloseRabbitMQ() {
	if rabbitMQConn != nil {
		rabbitMQConn.Close()
	}
	if rabbitMQChannel != nil {
		rabbitMQChannel.Close()
	}
}
