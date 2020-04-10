package store

import "github.com/inhumanLightBackend/app/models"

// UserRepository
type UserRepository interface {
	Create(*models.User) error
	FindByEmail(string) (*models.User, error)
	FindById(int) (*models.User, error)
	Update(*models.User) error
}

// BalanceRepository
type BalanceRepository interface {
	CreateBalance(uint) error
	AllTransactions(uint) ([]models.Balance, error)
	LookForBalance(uint) (*models.Balance, error)
	Add(uint, float32) (*models.Balance, error)
	Remove(uint, float32) (*models.Balance, error)
}

// TicketRepository
type TicketRepository interface {
	Create(*models.Ticket) error
	Accept(uint, *models.User) error
	Find(uint) (*models.Ticket, error)
	FindAll(uint) ([]*models.Ticket, error)
	ChangeStatus(uint, string) error
	TicketMessagesRepository
}

// TicketMessagesRepository
type TicketMessagesRepository interface {
	AddMessage(*models.TicketMessage) error
	TakeMessages(uint) ([]*models.TicketMessage, error)
}

type NotificationRepository interface {
	Create(*models.Notification) error
	FindById(uint) ([]*models.Notification, error)
	Check([]int) error
}