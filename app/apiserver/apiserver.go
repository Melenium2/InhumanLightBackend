package apiserver

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"time"

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
	notifs := telegram.New(config.TelegramUserId, config.TelegramToken).Notify()
	s := NewServer(store, config)
	notifs <- "Server started"

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	go func() {
		<-exit
		notifs <- "Server shutdown"
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			println("Server shutdown with error")
		}
		println("Server shutdown")
		close(notifs)

		os.Exit(1)
	}()

	println(fmt.Sprintf("Api server started on port %s", config.Port))
	println("Telegram bot sent " + <-notifs)
	
	if err := s.ListenAndServe(); err != nil {
		notifs <- "Server drops with error " + err.Error()
		<-notifs
		close(notifs)
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