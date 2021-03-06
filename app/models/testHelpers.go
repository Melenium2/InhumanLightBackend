package models

import (
	"testing"
	"time"

	"github.com/inhumanLightBackend/app/models/notificationStatus"
	"github.com/inhumanLightBackend/app/models/roles"
	"github.com/inhumanLightBackend/app/models/ticketStatus"
)

func NewTestUser(t *testing.T) *User {
	return &User{
		Email: "testUser@gmail.com",
		Login: "Usernmae",
		Password: "123456",
		Contacts: "Contacts",
		CreatedAt: time.Now(),
		IsActive: true,
		Role: roles.USER,
		Token: "supermegatoken",
	}
}

func NewTestInactiveUser(t *testing.T) *User {
	return &User{
		Email: "testUser@gmail.com",
		Login: "Usernmae",
		Password: "123456",
		Contacts: "Contacts",
		CreatedAt: time.Now(),
		IsActive: false,
		Role: roles.USER,
		Token: "supermegatoken",
	}
}

func NewTestUserEmptyFields(t *testing.T) *User {
	return &User{
		Email: "testUser@gmail.com",
		Login: "Usernmae",
		Password: "123456",
		Contacts: "Contacts",
	}
}

func NewTestBalanceEmpty(t *testing.T) *Balance {
	return &Balance{
		Transaction: 0,
		BalanceNow: 0,
		Date: time.Now().UTC(),
		From: "Bank",
		User: 1,
		AddInfo: "Kakoyto balance",
	}
}

func NewTestTicket(t *testing.T) *Ticket {
	return &Ticket{
		Title: "Zagolovok",
		Description: "Kkakoyto description",
		Section: "Super question",
		From: 33,
		Helper: -1,
		Created_at: time.Now().UTC(),
		Status: ticketStatus.Opened,
	}
}

func NewTestTicketMessage(t *testing.T) *TicketMessage {
	return &TicketMessage{
		Who: 43,
		TicketId: 1,
		Message: "Message",
		Date: time.Now().UTC(),
	}
}

func NewTestNotification(t *testing.T) *Notification {
	return &Notification{
		Message: "Message",
		Status: notificationStatus.Info,
		For: 3,
	}
}