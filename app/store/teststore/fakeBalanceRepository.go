package teststore

import "github.com/inhumanLightBackend/app/models"

type FakeBalanceRepository struct {
	store *Store
	balances map[int]*models.Balance
}

func (repo *FakeBalanceRepository) CreateBalance(userId uint) error {
	return nil
}

func (repo *FakeBalanceRepository) Add(userId uint, value float32) (*models.Balance, error) {
	return nil, nil
}

func (repo *FakeBalanceRepository) Remove(userId uint, value float32) (*models.Balance, error) {
	return nil, nil
}

func (repo *FakeBalanceRepository) LookForBalance(userId uint) (*models.Balance, error) {
	return nil, nil
}

func (repo *FakeBalanceRepository) AllTransactions(userId uint) ([]models.Balance, error) {
	return nil, nil
}