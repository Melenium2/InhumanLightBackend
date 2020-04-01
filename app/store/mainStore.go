package store

type Store interface {
	User() UserRepository
	Balance() BalanceRepository
	Tickets() TicketRepository
}