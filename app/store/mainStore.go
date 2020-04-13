package store

import "context"

// Main store with all repos
type Store interface {
	User(ctx context.Context) UserRepository
	Balance(ctx context.Context) BalanceRepository
	Tickets(ctx context.Context) TicketRepository
	Notifications(ctx context.Context) NotificationRepository
}