package models

import (
	"time"

	"github.com/inhumanLightBackend/app/models/notificationStatus"
)

type Notification struct {
	ID      int       `json:"id"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
	Status  string    `json:"status"`
	For     int       `json:"for"`
	Checked bool      `json:"checked"`
}

func (n *Notification) Validate() error {
	switch n.Status {
	case notificationStatus.Info:
	case notificationStatus.Warnign:
	case notificationStatus.Error:
	default:
		return notificationStatus.ErrTicketStatusNotFound
	}

	return nil
}

func (n *Notification) BeforeCreate() {
	n.Date = time.Now().UTC()
	n.Checked = false
}
