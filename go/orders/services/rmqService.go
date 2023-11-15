package services

import (
	"context"
	"encoding/json"
	"log"

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

func (rmq *RabbitMQ) PublishUpdateSeatStatus(request *LockingRequest, exchange string) error {
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	queue, err := rmq.channel.QueueDeclare(
		"seat_requests",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Error in declaring queue: %s", err)
		return err
	}
	ctx := context.Background()
	err = rmq.channel.PublishWithContext(
		ctx,
		"",
		exchange,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return err
	}
	log.Println("Queue status: ", queue)
	log.Printf("Send request thru rabbitmq: %s", body)

	return nil
}

func ConnectRMQ() *RabbitMQ {
	rmq, err := NewRabbitMQ("amqp://guest:guest@127.0.0.1:5672/")
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %s", err)
	}
	return rmq
}
