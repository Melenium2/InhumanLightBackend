package teststore

import "github.com/inhumanLightBackend/app/models"

type TicketRepository struct {
	store *Store
	tickets map[int]*models.Ticket
	ticketMessages map[int]*models.TicketMessage
}

func (repo *TicketRepository) Create(ticket *models.Ticket) error {
	return nil
}

func (repo *TicketRepository) Accept(ticketId uint, helper *models.User) error {
	return nil
}

func (repo *TicketRepository) Find(ticketId uint) (*models.Ticket, error) {
	return nil, nil
}

func (repo *TicketRepository) ChangeStatus(ticketId uint, status string) error {
	return nil
}

func (repo *TicketRepository) AddMessage(tm *models.TicketMessage) error {
	return nil
}

func (repo *TicketRepository) TakeMessages(ticketId uint) ([]*models.TicketMessage, error) {
	return nil, nil
}