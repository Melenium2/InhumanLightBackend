package apiserver

import (
	"database/sql"
	"net/http"
	"os"
	"os/signal"

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

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	go func() {
		<-exit
		s.notificationChannel <- "Server shutdown"
		println("Server shutdown")
		close(s.notificationChannel)

		os.Exit(1)
	}()

	s.notificationChannel <- "Server started"
	println("Api server started. Telegram bot sended " + <-s.notificationChannel)
	
	if err := http.ListenAndServe(config.Port, s); err != nil {
		s.notificationChannel <- "Server offline with error " + err.Error()
		<-s.notificationChannel
		close(s.notificationChannel)
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