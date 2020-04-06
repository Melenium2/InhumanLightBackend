package teststore_test

import (
	"testing"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T)  {
	store := teststore.New()
	user := models.NewTestUser(t)
	assert.NoError(t, store.User().Create(user))
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	store := teststore.New()
	user := models.NewTestUser(t)
	assert.NoError(t, store.User().Create(user))
	user1, err := store.User().FindByEmail(user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user1)
}

func TestUserRepository_FindById(t *testing.T) {
	store := teststore.New()
	user := models.NewTestUser(t)
	assert.NoError(t, store.User().Create(user))
	user1, err := store.User().FindById(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user1)
}

func TestUserRepository_Update(t *testing.T) {
	store := teststore.New()
	user := models.NewTestUser(t)
	assert.NoError(t, store.User().Create(user))
	email := "supermegamen@gmail.com"
	user.Email = email
	assert.NoError(t, store.User().Update(user))
}