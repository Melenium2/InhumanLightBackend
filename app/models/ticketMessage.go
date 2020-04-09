package models

import "time"

// Ticket message model
type TicketMessage struct {
	ID       uint      `json:"id"`
	Who      uint      `json:"who"`
	TicketId uint      `json:"ticket_id"`
	Message  string    `json:"message"`
	Date     time.Time `json:"date"`
}

// Fill fields before message create
func (tm *TicketMessage) BeforeCreate() {
	tm.Date = time.Now().UTC()
}
