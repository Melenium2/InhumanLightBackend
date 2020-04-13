package sqlstore

import (
	"context"
	"database/sql"

	"github.com/inhumanLightBackend/app/store"
	_ "github.com/lib/pq" //
)

// Store struct
type Store struct {
	db                     *sql.DB
	userRepository         *UserRepository
	balanceRepository      *BalanceRepository
	ticketRepository       *TicketRepository
	notificationRepositroy *NotificationRepository
}

// Create new store
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Return user functionality
func (store *Store) User(ctx context.Context) store.UserRepository {
	// Передовать контекст как параметр и класть в репозиторий. А дальше прописать у всех запросов
	if store.userRepository == nil {
		store.userRepository = &UserRepository{
			store: store,
			ctx:   ctx,
		}
	}

	return store.userRepository
}

// Return Balance transaction history functionality
func (store *Store) Balance(ctx context.Context) store.BalanceRepository {
	if store.balanceRepository == nil {
		store.balanceRepository = &BalanceRepository{
			store: store,
			ctx:   ctx,
		}
	}

	return store.balanceRepository
}

// Return Ticket history functionality
func (store *Store) Tickets(ctx context.Context) store.TicketRepository {
	if store.ticketRepository == nil {
		store.ticketRepository = &TicketRepository{
			store: store,
			ctx:   ctx,
		}
	}

	return store.ticketRepository
}

// Return Notification functionality
func (store *Store) Notifications(ctx context.Context) store.NotificationRepository {
	if store.notificationRepositroy == nil {
		store.notificationRepositroy = &NotificationRepository{
			store: store,
			ctx:   ctx,
		}
	}

	return store.notificationRepositroy
}
