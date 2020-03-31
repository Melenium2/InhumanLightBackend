package sqlstore_test

import (
	"testing"

	"github.com/inhumanLightBackend/app/store/sqlstore"
	"github.com/magiconair/properties/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("users")
	
	assert.Equal(t, "1", "1")
}