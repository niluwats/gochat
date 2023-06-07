package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/niluwats/gochat/pkg/dto"
	"github.com/niluwats/gochat/pkg/middleware"
	"github.com/niluwats/gochat/pkg/service"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GOCHAT http server")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &dto.UserReq{}

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoding request object", http.StatusBadRequest)
		return
	}

	log.Println(u)

	res := service.Register(u)
	json.NewEncoder(w).Encode(res)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &dto.UserReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoding request object", http.StatusBadRequest)
		return
	}

	res := service.Login(u)
	if res.Status {
		token, err := middleware.GenerateJWT(u.Username)
		if err != nil {
			res.Message = err.Error()
			res.Status = false
		}
		res.Data = token
	}
	json.NewEncoder(w).Encode(res)
}

func verifyContactHandler(w http.ResponseWriter, r *http.Request) {
	if middleware.VerifyJWT(w, r) != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	u := &dto.UserReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoding request object", http.StatusBadRequest)
		return
	}

	res := service.VerifyContact(u.Username)
	json.NewEncoder(w).Encode(res)
}

func chatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if middleware.VerifyJWT(w, r) != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	u1 := r.URL.Query().Get("u1")
	u2 := r.URL.Query().Get("u2")

	fromTS, toTS := "0", "+inf"

	if r.URL.Query().Get("from-ts") != "" && r.URL.Query().Get("to-ts") != "" {
		fromTS = r.URL.Query().Get("from-ts")
		toTS = r.URL.Query().Get("to-ts")
	}

	res := service.ChatHistory(u1, u2, fromTS, toTS)
	json.NewEncoder(w).Encode(res)
}

func contactListHandler(w http.ResponseWriter, r *http.Request) {
	if middleware.VerifyJWT(w, r) != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	u := r.URL.Query().Get("username")
	res := service.ContactList(u)
	json.NewEncoder(w).Encode(res)
}
