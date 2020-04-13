package teststore

import (
	"context"

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

func (s *Store) User(ctx context.Context) store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &FakeUserRepository{
		store: s,
		ctx:   ctx,
		users: make(map[int]*models.User),
	}

	return s.userRepository
}

func (s *Store) Balance(ctx context.Context) store.BalanceRepository {
	if s.balanceRepository != nil {
		return s.balanceRepository
	}

	s.balanceRepository = &FakeBalanceRepository{
		store:    s,
		ctx:      ctx,
		balances: make(map[int]*models.Balance),
	}

	return s.balanceRepository
}

func (s *Store) Tickets(ctx context.Context) store.TicketRepository {
	if s.ticketRepository != nil {
		return s.ticketRepository
	}

	s.ticketRepository = &FakeTicketRepository{
		store:          s,
		ctx:            ctx,
		tickets:        make(map[int]*models.Ticket),
		ticketMessages: make(map[int]*models.TicketMessage),
	}

	return s.ticketRepository
}

func (s *Store) Notifications(ctx context.Context) store.NotificationRepository {
	if s.notificationRepository != nil {
		return s.notificationRepository
	}

	s.notificationRepository = &FakeNotificationRepository{
		store:         s,
		ctx:           ctx,
		notifications: make(map[int]*models.Notification),
	}

	return s.notificationRepository
}
