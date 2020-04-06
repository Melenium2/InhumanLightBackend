package teststore

import "github.com/inhumanLightBackend/app/models"

type BalanceRepository struct {
	store *Store
	balances map[int]*models.Balance
}

func (repo *BalanceRepository) CreateBalance(userId uint) error {
	return nil
}

func (repo *BalanceRepository) Add(userId uint, value float32) (*models.Balance, error) {
	return nil, nil
}

func (repo *BalanceRepository) Remove(userId uint, value float32) (*models.Balance, error) {
	return nil, nil
}

func (repo *BalanceRepository) LookForBalance(userId uint) (*models.Balance, error) {
	return nil, nil
}

func (repo *BalanceRepository) AllTransactions(userId uint) ([]models.Balance, error) {
	return nil, nil
}