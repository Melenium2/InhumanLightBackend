package models

import "time"

type TicketMessage struct {
	ID       uint      `json:"id"`
	Who      uint      `json:"who"`
	TicketId uint      `json:"ticket_id"`
	Message  string    `json:"message"`
	Date     time.Time `json:"date"`
}

func (tm *TicketMessage) BeforeCreate() {
	tm.Date = time.Now().UTC()
}
