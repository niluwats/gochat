package httpserver

import (
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

	router.HandleFunc("/", rootHandler).Methods(http.MethodGet)
	router.HandleFunc("/register", registerHandler).Methods(http.MethodPost)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	router.HandleFunc("/verifycontact", verifyContactHandler).Methods(http.MethodPost)
	router.HandleFunc("/chathistory", chatHistoryHandler).Methods(http.MethodGet)
	router.HandleFunc("/contactlist", contactListHandler).Methods(http.MethodGet)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	http.ListenAndServe(":8080", handler)
}
