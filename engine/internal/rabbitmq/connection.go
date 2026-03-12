package rabbitmq

import (
	"log"

	"github.com/Prayas-35/ragkit/engine/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() *amqp.Connection {
	conn, err := amqp.Dial(config.LoadConfig().RABBITMQ_URL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	return conn
}
