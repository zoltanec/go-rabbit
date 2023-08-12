package web

import (
	"context"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Api struct {
	log     *logger.Logger
	context context.Context
	port    int
}

func NewApi(log *logger.Logger) (a Api, err error) {
	a.log = log
	a.log.Info("API initialized")

	a.context = context.Background()
	a.port = 8080

	return a, nil
}

func (a *Api) SetPort(port int) {
	a.port = port
	a.log.Infof("Api port is configured to: %d", a.port)
}

func (a *Api) Run() {
	c := xhandler.Chain{}
	c.UseC(xhandler.TimeoutHandler(10 * time.Second))

	mux := xmux.New()
	mux.GET("/status", xhandler.HandlerFuncC(a.Status))
	http.Handle("/", c.HandlerCtx(a.context, mux))

	srv := &http.Server{Addr: ":" + strconv.Itoa(a.port), Handler: xhandler.New(a.context, mux)}
	log.Print("Starting HTTP API")

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Got listen error: %s", err.Error())
	}
}
