package teststore

import (
	"errors"

	"github.com/inhumanLightBackend/app/models"
)

type FakeTicketRepository struct {
	store *Store
	tickets map[int]*models.Ticket
	ticketMessages map[int]*models.TicketMessage
}

func (repo *FakeTicketRepository) Create(ticket *models.Ticket) error {
	ticket.BeforeCreate()

	mapCount := len(repo.tickets) + 1
	ticket.ID = uint(mapCount)
	repo.tickets[mapCount] = ticket

	return nil
}

func (repo *FakeTicketRepository) Accept(ticketId uint, helper *models.User) error {
	repo.tickets[int(ticketId)].Helper = helper.ID

	return nil
}

func (repo *FakeTicketRepository) Find(ticketId uint) (*models.Ticket, error) {
	ticket, ok := repo.tickets[int(ticketId)]
	if !ok {
		return nil, errors.New("Not found")
	} 
	return ticket, nil
}

func (repo *FakeTicketRepository) FindAll(userId uint) ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	for _, item := range repo.tickets {
		if item.From == userId {
			tickets = append(tickets, item)
		}
	} 

	return tickets, nil
}

func (repo *FakeTicketRepository) ChangeStatus(ticketId uint, status string) error {
	repo.tickets[int(ticketId)].Status = status

	return nil
}

func (repo *FakeTicketRepository) AddMessage(tm *models.TicketMessage) error {
	tm.BeforeCreate()

	nextId := len(repo.ticketMessages) + 1
	tm.ID = uint(nextId)
	repo.ticketMessages[nextId] = tm

	return nil
}

func (repo *FakeTicketRepository) TakeMessages(ticketId uint) ([]*models.TicketMessage, error) {
	messages := make([]*models.TicketMessage, 0)
	for _, item := range repo.ticketMessages {
		if item.TicketId == ticketId {
			messages = append(messages, item)
		}
	}
	
	return messages, nil
}