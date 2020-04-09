package teststore_test

import (
	"testing"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/models/ticketStatus"
	"github.com/inhumanLightBackend/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestFakeTicketRepository_Craete(t *testing.T) {
	ticket := models.NewTestTicket(t)
	store := teststore.New()
	assert.NoError(t, store.Tickets().Create(ticket))
}

func TestFakeTicketRepository_Accept(t *testing.T) {
	ticket := models.NewTestTicket(t)
	user := models.NewTestUser(t)
	store := teststore.New()
	assert.NoError(t, store.Tickets().Create(ticket))
	assert.NoError(t, store.Tickets().Accept(ticket.ID, user))
}

func TestFakeTicketRepository_Find(t *testing.T) {
	ticket := models.NewTestTicket(t)
	store := teststore.New()
	assert.NoError(t, store.Tickets().Create(ticket))
	ti, err := store.Tickets().Find(ticket.ID)
	assert.NoError(t, err)
	assert.NotNil(t,ti)
}

func TestFakeTicketRepository_FindAll(t *testing.T) {
	var userId uint = 33
	ticketCount := 5
	store := teststore.New()
	for i := 0; i < ticketCount; i++ {
		ticket := models.NewTestTicket(t)
		store.Tickets().Create(ticket)
	}
	tickets, err := store.Tickets().FindAll(userId)
	assert.NoError(t, err)
	assert.NotNil(t, tickets)
	assert.Equal(t, ticketCount, len(tickets))
}

func TestFakeTicketRepository_ChangeStatus(t *testing.T) {
	ticket := models.NewTestTicket(t)
	store := teststore.New()
	store.Tickets().Create(ticket)
	assert.NoError(t, store.Tickets().ChangeStatus(ticket.ID, ticketStatus.InProcess))
}

func TestFakeTicketRepository_AddMessage(t *testing.T) {
	message := models.NewTestTicketMessage(t)
	store := teststore.New()
	assert.NoError(t, store.Tickets().AddMessage(message))
}

func TestFakeTicketRepository_TakeMessages(t *testing.T) {
	store := teststore.New()
	messgesCount := 5
	var ticketId uint = 1
	for i := 0; i < messgesCount; i++ {
		message := models.NewTestTicketMessage(t)
		store.Tickets().AddMessage(message)
	}
	messages, err := store.Tickets().TakeMessages(ticketId)
	assert.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Equal(t, messgesCount, len(messages))
}