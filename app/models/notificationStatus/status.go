package notificationStatus

import "errors"

var (
	ErrTicketStatusNotFound = errors.New("Status not found")	
)

const (
	Info = "info"
	Warnign = "warning"
	Error = "error"
)