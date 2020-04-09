package sqlstore_test

import (
	"testing"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/models/ticketStatus"
	"github.com/inhumanLightBackend/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestTicketRepository_Create(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("tickets")

	store := sqlstore.New(db)
	ticket := models.NewTestTicket(t)

	assert.NoError(t, store.Tickets().Create(ticket))
} 


func TestTicketRepository_Find(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("tickets")

	store := sqlstore.New(db)
	ticket := models.NewTestTicket(t)
	assert.NoError(t, store.Tickets().Create(ticket))
	ticket1, err := store.Tickets().Find(ticket.ID)
	assert.NoError(t, err)
	assert.NotNil(t, ticket1)
}

func TestTicketRepository_FindAll(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("tickets")

	store := sqlstore.New(db)
	ticketsCount := 3
	for i := 0; i < ticketsCount; i++ {
		ticket := models.NewTestTicket(t)
		assert.NoError(t, store.Tickets().Create(ticket))
	}
	tickets, err := store.Tickets().FindAll(uint(33))
	assert.NoError(t, err)
	assert.NotNil(t, tickets)
	assert.Equal(t, len(tickets), ticketsCount)
}

func TestTicketRepository_Accept(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("tickets")

	store := sqlstore.New(db)
	ticket := models.NewTestTicket(t)
	assert.NoError(t, store.Tickets().Create(ticket))
	helper := models.NewTestUser(t)

	assert.NoError(t, store.Tickets().Accept(ticket.ID, helper))
	ticket1, err := store.Tickets().Find(ticket.ID)
	assert.NoError(t, err)
	assert.NotNil(t, ticket1)

	assert.NotEqual(t, ticket1.Helper, -1)
	assert.Equal(t, ticket1.Status, ticketStatus.InProcess)
}

func TestTicketRepository_ChangeStatus(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("tickets")

	store := sqlstore.New(db)
	ticket := models.NewTestTicket(t)
	assert.NoError(t, store.Tickets().Create(ticket))
	assert.NoError(t, store.Tickets().ChangeStatus(ticket.ID, ticketStatus.Closed))
	ticket1, err := store.Tickets().Find(ticket.ID)
	assert.NoError(t, err)
	assert.NotNil(t, ticket1)

	assert.Equal(t, ticket1.Status, ticketStatus.Closed)
}

func TestTicketRepository_AddMessage(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("ticket_messages")

	store := sqlstore.New(db)
	message := models.NewTestTicketMessage(t)
	assert.NoError(t, store.Tickets().AddMessage(message))
	assert.NotEmpty(t, message.ID)
}

func TestTicketRepository_TakeMessages(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("ticket_messages")

	store := sqlstore.New(db)
	
	messagesCount := 30
	message := models.NewTestTicketMessage(t)
	for i := 0; i < messagesCount; i++ {
		assert.NoError(t, store.Tickets().AddMessage(message))
	}

	messages, err := store.Tickets().TakeMessages(message.TicketId)
	assert.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Equal(t, len(messages), messagesCount)
}


