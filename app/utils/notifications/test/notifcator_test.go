package notifications_test

import (
	"os"
	"strconv"
	"testing"
)

var (
	apiToken string
	userId   int
)

func TestMain(m *testing.M) {
	apiToken = os.Getenv("TELEGRAM_API_TOKEN")
	if apiToken == "" {
		apiToken = "1293039613:AAER81Qqklo9JZQa3kt2iHrKBA9ptPpJ8IY"
	}
	userIdStr := os.Getenv("TELEGRAM_USER_ID")
	if userIdStr == "" {
		userId = 708015155
	} else {
		userId, _ = strconv.Atoi(userIdStr)
	}

	os.Exit(m.Run())
}
