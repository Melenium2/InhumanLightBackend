package teststore

import (
	"context"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store"
)

type FakeUserRepository struct {
	store *Store
	ctx context.Context
	users map[int]*models.User
}

func (repo *FakeUserRepository) Create(newUser *models.User) error {
	if err := newUser.Validate(); err != nil {
		return err
	}

	if err := newUser.BeforeCreate(); err != nil {
		return err
	}

	newUser.ID = len(repo.users) + 1
	repo.users[newUser.ID] = newUser

	return nil
}

func (repo *FakeUserRepository) FindByEmail(email string) (*models.User, error) {
	for _, u := range repo.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound
}
func (repo *FakeUserRepository) FindById(id int) (*models.User, error) {
	user, ok := repo.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return user, nil
}
func (repo *FakeUserRepository) Update(user *models.User) error {
	_, err := repo.FindById(user.ID)
	if err != nil {
		return err
	}
	repo.users[user.ID] = user

	return nil
}