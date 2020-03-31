package sqlstore

import "github.com/inhumanLightBackend/app/models"

type UserRepository struct {
	store *Store
}

func (repo *UserRepository) Create(newUser *models.User) error {
	
	return nil
}
func (reop *UserRepository) FindByEmail(email string) (*models.User, error) {
	return nil, nil
}
func (repo *UserRepository)	FundById(id int) (*models.User, error) {
	return nil, nil
}
func (repo *UserRepository) Update(user *models.User) error {
	return nil
}