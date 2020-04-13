package teststore_test

import (
	"context"
	"testing"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestFakeUserRepository_Create(t *testing.T)  {
	store := teststore.New()
	user := models.NewTestUser(t)
	assert.NoError(t, store.User(context.Background()).Create(user))
	assert.NotNil(t, user)
}

func TestFakeUserRepository_FindByEmail(t *testing.T) {
	store := teststore.New()
	ctx := context.Background()
	user := models.NewTestUser(t)
	assert.NoError(t, store.User(ctx).Create(user))
	user1, err := store.User(ctx).FindByEmail(user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user1)
}

func TestFakeUserRepository_FindById(t *testing.T) {
	store := teststore.New()
	ctx := context.Background()
	user := models.NewTestUser(t)
	assert.NoError(t, store.User(ctx).Create(user))
	user1, err := store.User(ctx).FindById(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user1)
}

func TestFakeUserRepository_Update(t *testing.T) {
	store := teststore.New()
	ctx := context.Background()
	user := models.NewTestUser(t)
	assert.NoError(t, store.User(ctx).Create(user))
	email := "supermegamen@gmail.com"
	user.Email = email
	assert.NoError(t, store.User(ctx).Update(user))
}