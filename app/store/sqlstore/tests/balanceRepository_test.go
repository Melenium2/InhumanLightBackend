package sqlstore_test

import (
	"context"
	"testing"

	"github.com/inhumanLightBackend/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestBalanceRepository_Create(t *testing.T) {
	db, cleaner := sqlstore.TestDb(t, databaseUrl)
	defer cleaner("user", "balance")

	store := sqlstore.New(db)
	assert.NoError(t, store.Balance(context.Background()).CreateBalance(23))
}