package ws

import (
	"log"
	"net/http"

	"github.com/niluwats/gochat/pkg/redisrepo"
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Println("serveWs called")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading to tcp ", err)
	}

	client := &Client{Conn: ws}
	clients[client] = true

	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	log.Println("clients ", len(clients), clients, ws.RemoteAddr())

	receiver(client)
	println("existing", ws.RemoteAddr().String())
	delete(clients, client)
}

func StartWebSocketServer() {
	log.Println("websocket server")
	redisClient := redisrepo.InitializeRedis()
	defer redisClient.Close()

	http.HandleFunc("/ws", serveWs)

	go broadcaster()
	http.ListenAndServe(":8081", nil)
}
