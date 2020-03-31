package sqlstore

import "github.com/inhumanLightBackend/app/models"

type UserRepository struct {
	store *Store
}

func (repo *UserRepository) Create(*models.User) error {
	return nil
}
func (reop *UserRepository) FindByEmail(string) (*models.User, error) {
	return nil, nil
}
func (repo *UserRepository)	FundById(int) (*models.User, error) {
	return nil, nil
}
func (repo *UserRepository) Update(*models.User) error {
	return nil
}