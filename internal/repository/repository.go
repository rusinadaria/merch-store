package repository

import (
	"database/sql"
	"merch-store/models"
)

type Authorization interface {
	CreateUser(user models.AuthRequest) (int, error)
	AddCoins(user_wallet models.Wallet) error
	GetUserPassword(username string) (string, error)
	GetUserID(username string) (int, error)
}

type Transaction interface {
	SendCoin(id int, req models.SendCoinRequest) error
	BuyItem(userId int, name string) error
}

type Info interface {
	GetUserInfo(userId int) (models.InfoResponse, error)
}

type Repository struct {
	Authorization
	Transaction
	Info
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Transaction: NewTransactionPostgres(db),
		Info: NewInfoPostgres(db),
	}
}