package teststore

import (
	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store"
)

type Store struct {
	userRepository         *FakeUserRepository
	balanceRepository      *FakeBalanceRepository
	ticketRepository       *FakeTicketRepository
	notificationRepository *FakeNotificationRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &FakeUserRepository{
		store: s,
		users: make(map[int]*models.User),
	}

	return s.userRepository
}

func (s *Store) Balance() store.BalanceRepository {
	if s.balanceRepository != nil {
		return s.balanceRepository
	}

	s.balanceRepository = &FakeBalanceRepository{
		store:    s,
		balances: make(map[int]*models.Balance),
	}

	return s.balanceRepository
}

func (s *Store) Tickets() store.TicketRepository {
	if s.ticketRepository != nil {
		return s.ticketRepository
	}

	s.ticketRepository = &FakeTicketRepository{
		store:          s,
		tickets:        make(map[int]*models.Ticket),
		ticketMessages: make(map[int]*models.TicketMessage),
	}

	return s.ticketRepository
}

func (s *Store) Notifications() store.NotificationRepository {
	if s.notificationRepository != nil {
		return s.notificationRepository
	}

	s.notificationRepository = &FakeNotificationRepository{
		store:         s,
		notifications: make(map[int]*models.Notification),
	}

	return s.notificationRepository
}
