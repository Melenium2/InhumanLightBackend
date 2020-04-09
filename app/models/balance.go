package models

import "time"

// Balance model
type Balance struct {
	ID          uint      `json:"id"`
	Transaction float32   `json:"transaction"`
	BalanceNow  float32   `json:"balance_now"`
	From        string    `json:"from"`
	Date        time.Time `json:"date"`
	AddInfo     string    `json:"additional_info"`
	User        uint      `json:"user_id"`
}

// Init new instance of balance
func CreateBalance() *Balance {
	return &Balance{
		Transaction: 0,
		BalanceNow: 0,
		From: "Service",
		Date: time.Now().UTC(),
		AddInfo: "Init balance account",
	}
}
