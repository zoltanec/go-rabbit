package rabbit

import (
	"context"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"math/rand"
	"time"
)

type Producer struct {
	log *logger.Logger
	queue amqp.Queue
}

var resources = []string{
    "http://hui.vam",
    "https://api.myip.com",
    "http://edns.ip-api.com/json",
    "https://v3.football.api-sports.io/",
}

func NewProducer(channel *amqp.Channel) (p *Producer) {
	p = &Producer{}
	p.log = logger.New()

	return
}

func (p *Producer) Produce(channel *amqp.Channel, ctx context.Context) {
	fmt.Println("start producer ...")

    //Кажется это не пересоздает очередь
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

	fmt.Println("-> Running cycle")
	for i := 1; i < 6; i++ {
		// publishing a message
		err := channel.Publish(
			"",         // exchange
			"nobrakes", // key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(resources[rand.Intn(2)]),
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
