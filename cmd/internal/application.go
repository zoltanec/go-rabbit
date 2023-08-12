package internal

import (
	"arbitrag/queue-app/cmd/internal/web"
	logger "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Application struct {
	log        *logger.Logger
	config     Configuration
	api        web.Api
	connection *amqp.Connection
	Channel    *amqp.Channel
}

type Configuration struct {
	Port int
}

func NewApplication(connection *amqp.Connection) (a *Application, err error) {
	app := &Application{}
	formatter := &logger.TextFormatter{FullTimestamp: true}
	logger.SetFormatter(formatter)

	app.log = logger.New()

	api, err := web.NewApi(app.log)
	if err != nil {
		return
	}

	api.SetPort(18000)
	app.api = api

	app.connection = connection
	app.Channel, err = app.connection.Channel()
	//defer app.Channel.Close()

	//defer app.channel.Close()
	if err != nil {
		panic(err)
	}

	return app, nil
}

func (a *Application) Run() {
	a.log.Info("Application started")
	a.api.Run()
}
