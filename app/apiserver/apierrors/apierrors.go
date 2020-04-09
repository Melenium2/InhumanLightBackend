package apierrors

import "errors"

var (
	ErrEmptyParam               = errors.New("Invalid param")
	ErrNotValidBody             = errors.New("Invalid json")
	ErrIncorrectEmailOrPassword = errors.New("Incorrect email or password")
	ErrNotAuthenticated         = errors.New("Not authenticated")
	ErrPermissionDenied         = errors.New("Permission denied")
)
