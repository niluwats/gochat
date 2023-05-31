package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/niluwats/gochat/pkg/httpserver"
	"github.com/niluwats/gochat/pkg/ws"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("unable to load environment variables ", err)
	}
}

func main() {
	go httpserver.StartHTTPServer()
	go ws.StartWebSocketServer()
	select {}
}
