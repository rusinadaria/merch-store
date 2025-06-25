package services

import (
	"merch-store/internal/repository"
	"merch-store/models"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Auth interface {
	CreateUser(username string, password string) (int, error)
	GenerateToken(userId int) (string, error)
	ParseToken(accessToken string) (int, error)
	FindUser(username string) (int, error)
	SignIn(username string, password string) (int, error)
}

type Transaction interface {
	SendCoin(senderId int, username string, amount int) error
	BuyItem(userId int, name string) error
}

type Info interface {
	UserInfo(userId int) (models.InfoResponse, error)
}

type Service struct {
	Auth
	Transaction
	Info
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repos.Authorization),
		Transaction: NewTransactionService(repos.Transaction),
		Info: NewInfoService(repos.Info),
	}
}
