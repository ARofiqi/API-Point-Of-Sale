package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection

func InitRabbitMQ(url string) error {
	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		log.Printf("Gagal konek ke RabbitMQ: %s", err)
		return err
	}
	return nil
}

func PublishMessage(queueName string, body string) error {
	if conn == nil {
		return amqp.ErrClosed
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	return err
}
