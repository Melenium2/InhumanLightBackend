package sqlstore

import (
	"github.com/inhumanLightBackend/app/models"
)

type BalanceRepository struct {
	store *Store
}

func (repo *BalanceRepository) CreateBalance(userId uint) error {
	balance := models.CreateBalance()
	balance.User = userId

	return repo.store.db.QueryRow(
		`insert into balance (transaction_value, balance_now, from_market, transaction_at, additional_info, user_id)
		 values ($1, $2, $3, $4, $5, $6) returning id`,
		 balance.Transaction,
		 balance.BalanceNow,
		 balance.From,
		 balance.Date,
		 balance.AddInfo,
		 balance.User,
	).Scan(&balance.ID)
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