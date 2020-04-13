package teststore_test

import (
	"context"
	"testing"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/models/ticketStatus"
	"github.com/inhumanLightBackend/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestFakeTicketRepository_Craete(t *testing.T) {
	ticket := models.NewTestTicket(t)
	store := teststore.New()
	assert.NoError(t, store.Tickets(context.Background()).Create(ticket))
}

func TestFakeTicketRepository_Accept(t *testing.T) {
	ticket := models.NewTestTicket(t)
	ctx := context.Background()
	user := models.NewTestUser(t)
	store := teststore.New()
	assert.NoError(t, store.Tickets(ctx).Create(ticket))
	assert.NoError(t, store.Tickets(ctx).Accept(ticket.ID, user))
}

func TestFakeTicketRepository_Find(t *testing.T) {
	ticket := models.NewTestTicket(t)
	ctx := context.Background()
	store := teststore.New()
	assert.NoError(t, store.Tickets(ctx).Create(ticket))
	ti, err := store.Tickets(ctx).Find(ticket.ID)
	assert.NoError(t, err)
	assert.NotNil(t,ti)
}

func TestFakeTicketRepository_FindAll(t *testing.T) {
	ctx := context.Background()
	var userId uint = 33
	ticketCount := 5
	store := teststore.New()
	for i := 0; i < ticketCount; i++ {
		ticket := models.NewTestTicket(t)
		store.Tickets(ctx).Create(ticket)
	}
	tickets, err := store.Tickets(ctx).FindAll(userId)
	assert.NoError(t, err)
	assert.NotNil(t, tickets)
	assert.Equal(t, ticketCount, len(tickets))
}

func TestFakeTicketRepository_ChangeStatus(t *testing.T) {
	ticket := models.NewTestTicket(t)
	ctx := context.Background()
	store := teststore.New()
	store.Tickets(ctx).Create(ticket)
	assert.NoError(t, store.Tickets(ctx).ChangeStatus(ticket.ID, ticketStatus.InProcess))
}

func TestFakeTicketRepository_AddMessage(t *testing.T) {
	message := models.NewTestTicketMessage(t)
	store := teststore.New()
	assert.NoError(t, store.Tickets(context.Background()).AddMessage(message))
}

func TestFakeTicketRepository_TakeMessages(t *testing.T) {
	store := teststore.New()
	ctx := context.Background()
	messgesCount := 5
	var ticketId uint = 1
	for i := 0; i < messgesCount; i++ {
		message := models.NewTestTicketMessage(t)
		store.Tickets(ctx).AddMessage(message)
	}
	messages, err := store.Tickets(ctx).TakeMessages(ticketId)
	assert.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Equal(t, messgesCount, len(messages))
}