package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/inhumanLightBackend/app/store/sqlstore"
	"github.com/inhumanLightBackend/app/utils/notifications/telegram"
)

// Start and configure server
func Start(config *Config) error {
	db, err := newDb(config.DatabaseURL)
	if err != nil {
		return err
	}

	store := sqlstore.New(db)
	notificator := telegram.New(config.TelegramUserId, config.TelegramToken)
	s := NewServer(store, notificator)
	s.mc <- "Api server started"
	println("Api server started. Telegram bot sended " + <-s.mc)
	if err := http.ListenAndServe(config.Port, s); err != nil {
		s.mc <- err.Error()
		<-s.mc
		return err
	}

	return nil
}

// Init db connection
func newDb(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}