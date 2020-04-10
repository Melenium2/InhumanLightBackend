package notifications_test

import (
	"testing"

	"github.com/inhumanLightBackend/app/utils/notifications/telegram"
	"github.com/stretchr/testify/assert"
)

func TestTelegramNotificator_Notify(t *testing.T) {
	notificatior := telegram.New(userId, apiToken)
	mc := notificatior.Notify()
	mc <- "Message from telegram"
	assert.Equal(t, <-mc, "200 OK")
}