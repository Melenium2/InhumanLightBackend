package teststore

import (
	"errors"

	"github.com/inhumanLightBackend/app/models"
)

type TicketRepository struct {
	store *Store
	tickets map[int]*models.Ticket
	ticketMessages map[int]*models.TicketMessage
}

func (repo *TicketRepository) Create(ticket *models.Ticket) error {
	ticket.BeforeCreate()

	mapCount := len(repo.tickets) + 1
	ticket.ID = uint(mapCount)
	repo.tickets[mapCount] = ticket

	return nil
}

func (repo *TicketRepository) Accept(ticketId uint, helper *models.User) error {
	repo.tickets[int(ticketId)].Helper = helper.ID

	return nil
}

func (repo *TicketRepository) Find(ticketId uint) (*models.Ticket, error) {
	ticket, ok := repo.tickets[int(ticketId)]
	if !ok {
		return nil, errors.New("Not found")
	} 
	return ticket, nil
}

func (repo *TicketRepository) FindAll(userId uint) ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	for _, item := range repo.tickets {
		if item.From == userId {
			tickets = append(tickets, item)
		}
	} 

	return tickets, nil
}

func (repo *TicketRepository) ChangeStatus(ticketId uint, status string) error {
	repo.tickets[int(ticketId)].Status = status

	return nil
}

func (repo *TicketRepository) AddMessage(tm *models.TicketMessage) error {
	return nil
}

func (repo *TicketRepository) TakeMessages(ticketId uint) ([]*models.TicketMessage, error) {
	return nil, nil
}