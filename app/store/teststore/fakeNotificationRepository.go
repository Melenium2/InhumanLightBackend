package teststore

import (
	"context"

	"github.com/inhumanLightBackend/app/models"
)

type FakeNotificationRepository struct {
	store         *Store
	ctx           context.Context
	notifications map[int]*models.Notification
}

func (repo *FakeNotificationRepository) Create(newModel *models.Notification) error {
	nextId := len(repo.notifications) + 1
	newModel.ID = nextId
	repo.notifications[nextId] = newModel

	return nil
}

func (repo *FakeNotificationRepository) FindById(userId uint) ([]*models.Notification, error) {
	notifications := make([]*models.Notification, 0)
	for _, item := range repo.notifications {
		if item.For == int(userId) {
			notifications = append(notifications, item)
		}
	}

	return notifications, nil
}

func (repo *FakeNotificationRepository) Check(notifications []int, userId uint) error {
	for _, n := range notifications {
		for _, item := range repo.notifications {
			if n == item.ID && int(userId) == item.For {
				item.Checked = true
			}
		}
	}

	return nil
}
