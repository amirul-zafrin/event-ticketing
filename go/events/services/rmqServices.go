package services

import (
	"encoding/json"
	"fmt"
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
func (rmq *RabbitMQ) ConsumeUpdate(queueName string, callback func(*LockingRequest) error) {
	msgs, err := rmq.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Printf("Failed to create consume channel: %s", err)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			updReq := &LockingRequest{}
			err := json.Unmarshal(msg.Body, updReq)
			if err != nil {
				log.Printf("Error pasing update request %s", err)
			} else {
				err = callback(updReq)
				if err != nil {
					log.Printf("Error processing update request %s", err)
				}
			}
			fmt.Printf("Received Message: %s\n", msg.Body)
		}
	}()
	fmt.Println("Waiting for messages...")
	<-forever
}
