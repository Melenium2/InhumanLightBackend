package sqlstore

import (
	"database/sql"
	"sort"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store"
)

type TicketRepository struct {
	store *Store
}

func (repo *TicketRepository) Create(ticket *models.Ticket) error {
	ticket.BeforeCreate()

	return repo.store.db.QueryRow(
		`insert into tickets (title, description, section, from_user, helper, created_at, status) 
		values ($1, $2, $3, $4, $5, $6, $7) returning id`,
		ticket.Title,
		ticket.Description,
		ticket.Section,
		ticket.From,
		ticket.Helper,
		ticket.Created_at,
		ticket.Status,
	).Scan(&ticket.ID)
}

func (repo *TicketRepository) Accept(ticketId uint, helper *models.User) error {
	_, err := repo.store.db.Exec(
		"update tickets set helper = $2, status = $3 where id = $1",
		ticketId,
		helper.ID,
		models.TicketProcessStatus[1],
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *TicketRepository) Find(ticketId uint) (*models.Ticket, error) {
	ticket := &models.Ticket{}

	if err := repo.store.db.QueryRow(
		"select id, title, description, section, from_user, helper, created_at, status from tickets where id = $1",
		ticketId,
	).Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Section, 
	&ticket.From, &ticket.Helper, &ticket.Created_at, &ticket.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return ticket, nil
}

func (repo *TicketRepository) ChangeStatus(ticketId uint, status string) error {
	index := sort.SearchStrings(models.TicketProcessStatus, status)
	if index < 0 {
		return store.ErrProccessingStatusNotFound
	}

	_, err := repo.store.db.Exec(
		"update tickets set status = $2 where id = $1",
		ticketId,
		status,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *TicketRepository) AddMessage(tm *models.TicketMessage) error {
	tm.BeforeCreate()

	return repo.store.db.QueryRow(
		"insert into ticket_messages (who, ticket_id, message_text, reply_at) values ($1, $2, $3, $4) returning id" ,
		tm.Who,
		tm.TicketId,
		tm.Message,
		tm.Date,
	).Scan(&tm.ID)
}

func (repo *TicketRepository) TakeMessages(ticketId uint) ([]*models.TicketMessage, error) {
	rows, err :=repo.store.db.Query(
		"select * from ticket_messages where ticket_id = $1",
		ticketId,
	)
	if err != nil {
		return nil, err
	}
	var messages []*models.TicketMessage
	for rows.Next() {	
		message := &models.TicketMessage{}
		if err := rows.Scan(&message.ID, &message.Who, &message.TicketId, &message.Message, &message.Date); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}