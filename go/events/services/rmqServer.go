package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/amirul-zafrin/event-ticketing/events.git/database"
	"github.com/amirul-zafrin/event-ticketing/events.git/models"
	"github.com/amirul-zafrin/event-ticketing/events.git/utilities"
	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s:%s", msg, err)
	}
}

func UpdateSeat(opt *LockingRequest) error {
	event := models.Events{}
	if err := utilities.FindEvent(opt.EventID, &event); err != nil {
		log.Printf("Failed to update seat: %s", err)
		return err
	}
	if event.Seats == nil {
		return errors.New("no seat category were set")
	}
	var resultMap map[string][]int
	jsonbytes, err := opt.Details.MarshalJSON()
	if err != nil {
		log.Printf("Error when marshalling")
		return err
	}
	err = json.Unmarshal(jsonbytes, &resultMap)
	if err != nil {
		log.Printf("Error when unmarshalling")
		return err
	}
	status := "locked"
	if opt.IsPaid {
		status = "occupied"
	}
	for _, element := range resultMap {
		for _, val := range element {
			if seatInfo, isMap := event.Seats[strconv.Itoa(val)].(map[string]interface{}); isMap {
				seatInfo["status"] = status
			}
		}
	}
	database.Database.Db.Save(&event)
	return nil
}

func GetTotalPrices(req *LockingRequest) float64 {
	err := UpdateSeat(req)
	if err != nil {
		log.Println(err)
	}
	event := models.Events{}
	if err := utilities.FindEvent(req.EventID, &event); err != nil {
		failOnError(err, "Cannot Find Event")
	}

	var resultMap map[string][]int
	jsonbytes, err := req.Details.MarshalJSON()
	if err != nil {
		failOnError(err, "Error when marshalling")
		return -1
	}
	err = json.Unmarshal(jsonbytes, &resultMap)
	if err != nil {
		failOnError(err, "Error when unmarshalling")
		return -1
	}

	total_cost := 0.00
	for key, element := range resultMap {
		prices, err := utilities.GetPricesByEventClass(event.ID, key)
		if err != nil {
			return -1
		}
		total_cost += (prices * float64(len(element)))
	}
	return total_cost
	// log.Println("Detail from get total prices: ", req)
	// return 0.001
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

func (rmq *RabbitMQ) RMQReply(queueName string, callback func(*LockingRequest) float64) {
	queue, err := rmq.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := rmq.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	forever := make(chan bool)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		for msg := range msgs {
			updReq := &LockingRequest{}
			err := json.Unmarshal(msg.Body, updReq)
			if err != nil {
				log.Printf("Error pasing update request %s", err)
			} else {
				response := callback(updReq)
				responseStr := strconv.FormatFloat(response, 'f', -1, 64)
				if err != nil {
					log.Printf("Error processing update request %s", err)
				}
				err = rmq.channel.PublishWithContext(
					ctx,
					"",
					msg.ReplyTo,
					false,
					false,
					amqp091.Publishing{
						ContentType:   "text/plain",
						CorrelationId: msg.CorrelationId,
						Body:          []byte(responseStr),
					})
				if err != nil {
					fmt.Println("Error from publisher RMQ Reply: ", err)
				}
			}
			fmt.Printf("Received Message: %s\n", msg.Body)
		}
	}()
	fmt.Println("Waiting for messages...")
	<-forever
}
