package models

import (
	"testing"
	"time"
)

func NewTestUser(t *testing.T) *User {
	return &User{
		Email: "testUser@gmail.com",
		Login: "Usernmae",
		Password: "123456",
		Contacts: "Contacts",
		CreatedAt: time.Now(),
		IsActive: true,
		Role: Roles[0],
		Token: "supermegatoken",
	}
}

func NewInactiveUser(t *testing.T) *User {
	return &User{
		Email: "testUser@gmail.com",
		Login: "Usernmae",
		Password: "123456",
		Contacts: "Contacts",
		CreatedAt: time.Now(),
		IsActive: false,
		Role: Roles[0],
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