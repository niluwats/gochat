package main

import (
	"github.com/niluwats/gochat/pkg/httpserver"
	"github.com/niluwats/gochat/pkg/ws"
)

func main() {
	go ws.StartWebSocketServer()
	go httpserver.StartHTTPServer()
	select {}
}
