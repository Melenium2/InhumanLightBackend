package store

import "errors"

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("Record not found")
	// ErrProccessingStatusNotFound
	ErrProccessingStatusNotFound = errors.New("Proccessing status not found")
)