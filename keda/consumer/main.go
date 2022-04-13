package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

var (
	RABBITMQ_USER = os.Getenv("RABBITMQ_USER")
	RABBITMQ_PASS = os.Getenv("RABBITMQ_PASS")
	RABBITMQ_HOST = os.Getenv("RABBITMQ_HOST")
	RABBITMQ_PORT = os.Getenv("RABBITMQ_PORT")
)

func failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s %s", err, message)
	}
}

func main() {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", RABBITMQ_USER, RABBITMQ_PASS, RABBITMQ_HOST, RABBITMQ_PORT))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to connect to channel")
	defer ch.Close()

	forever := make(chan bool)

	channelName := "scale_out"
	q, err := ch.QueueDeclare(
		channelName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue for channel "+channelName)

	err = ch.Qos(
		1,
		0,
		false,
	)
	failOnError(err, "Failed to set QoS for channel "+channelName)

	messages, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer for channel "+channelName)

	for message := range messages {
		log.Printf("Channel %s received message %s", channelName, string(message.Body))
		if err != nil {
			log.Panicln(err)
		}
		message.Ack(false)
	}

	<-forever
}
