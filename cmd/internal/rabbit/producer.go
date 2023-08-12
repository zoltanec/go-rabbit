package rabbit

import (
	"context"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type Producer struct {
	log *logger.Logger
}

func NewProducer() (p *Producer) {
	p = &Producer{}
	p.log = logger.New()
	return
}

func Produce(channel *amqp.Channel, ctx context.Context) {
	fmt.Println("start producer ...")
	//channel := app.channel
	/*if err != nil {
		panic(err)
	}*/
	//defer channel.Close()

	// declaring queue with its properties over the the channel opened
	queue, err := channel.QueueDeclare(
		"nobrakes", // name
		true,       // durable
		false,      // auto delete
		false,      // exclusive
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		panic(err)
	}

	//go func() {
	//ctx.Err() == nil

	fmt.Println("-> Running cycle")
	for i := 1; i < 6; i++ {
		//fmt.Println("err", err)
		if err != nil {
			fmt.Println("-> Get error and sleep some")
			time.Sleep(time.Minute)
			continue
		}

		// publishing a message
		err = channel.Publish(
			"",         // exchange
			"nobrakes", // key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("http://ya.ru"),
			},
		)

		if err != nil {
			fmt.Println("-> Panic")
			panic(err)
		}

		queue, err = channel.QueueInspect("nobrakes")
		if err != nil {
			fmt.Println("-> Panic")
			panic(err)
		}
		fmt.Println("-> Queue status:", queue)
		fmt.Println("-> Successfully published message, sleep 3sec")

		time.Sleep(3 * time.Second)
	}
}
