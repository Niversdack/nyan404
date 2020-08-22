package apiserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Oringik/nyan404-libs/database"
	models "github.com/Oringik/nyan404-libs/models"
	"github.com/gorilla/websocket"

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

func (s *server) init() {
	playersArray := []*models.User{
		{
			ID: 1,
			FIuser: models.FIuser{
				Name:    "Дмитрий",
				Surname: "Степанов",
			},
		},
		{
			ID: 2,
			FIuser: models.FIuser{
				Name:    "Василий",
				Surname: "Тапочкин",
			},
		},
		{
			ID: 3,
			FIuser: models.FIuser{
				Name:    "Артем",
				Surname: "Пупкин",
			},
		},
		{
			ID: 4,
			FIuser: models.FIuser{
				Name:    "Елезавета",
				Surname: "Орлова",
			},
		},
		{
			ID: 5,
			FIuser: models.FIuser{
				Name:    "Екатерина",
				Surname: "Коврова",
			},
		},
		{
			ID: 6,
			FIuser: models.FIuser{
				Name:    "Адрей",
				Surname: "Данник",
			},
		},
		{
			ID: 7,
			FIuser: models.FIuser{
				Name:    "Дарья",
				Surname: "Дворцова",
			},
		},
		{
			ID: 8,
			FIuser: models.FIuser{
				Name:    "Евгений",
				Surname: "Мартынов",
			},
		},
		{
			ID: 9,
			FIuser: models.FIuser{
				Name:    "Мария",
				Surname: "Вишнева",
			},
		},
		{
			ID: 10,
			FIuser: models.FIuser{
				Name:    "Виктория",
				Surname: "Халтурина",
			},
		},
	}
	s.db.SetArray(playersArray)
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
	// s.router.HandleFunc("/getCards", s.handleGetCards())
	// s.router.HandleFunc("/sendAnswer", s.handleSendAnswer())
	s.router.HandleFunc("/setUser", s.handleSetUser())
}

// func generateId() int {
// 	return 1
// }

// func (s *server) handleGetCards() http.HandlerFunc {
// 	var userCase *models.UserCase
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		err := db.Model(userCase).Field(userCase.ID).Equal(generateId()).Get()
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			log.Println(err)
// 		}
// 		card, err = json.Marshal(&userCase)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			log.Println(err)
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		NewResponseWriter(card, w)
// 	}
// }

// func (s *server) handleSendAnswer() http.HandlerFunc {
// 	type request struct {
// 		Answer string `json:"answer"`
// 	}
// 	var req request
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		body := NewRequestReader(r)
// 		json.Unmarshal(body, &req)
// 		err := db.Model(req).Set()
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			log.Println(err)
// 		}
// 		card, err = json.Marshal()
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			log.Println(err)
// 		}
// 		w.WriteHeader(http.StatusOK)
// 		NewResponseWriter(card, w)
// 	}
// }

func (s *server) handleGetUserCase() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *server) handleSetUser() http.HandlerFunc {
	ID := 1
	return func(w http.ResponseWriter, r *http.Request) {
		var p models.User
		_, err := s.db.Model(&p).Field("ID").Equal(ID).Get()
		player, err := json.Marshal(p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		ID++
		NewResponseWriter(player, w)
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
		//msg := []byte("Let's start to talk something.")
		//
		//err = conn.WriteMessage(websocket.TextMessage, msg)
		//if err != nil {
		//	log.Println(err)
		//}
	}
	// do other stuff...
}
