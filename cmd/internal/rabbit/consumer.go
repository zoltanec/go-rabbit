package rabbit

import (
	logger "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"time"
	"fmt"
	"context"
	"net/http"
	"io"
	"errors"
)

type Consumer struct {
	log *logger.Logger
}

const (
    username = "root"
    password = "mroot"
    hostname = "127.0.0.1:3306"
    dbname   = "cons_requests"
)

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
				data.Ack(false)
				cons.log.Printf("<- Received Message: %s\n", data.Body)
				url, err := cons.getUrlFromMsg(data.Body)
				res, err := cons.makeRequest(url)
				if err != nil {
				    cons.log.Printf("Hui na ni: %s\n", res)
				    //I will publish new ampq message with url
				    channel.Publish(
                        "",         // exchange
                        "nobrakes", // key
                        false,      // mandatory
                        false,      // immediate
                        amqp.Publishing{
                            ContentType: "text/plain",
                            Body:        []byte(url),
                        },
                    )
				}
				cons.log.Printf("<- Received Data: %s\n", res)
				time.Sleep(time.Second * 10)
			}
		}()
	}

	cons.log.Printf("<- Waiting for messages...")
	<-forever
}

func dsn(dbName string) string {
    return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func (cons *Consumer) getUrlFromMsg(body []byte) (url string, err error) {
    //todo: write regexp checker for url
    url = string(body)
    if url == "" {
        err = errors.New("Wrong ampq message");
    }

    return url, err
}

func (cons *Consumer) makeRequest(url string) (p []byte, err error) {
    var client = &http.Client{
        Timeout: time.Second * 10,
    }

    cons.log.Info("ContentInfo request to " + url)

    res, err := client.Get(url)
    if err != nil {
        cons.log.Info("Unable to load content settings: " + err.Error())
        return
    }

    p, err = io.ReadAll(res.Body)
    if err != nil {
        cons.log.Info("Unable to read client settings reply: " + err.Error())
        return
    }

    /*
    err = json.Unmarshal(content, &reply)
    if err != nil {
    	g.Info("Incorrect json reply for contentId: " + err.Error())
    	return
    }
    */

    return p, nil
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