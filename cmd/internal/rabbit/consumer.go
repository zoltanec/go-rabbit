package rabbit

import (
	logger "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"time"
	"fmt"
	"context"
)

type Consumer struct {
	log *logger.Logger
}

func NewConsumer() (cons *Consumer) {
	cons = &Consumer{}
	cons.log = logger.New()
	return
}

func (cons *Consumer) Consume(channel *amqp.Channel) {
	time.Sleep(time.Minute / 5)
	cons.log.Printf("start consuminggg ...")
	cons.saveDb()
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
				cons.log.Printf("<- Received Message: %s\n", data.Body)
				data.Ack(false)
				time.Sleep(time.Second * 10)
			}
		}()
	}

	cons.log.Printf("<- Waiting for messages...")
	<-forever
}

const (
    username = "root"
    password = "mroot"
    hostname = "127.0.0.1:3306"
    dbname   = "cons_requests"
)

func dsn(dbName string) string {
    return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}


//*sql.DB
func (cons *Consumer) saveDb() {

    db, err := sql.Open("mysql", dsn(""))
    defer db.Close()

    ctx, cancelfunc := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancelfunc()
    res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS " + dbname)

    if err != nil {
        cons.log.Printf("Error %s when creating DB\n", err)
        //return nil, err
    }

    no, err := res.RowsAffected()
    if err != nil {
        cons.log.Printf("Error %s when fetching rows", err)
        //return nil, err
    }
    cons.log.Printf("rows affected: %d\n", no)



    if err != nil {
        cons.log.Fatal(err)
    }

    var version string
    db.QueryRow("SELECT VERSION()").Scan(&version)
    fmt.Println("Connected to:", version)
    db.Close()
}

//migrate pack https://dev.to/techschoolguru/how-to-write-run-database-migration-in-golang-5h6g
