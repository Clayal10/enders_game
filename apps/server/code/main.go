package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Clayal10/enders_game/apps/server/code/server"
)

const (
	defaultPort = 5069
)

var (
	cfg = &server.ServerConfig{
		Port: defaultPort,
	}
)

func main() {
	cancelFunctions, err := server.New(cfg)
	fatalOnErr(err)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	<-ch
	log.Println("Terminating Server")

	for _, cf := range cancelFunctions {
		cf()
	}
}

func fatalOnErr(err error) {
	if err != nil {
		log.Fatalf("%v: Could not start server!", err.Error())
	}
}
