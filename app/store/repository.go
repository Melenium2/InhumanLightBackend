package store

import "github.com/inhumanLightBackend/app/models"

// UserRepository ...
type UserRepository interface {
	Create(*models.User) error
	FindByEmail(string) (*models.User, error)
	FindById(int) (*models.User, error)
	Update(*models.User) error
}

type BalanceRepository interface {
	CreateBalance(uint) error
	AllTransactions(uint) ([]models.Balance, error)
	LookForBalance(uint) (*models.Balance, error)
	Add(uint, float32) (*models.Balance, error)
	Remove(uint, float32) (*models.Balance, error)
}