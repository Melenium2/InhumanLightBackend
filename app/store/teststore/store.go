package teststore

import (
	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store"
)

type Store struct {
	userRepository *UserRepository 
	balanceRepository *BalanceRepository
	ticketRepository *TicketRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*models.User),
	}

	return s.userRepository
}

func (s *Store) Balance() store.BalanceRepository {
	if s.balanceRepository != nil {
		return s.balanceRepository
	}

	s.balanceRepository = &BalanceRepository{
		store: s,
		balances: make(map[int]*models.Balance),
	}

	return s.balanceRepository
}

func (s *Store) Tickets() store.TicketRepository {
	if s.ticketRepository != nil {
		return s.ticketRepository
	}

	s.ticketRepository = &TicketRepository{
		store: s,
		tickets: make(map[int]*models.Ticket),
		ticketMessages: make(map[int]*models.TicketMessage),
	}

	return s.ticketRepository
}