package services

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

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

func (rmq *RabbitMQ) RMQRPC(request *LockingRequest, exchange string, corrId string) (res float64, err error) {
	body, err := json.Marshal(request)
	queue, err := rmq.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	msgs, err := rmq.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = rmq.channel.PublishWithContext(
		ctx,
		"",
		exchange,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       queue.Name,
			Body:          body,
		})

	log.Println("Queue status: ", queue)
	log.Printf("Send request thru rabbitmq: %s", body)

	for msg := range msgs {
		if corrId == msg.CorrelationId {
			res, err = strconv.ParseFloat(string(msg.Body), 64)
			break
		}
	}
	return
}

func ConnectRMQ() *RabbitMQ {
	rmq, err := NewRabbitMQ("amqp://guest:guest@127.0.0.1:5672/")
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %s", err)
	}
	return rmq
}

func GetTotalAmount(req *LockingRequest) (float64, error) {
	rmq := ConnectRMQ()
	log.Println("Trying to publish!")
	res, err := rmq.RMQRPC(req, "total_price_rmq", "unique_id_1")

	if err != nil {
		log.Printf("Failed to publish seat update request: %s", err)
		return -1.99, err
	}
	return res, nil
}
