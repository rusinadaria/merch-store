package services

import (
	"merch-store/internal/repository"
	"merch-store/models"
)

type TransactionService struct {
	repo repository.Transaction
}

func NewTransactionService(repo repository.Transaction) *TransactionService {
	return &TransactionService{repo: repo}
}

func (t *TransactionService) SendCoin(senderId int, username string, amount int) error {
	req := models.SendCoinRequest {
		ToUser: username,
		Amount: amount,
	}

	err := t.repo.SendCoin(senderId, req)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionService) BuyItem(userId int, name string) error {
	err := t.repo.BuyItem(userId, name)
	if err != nil {
		return err
	}
	return nil
}