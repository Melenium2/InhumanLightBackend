package store

// Main store with all repos
type Store interface {
	User() UserRepository
	Balance() BalanceRepository
	Tickets() TicketRepository
}