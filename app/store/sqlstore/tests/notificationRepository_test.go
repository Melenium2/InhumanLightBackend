package sqlstore_test

import (
	"context"
	"testing"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestNotificationRepository_Create(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("notifications")
	
	store := sqlstore.New(db)
	notification := models.NewTestNotification(t)
	assert.NoError(t, store.Notifications(context.Background()).Create(notification))
	assert.NotEmpty(t, notification.ID)
}

func TestNotificationRepository_FindById(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	ctx := context.Background()
	defer cleaner("notifications")

	var userId uint = 3
	count := 5
	store := sqlstore.New(db)
	for i := 0; i < count; i++ {
		notification := models.NewTestNotification(t)
		assert.NoError(t, store.Notifications(ctx).Create(notification))
	} 
	notifs, err := store.Notifications(ctx).FindById(userId)
	assert.NoError(t, err)
	assert.NotNil(t, notifs)
	assert.Equal(t, count, len(notifs))
}

func TestNotificationRepository_Check(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	ctx := context.Background()
	defer cleaner("notifications")

	count := 5
	store := sqlstore.New(db)
	for i := 0; i < count; i++ {
		notification := models.NewTestNotification(t)
		assert.NoError(t, store.Notifications(ctx).Create(notification))
	} 
	
	var userId uint = 3
	notifs, _ := store.Notifications(ctx).FindById(userId)
	indexes := make([]int, 0)
	for _, item := range notifs {
		indexes = append(indexes, item.ID)
	}
	assert.NoError(t, store.Notifications(ctx).Check(indexes, userId))

	notifs, _ = store.Notifications(ctx).FindById(userId)
	for _, item := range notifs {
		assert.True(t, item.Checked)
	}
}