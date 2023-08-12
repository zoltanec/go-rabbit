package rabbit

import (
	logger "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Consumer struct {
	log *logger.Logger
}

func NewConsumer() (c *Consumer) {
	c = &Consumer{}
	c.log = logger.New()
	return
}

func (c *Consumer) Consume(channel *amqp.Channel) {
	time.Sleep(time.Minute / 5)
	c.log.Printf("start consuming ...")
	//channel, err := app.channel
	//channel, err := connection.Channel()
	msgs, err := channel.Consume(
		"nobrakes", // queue
		"",         // consumer
		false,      // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        //args
	)

	if err != nil {
		panic(err)
	}
	// print consumed messages from queue
	forever := make(chan bool)

	numThreads := 6
	for i := 1; i <= numThreads; i++ {
		go func() {
			for {
				data := <-msgs
				c.log.Printf("<- Received Message: %s\n", data.Body)
				data.Ack(false)
				time.Sleep(time.Second * 10)
			}
		}()
	}

	c.log.Printf("<- Waiting for messages...")
	<-forever
}
