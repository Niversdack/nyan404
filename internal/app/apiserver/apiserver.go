package apiserver

import (
	"github.com/Oringik/nyan404-libs/database"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"net/http"
)

// Start ...
func Start(config *Config) error {

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(sessionStore)

	srv.logger.Info("Server starting")

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
