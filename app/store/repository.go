package store

import "github.com/inhumanLightBackend/app/models"

// UserRepository ...
type UserRepository interface {
	Create(*models.User) error
	FindByEmail(string) (*models.User, error)
	FindById(int) (*models.User, error)
	Update(*models.User) error
}