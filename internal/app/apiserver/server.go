package apiserver

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	models "github.com/Oringik/nyan404-libs/tree/master/models"
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
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
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
