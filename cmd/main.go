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
	//producer := rabbit.NewProducer()
	consumer := rabbit.NewConsumer()
	go rabbit.Produce(appInstance.Channel, ctx)
	//Запустить в пять потоков
	go consumer.Consume(appInstance.Channel)
	appInstance.Run()
}
