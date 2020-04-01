package models

import "time"

type Ticket struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Section     string    `json:"section"`
	From        uint      `json:"from"`
	Helper      int       `json:"helper"`
	Created_at  time.Time `json:"created_at"`
	Status      string    `json:"status"`
}

var (
	TicketProcessStatus = []string{"opened", "in process", "closed"}
)

func (t *Ticket) BeforeCreate() {
	t.Helper = -1
	t.Created_at = time.Now().UTC()
	t.Status = TicketProcessStatus[0]
}
