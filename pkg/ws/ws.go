package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/niluwats/gochat/model"
	"github.com/niluwats/gochat/pkg/redisrepo"
)

type Client struct {
	Conn     *websocket.Conn
	Username string
}

type Message struct {
	Type string     `json:"type"`
	User string     `json:"user,omitempty"`
	Chat model.Chat `json:"chat,omitempty"`
}

var clients = make(map[*Client]bool)
var broadcast = make(chan *model.Chat)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		allowedOrigin := "http://localhost:3000"
		return r.Header.Get("Origin") == allowedOrigin
	},
}

func receiver(client *Client) {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		m := &Message{}

		err = json.Unmarshal(p, m)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("HOST ", client.Conn.RemoteAddr())

		if m.Type == "bootup" {
			client.Username = m.User
			log.Println("client successfully mapped ", client.Username)
		} else {
			log.Println("recevied chat message", m.Type, m.Chat)

			chat := m.Chat
			chat.Timestamp = time.Now().Unix()

			id, err := redisrepo.CreateChat(&chat)
			if err != nil {
				log.Println("error saving chat in redis ", err)
				return
			}

			chat.ID = id
			broadcast <- &chat
		}
	}

}

func broadcaster() {
	for {
		message := <-broadcast
		fmt.Println("new message", message)

		for client := range clients {
			if client.Username == message.From || client.Username == message.To {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("web socket error %s ", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}
