package services

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/datatypes"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type LockingRequest struct {
	EventID int
	IsPaid  bool
	Details datatypes.JSON
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}
