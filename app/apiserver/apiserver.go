package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/inhumanLightBackend/app/store/sqlstore"
)

// Start and configure server
func Start(config *Config) error {
	db, err := newDb(config.DatabaseURL)
	if err != nil {
		return err
	}

	store := sqlstore.New(db)
	s := NewServer(store)
	
	return http.ListenAndServe(config.Port, s)
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