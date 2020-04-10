package telegram

import (
	"fmt"
	"net/http"
)

type (
	TelegramNotificator struct {
		UserId   int
		ApiToken string
	}
)

func New(userId int, token string) *TelegramNotificator {
	return &TelegramNotificator{
		UserId: userId,
		ApiToken: token,
	}
}

func (ta *TelegramNotificator) Notify() chan string {
	mc := make(chan string)
	
	go func(m chan string) {
		for {
			message := <-m
			resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&parse_mode=Markdown&text=%s", ta.ApiToken, ta.UserId, message))
			if err != nil {
				m <- err.Error()
			} else {
				m <- resp.Status
			}
		}
	}(mc)

	return mc
}
