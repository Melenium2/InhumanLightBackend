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
	balanceRepository *BalanceRepository
	ticketRepository *TicketRepository
}

// Create new store
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Return user functionality
func (store *Store) User() store.UserRepository {
	if store.userRepository == nil {
		store.userRepository = &UserRepository{
			store: store,
		}
	}
	
	return store.userRepository
}

// Return Balance transaction history functionality
func (store *Store) Balance() store.BalanceRepository {
	if store.balanceRepository == nil {
		store.balanceRepository = &BalanceRepository{
			store: store,
		}
	}

	return store.balanceRepository
}

// Return Ticket history functionality
func (store *Store) Tickets() store.TicketRepository {
	if store.ticketRepository == nil {
		store.ticketRepository = &TicketRepository{
			store: store,
		}
	}

	return store.ticketRepository
}