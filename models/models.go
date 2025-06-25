package models

import (
	"time"
)

type AuthRequest struct {
	Id int 
	Username string `json:"username"`
	Password string `json:"password"`
}

type Wallet struct {
	Id int
	User_id int
	Coins int
}

type Item struct {
	Id int
	Name string
	Price int
}

type ErrorResponse  struct {
	Errors string `json:"errors"`
}

type Inventory struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"` // Количество предметов
}

type Sent struct {
	ToUser string `json:"toUser"`
	Amount int `json:"amount"`
}

type Received struct {
	FromUser string `json:"fromUser"`
	Amount int `json:"amount"`
}

type CoinHistory struct {
	Received []Received `json:"received"`
	Sent     []Sent `json:"sent"`
}

type InfoResponse struct {
	Coins int `json:"coins"`
	Inventory []Inventory `json:"inventory"`
	CoinHistory []CoinHistory `json:"coinHistory"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type Transaction struct {
	ID        int       `json:"id"`
	FromUser  int       `json:"from_user"`
	ToUser    int       `json:"to_user"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}