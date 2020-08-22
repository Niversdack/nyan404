package apiserver

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
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
		{
			UserInfo: models.UserInfo{
				Name:    "Мария",
				Surname: "Ерошкина",
				Gender:  "Женщина",
				Age:     30,
				Job:     models.JOB_SECRETARY,
			},
			Cases: []*models.Case{
				{
					ID:       2,
					AnswerID: 0,
					Description: models.Description{
						Title: "Женщина и Покупки",
						Text:  "Вы стоите за прилавком в крупном супермаркете и видите женщину,которая стоит недалеко от кассы и пытается что-то сделать в телефоне.Она кричит что-то в трубку и жалуется на проблемы со связью.Она хочет дозвониться до мужа,чтобы тот помог ей с покупками.Ваши действия...",
					},
					Ans: []models.Answer{
						{
							ID:           1,
							Text:         "Спросить,что произошло и в чем проблема",
							Significance: 3,
						},
						{
							ID:           2,
							Text:         "Предложить сходу купить новую симку",
							Significance: -3,
						},
						{
							ID:           3,
							Text:         "Пожаловаться охране на сумасшедшую женщину",
							Significance: -5,
						},
						{
							ID:           4,
							Text:         "Предложить ей свой телефон для звонка",
							Significance: 4,
						},
					},
				},
			},
		},
		{
			UserInfo: models.UserInfo{
				Name:    "Мария",
				Surname: "Ерошкина",
				Gender:  "Женщина",
				Age:     30,
				Job:     models.JOB_SECRETARY,
			},
			Cases: []*models.Case{
				{
					ID:       3,
					AnswerID: 1,
					Description: models.Description{
						Title: "Женщина и Покупки",
						Text:  "Женщина рассказала,что не может позвонить мужу,так как сеть не ловит.Ваши действия...",
					},
					Ans: []models.Answer{
						{
							ID:           1,
							Text:         "Рассказать о своих услугах",
							Significance: 4,
						},
						{
							ID:           2,
							Text:         "Предложить перезагрузить телефон",
							Significance: 1,
						},
						{
							ID:           3,
							Text:         "Сказать,что на самом деле это заговор того оператора",
							Significance: -2,
						},
						{
							ID:           4,
							Text:         "Уйти и сказать,что сами разбирайтесь с проблемами",
							Significance: -5,
						},
					},
				},
			},
		},
		{
			UserInfo: models.UserInfo{
				Name:    "Мария",
				Surname: "Ерошкина",
				Gender:  "Женщина",
				Age:     30,
				Job:     models.JOB_SECRETARY,
			},
			Cases: []*models.Case{
				{
					ID:       4,
					AnswerID: 2,
					Description: models.Description{
						Title: "Женщина и Покупки",
						Text:  "Женщина разозлилась и сказала,что ваш оператор сотовой связи будет таким же...",
					},
					Ans: []models.Answer{
						{
							ID:           1,
							Text:         "Накричать и сказать,что нет",
							Significance: -5,
						},
						{
							ID:           2,
							Text:         "Рассказать о преимуществах тарифа",
							Significance: 3,
						},
						{
							ID:           3,
							Text:         "Задобрить и сделать комплимент внешности",
							Significance: -1,
						},
						{
							ID:           4,
							Text:         "Предложить конфетку",
							Significance: -3,
						},
					},
				},
			},
		},
		{
			Cases: []*models.Case{
				{
					ID:       5,
					AnswerID: 3,
					Description: models.Description{
						Title: "Женщина и Покупки",
						Text:  "Охранник подошел к женщине и начал узнавать,в чем проблема.Затем они подошли вдвоем и сказали,что вы были некомпетентны и просят извинений.",
					},
					Ans: []models.Answer{
						{
							ID:           1,
							Text:         "Извинитесь и попытаетесь прояснить ситуацию",
							Significance: 2,
						},
						{
							ID:           2,
							Text:         "Уйдете и крикните,что вы всех засудите",
							Significance: -5,
						},
						{
							ID:           3,
							Text:         "Сделаете глупое лицо и скажете,что это не вы",
							Significance: -2,
						},
						{
							ID:           4,
							Text:         "Расскажите о том,что хотели обратить внимание на себя",
							Significance: -1,
						},
					},
				},
			},
		},
		{
			Cases: []*models.Case{
				{
					ID:       6,
					AnswerID: 4,
					Description: models.Description{
						Title: "Женщина и Покупки",
						Text:  "Женщина приняла предложение.После быстрого звонка она спросила,как она может вас отблагодарить...",
					},
					Ans: []models.Answer{
						{
							ID:           1,
							Text:         "Попросите ее выслушать вас насчет смены оператора",
							Significance: 3,
						},
						{
							ID:           2,
							Text:         "Скажите,что это безвозмездно",
							Significance: -1,
						},
						{
							ID:           3,
							Text:         "Предложите сразу купить симку",
							Significance: -2,
						},
						{
							ID:           4,
							Text:         "Попросить деньги за звонок",
							Significance: -4,
						},
					},
				},
			},
		},
		{
			Cases: []*models.Case{
				{
					ID:       7,
					AnswerID: 1,
					Description: models.Description{
						Title: "Женщина и Покупки",
						Text:  "Женщина вас выслушивает и выглядит заинтересованной...",
					},
					Ans: []models.Answer{
						{
							ID:           1,
							Text:         "Предложить оформить симку с новым номером",
							Significance: 1,
						},
						{
							ID:           2,
							Text:         "Рассказать еще о преимуществах",
							Significance: 4,
						},
						{
							ID:           3,
							Text:         "Уйти,так как покупатель все равно не купит ничего",
							Significance: -5,
						},
						{
							ID:           4,
							Text:         "Оформить симку с переносом старого номера",
							Significance: 4,
						},
					},
				},
			},
		},
		{
			Cases: []*models.Case{
				{
					ID:       8,
					AnswerID: 2,
					Description: models.Description{
						Title: "Женщина и Покупки",
						Text:  "Женщина перезагрузила телефон,но ничего не изменилось",
					},
					Ans: []models.Answer{
						{
							ID:           1,
							Text:         "Покачать головой и сказать,что ничем не помочь",
							Significance: -3,
						},
						{
							ID:           2,
							Text:         "Отобрать телфон и сказать,что вы делаете не правильно",
							Significance: -5,
						},
						{
							ID:           3,
							Text:         "Предложить оформить симку",
							Significance: 2,
						},
						{
							ID:           4,
							Text:         "Уйти,сказав,что попросите у кого-нибудь телефон,чтобы позвонить",
							Significance: -2,
						},
					},
				},
			},
		},
		{},
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
		userCases, err := s.db.Model(&models.UserCase{}).GetArray()
		if err != nil {
			Response(err.Error(), w, http.StatusInternalServerError)
			return
		}
		randomIndex := rand.Intn(len(userCases))
		pick := userCases[randomIndex].(*models.UserCase)

		Reponse(pick, w, http.StatusOK)
		return

	}
}

func (s *server) handleSetUser() http.HandlerFunc {
	ID := 1
	return func(w http.ResponseWriter, r *http.Request) {
		var p *models.User
		value, err := s.db.Model(p).Field("ID").Equal(ID).Get()
		player, err := json.Marshal(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		if ID == 10 {
			ID = 1
		} else {
			ID++
		}
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
