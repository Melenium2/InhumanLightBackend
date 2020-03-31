package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseUrl string
)

func TestMain(m *testing.M) {
	databaseUrl = os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "postgres://postgres:123456@localhost/go_test?sslmode=disable"
	}

	os.Exit(m.Run())
}