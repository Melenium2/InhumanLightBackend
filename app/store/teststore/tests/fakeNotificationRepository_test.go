package teststore_test

import (
	"context"
	"testing"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestFakeNotificationRepository_Create(t *testing.T) {
	store := teststore.New()
	notification := models.NewTestNotification(t)
	assert.NoError(t, store.Notifications(context.Background()).Create(notification))
}

func TestFakeNotificationRepository_FindById(t *testing.T) {
	store := teststore.New()
	ctx := context.Background()
	count := 5
	var userId uint = 3

	for i := 0; i < count; i++ {
		notif := models.NewTestNotification(t)
		store.Notifications(ctx).Create(notif)
	}

	notifs, err := store.Notifications(ctx).FindById(userId)
	assert.NoError(t, err)
	assert.NotNil(t, notifs)
	assert.Equal(t, count, len(notifs))
}

func TestFakeNotificationRepository_Check(t *testing.T) {
	store := teststore.New()
	ctx := context.Background()
	count := 5
	var userId uint = 3

	for i := 0; i < count; i++ {
		notif := models.NewTestNotification(t)
		store.Notifications(ctx).Create(notif)
	}
	indexes := []int{1, 2, 3, 4, 5}
	
	assert.NoError(t, store.Notifications(ctx).Check(indexes, userId))
	notifs, _ := store.Notifications(ctx).FindById(userId)
	for _, item := range notifs {
		assert.True(t, item.Checked)
	}
}