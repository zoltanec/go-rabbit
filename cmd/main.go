package main

import (
	"arbitrag/queue-app/cmd/internal"
	"arbitrag/queue-app/cmd/internal/rabbit"
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	log.Printf("Initializing application")
	fmt.Println("RabbitMQ in Golang: Getting started tutorial")

	ctx := context.Background()
	//host.docker.internal
	connection, err := amqp.Dial("amqp://rbuser:rbpassword@localhost:5672")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	appInstance, err := internal.NewApplication(connection)
	if err != nil {
		log.Fatalf("failed to initialize app: %s", err.Error())
	}
	producer := rabbit.NewProducer(appInstance.Channel)
	consumer := rabbit.NewConsumer()
	go producer.Produce(appInstance.Channel, ctx)

	//Думаю тут можно без горутины, правильно?
	//Но для удобности можно горутину, если дополнительные модули запускать?
    consumer.Consume(appInstance.Channel)
	appInstance.Run()
}
