package sqlstore

import (
	"database/sql"

	"github.com/inhumanLightBackend/app/store"
	_ "github.com/lib/pq" //
)

// Store struct
type Store struct {
	db *sql.DB
	userRepository *UserRepository
}

// Create new store
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (store *Store) User() store.UserRepository {
	if store.userRepository == nil {
		store.userRepository = &UserRepository{
			store: store,
		}
	}
	
	return store.userRepository
}