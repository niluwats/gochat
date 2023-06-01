package ws

import (
	"log"
	"net/http"

	"github.com/niluwats/gochat/pkg/redisrepo"
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{Conn: ws}
	clients[client] = true

	log.Println("clients ", len(clients), clients, ws.RemoteAddr())

	receiver(client)
	println("existing", ws.RemoteAddr().String())
	delete(clients, client)
}

func setupRoutes() {
	http.HandleFunc("/ws", serveWs)
}

func StartWebSocketServer() {
	log.Println("websocket server")
	redisClient := redisrepo.InitializeRedis()
	defer redisClient.Close()

	go broadcaster()
	setupRoutes()
	http.ListenAndServe(":8081", nil)
}
