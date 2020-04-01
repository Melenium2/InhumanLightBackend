package sqlstore_test

import (
	"testing"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("users")

	store := sqlstore.New(db)
	user := models.NewTestUserEmptyFields(t)
	
	assert.NoError(t, store.User().Create(user))
	assert.NotNil(t, user.ID)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("users")

	store := sqlstore.New(db)
	user := models.NewTestUser(t)
	_, err := store.User().FindByEmail(user.Email)
	assert.Error(t, err)

	store.User().Create(user)
	user, err = store.User().FindByEmail(user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindById(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("users")

	store := sqlstore.New(db)
	user1 := models.NewTestUser(t)
	store.User().Create(user1)
	user2, err := store.User().FindById(user1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user2)
}

func TestUserRepository_Update(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("users")

	store := sqlstore.New(db)
	user := models.NewTestUser(t)
	assert.NoError(t, store.User().Create(user))
	newLogin := "Vasya"
	newContacts := "From UK"
	user.Login = newLogin
	user.Contacts = newContacts

	assert.NoError(t, store.User().Update(user))
	user1, err := store.User().FindByEmail(user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user1)
	assert.Equal(t, user1.Login, newLogin)
	assert.Equal(t, user1.Contacts, newContacts)
}