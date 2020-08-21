package apiserver

import (
	"encoding/json"
	"errors"
	"github.com/Oringik/nyan404-libs/database"
	"github.com/Oringik/nyan404-libs/models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nyan404/internal/app/store"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "fastexp"
	ctxKeyUser  ctxKey = iota
)

var (
	errIncorrectEmailOrPassword = errors.New("Incorrect email or password")
	errNotAuthenticated         = errors.New("Not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	db           *database.ModelStorage
	hub          *Hub
	sessionStore sessions.Store
}

func newServer(sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		db:           database.NewModelStorage(),
		hub:          newHub(),
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {

	s.router.HandleFunc("/ws", s.serveWs())
	s.router.HandleFunc("/getCards", s.handleGetCards())
	s.router.HandleFunc("/sendAnswer", s.handleSendAnswer())
}

func generateId() int {
	return 1
}

func (s *server) handleGetCards() http.HandlerFunc {
	var userCase *models.UserCase
	return func(w http.ResponseWriter, r *http.Request) {
		err := db.Model(userCase).Field(userCase.ID).Equal(generateId()).Get()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}
		card, err = json.Marshal(&userCase)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		NewResponseWriter(card, w)
	}
}

func (s *server) handleSendAnswer() http.HandlerFunc {
	type request struct {
		Answer string `json:"answer"`
	}
	var req request
	return func(w http.ResponseWriter, r *http.Request) {
		body := NewRequestReader(r)
		json.Unmarshal(body, &req)
		err := db.Model(req).Set()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}
		card, err = json.Marshal()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		NewResponseWriter(card, w)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (s *server) serveWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//func(s *server) handler(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		client := &Client{hub: s.hub, conn: conn, send: make(chan []byte, 256)}
		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()
		msg := []byte("Let's start to talk something.")

		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println(err)
		}
	}
	// do other stuff...
}
