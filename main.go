package main

import (
	"context"
	"github.com/autom8ter/thermomatic/internal/server"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s, err := server.NewServer(&server.Config{
		TcpPort:         1337,
		HttpPort:        1338,
		ClientLogPrefix: "Thermomatic-Client: ",
		ServerLogPrefix: "Thermomatic-Server: ",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	s.Listen(ctx)
}
