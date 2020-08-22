package apiserver

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"

	"github.com/Oringik/nyan404-libs/helpers"

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
	currentID    uint
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
	for _, player := range playersArray {
		s.db.Model(player).Set()

	}

	userCases := []*models.UserCase{
		{
			UserInfo: models.UserInfo{
				Name:    "Олег",
				Surname: "Ромашкин",
				Gender:  "Мужчина",
				Age:     21,
				Job:     models.JOB_POLICEMAN_BOY,
			},
			Cases: []*models.Case{
				{
					ID:       1,
					AnswerID: 0,
					Description: models.Description{
						Title: "Крикливый мужчина",
						Text:  "Вы встречаете мужчину на улицу, кричащего во весь голос. Из обрывков его фраз вы понимаете, что у него проблемы с сотовой связью. Окончательно прояснив для себя ситуацию, вы решаете подойти к молодому человеку и предложить свои услуги, но получаете резко-агрессивное настроение против вас...",
					},
					Ans: []models.Answer{
						{
							ID:           1,
							Text:         "Предложить мужчине успокоительное",
							Significance: -1,
						},
						{
							ID:           2,
							Text:         "Начать расспрашивать его о случившейся ситуации",
							Significance: 3,
						},
						{
							ID:           3,
							Text:         "Вызвать скорую помощь",
							Significance: -4,
						},
						{
							ID:           4,
							Text:         "Рассказать о своих услугах",
							Significance: -2,
						},
					},
				},
			},
		},
	}

	for _, userCase := range userCases {
		s.db.Model(userCase).Set()
	}
}

func newServer(sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		db:           database.NewModelStorage(),
		hub:          newHub(),
		sessionStore: sessionStore,
		currentID:    0,
	}
	s.init()
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
	s.router.HandleFunc("/getusercase", s.handleGetUserCase())
	s.router.HandleFunc("/inituser", s.handleInitUser())
	s.router.HandleFunc("/answer", s.handleSendAnswer())
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
	type request struct {
		UserID uint `json:"user_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		body := NewRequestReader(r)

		err := json.Unmarshal(body, req)
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}

		ansCounter := helpers.GetAnswerCounter()
		ansCounter.InitAnswerCounter()

		userCounter := &models.UserCounter{
			UserID:        req.UserID,
			AnswerCounter: ansCounter,
		}

		s.db.Model(userCounter).Set()

		userCases, err := s.db.Model(&models.UserCase{}).GetArray()
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}

		normallyUserCases := []*models.UserCase{}

		for _, userCase := range userCases.([]interface{}) {
			normallyUserCases = append(normallyUserCases, userCase.(*models.UserCase))
		}

		randomIndex := rand.Intn(len(normallyUserCases))
		pick := normallyUserCases[randomIndex]

		data, err := json.Marshal(pick)
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}

		Response(data, w, http.StatusOK)
		return

	}
}

func (s *server) handleInitUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := s.db.Model(&models.User{}).Field("ID").Equal(s.currentID).Get()
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}
		s.currentID++

		data, err := json.Marshal(user)
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}

		Response(data, w, http.StatusOK)
		return

	}
}

func (s *server) getUserCaseByAnswer() http.HandlerFunc {
	type request struct {
		UserCaseID uint `json:"user_case_id"`
		CaseID     uint `json:"case_id"`
		AnswerID   uint `json:"answer_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		body := NewRequestReader(r)

		err := json.Unmarshal(body, req)
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}

		userCase, err := s.db.Model(&models.UserCase{}).Field("ID").Equal(req.UserCaseID).Get()
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}

		normallyUserCase := userCase.(*models.UserCase)

		for _, singleCase := range normallyUserCase.Cases {
			if singleCase.ID == req.CaseID {
				if singleCase.AnswerID == req.AnswerID {
					data, err := json.Marshal(singleCase)
					if err != nil {
						Response([]byte(err.Error()), w, http.StatusInternalServerError)
						return
					}

					Response(data, w, http.StatusOK)
					return
				}
			}
		}

		Response([]byte("Value not found"), w, http.StatusBadRequest)
		return

	}
}

func (s *server) handleSendAnswer() http.HandlerFunc {
	type request struct {
		UserID     uint `json:"user_id"`
		UserCaseID uint `json:"user_case_id"`
		CaseID     uint `json:"case_id"`
		AnswerID   uint `json:"answer_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		body := NewRequestReader(r)

		err := json.Unmarshal(body, req)
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}

		userCase, err := s.db.Model(&models.UserCase{}).Field("ID").Equal(req.UserCaseID).Get()
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}

		normallyUserCase := userCase.(*models.UserCase)
		var answer *models.Answer

		for _, singleCase := range normallyUserCase.Cases {
			if singleCase.ID == req.CaseID {
				for _, ans := range singleCase.Ans {
					if ans.ID == req.AnswerID {
						answer = &ans
					}
				}
			}
		}

		userCounter, err := s.db.Model(&models.UserCounter{}).Field("UserID").Equal(req.UserID).Get()
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}

		normallyUserCounter := userCounter.(*models.UserCounter).AnswerCounter

		offset, err := normallyUserCounter.GenerateOffset(answer.Significance)
		if err != nil {
			Response([]byte(err.Error()), w, http.StatusInternalServerError)
			return
		}

		normallyUserCounter.RecountBalance(offset)

		if normallyUserCounter.KindOfBeyond() == helpers.FAIL {
			Response([]byte("FAIL"), w, http.StatusOK)
			return
		}

		if normallyUserCounter.KindOfBeyond() == helpers.SUCCESS {
			Response([]byte("SUCCESS"), w, http.StatusOK)
			return
		}

		Response([]byte("NOPE"), w, http.StatusOK)
		return

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
