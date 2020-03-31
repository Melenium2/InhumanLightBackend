package models

import "time"

type User struct {
	ID                int        `json:"id"`
	Login             string     `json:"login"`
	Email             string     `json:"email"`
	Password          string     `json:"password,omitempty"`
	EncryptedPassword string     `json:"-"`
	CreatedAt         *time.Time `json:"registration_date"`
	Token             string     `json:"api_token"`
	Contacts          string     `json:"contacts"`
	Role              string     `json:"user_role"`
	IsActive          bool       `json:"-"`
}

var (
	roles = []string{"USER", "ADMIN"}
)

func changePassword(newPassword string) error {
	return nil
}

func changeToken() {

}

func changeActiveStatus(newStatus bool) {

}
