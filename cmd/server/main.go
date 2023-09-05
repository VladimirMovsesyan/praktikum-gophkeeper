package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"praktikum-gophkeeper/pkg/configuration"
	"syscall"
)

var (
	flAddress = flag.String("a", ":8080", "Server's address.") // RUN_ADDRESS
	flDSN     = flag.String("d", "123", "Server's URI.")       // DSN
)

func main() {
	flag.Parse()

	config, err := configuration.NewServer(flAddress, flDSN)
	if err != nil {
		log.Println(err)
		return
	}

	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		log.Println("Server starting...")
		if err := config.Server.Serve(listener); err != nil {
			log.Println(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	sig := <-signals

	log.Println("Got signal:", sig)
	config.Server.GracefulStop()
}
