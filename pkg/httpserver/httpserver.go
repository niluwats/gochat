package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/niluwats/gochat/pkg/redisrepo"
	"github.com/rs/cors"
)

func StartHTTPServer() {
	log.Println("http server")
	redisClient := redisrepo.InitializeRedis()
	defer redisClient.Close()

	redisrepo.CreateFetchChatBetweenIndex()

	router := mux.NewRouter()
	router.HandleFunc("/status", func(w http.ResponseWriter, router *http.Request) {
		fmt.Fprintf(w, "Simple server")
	}).Methods(http.MethodGet)

	router.HandleFunc("/register", registerHandler).Methods(http.MethodPost)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	router.HandleFunc("/verify-contact", verifyContactHandler).Methods(http.MethodPost)
	router.HandleFunc("/chat-history", chatHistoryHandler).Methods(http.MethodGet)
	router.HandleFunc("/contact-list", contactListHandler).Methods(http.MethodGet)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	http.ListenAndServe(":8080", handler)
}
