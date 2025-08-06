package repository

import (
	"merch-store/models"
	"database/sql"
	"log"
	"errors"
)

type TransactionPostgres struct {
	db *sql.DB
}

func NewTransactionPostgres(db *sql.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (t *TransactionPostgres) withTx(fn func(*sql.Tx) error) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func logTxError(context string, err error) error {
	log.Printf("%s: %v", context, err)
	return err
}

func (t *TransactionPostgres) SendCoin(fromUserId int, req models.SendCoinRequest) error {
	return t.withTx(func(tx *sql.Tx) error {
		var senderCoins int
		if err := tx.QueryRow("SELECT coins FROM wallet WHERE user_id = $1", fromUserId).Scan(&senderCoins); err != nil {
			return logTxError("ошибка при получении монет отправителя", err)
		}

		if senderCoins < req.Amount {
			return errors.New("недостаточно монет для отправки")
		}

		if _, err := tx.Exec("UPDATE wallet SET coins = coins - $1 WHERE user_id = $2", req.Amount, fromUserId); err != nil {
			return logTxError("ошибка при обновлении баланса отправителя", err)
		}

		var toUserId int
		if err := tx.QueryRow(`
			SELECT user_id FROM wallet
			WHERE user_id = (SELECT id FROM "user" WHERE username = $1)`, req.ToUser).Scan(&toUserId); err != nil {
			return logTxError("ошибка при получении ID получателя", err)
		}

		if _, err := tx.Exec("UPDATE wallet SET coins = coins + $1 WHERE user_id = $2", req.Amount, toUserId); err != nil {
			return logTxError("ошибка при обновлении баланса получателя", err)
		}

		if _, err := tx.Exec(`
			INSERT INTO transaction (from_user, to_user, amount)
			VALUES ($1, $2, $3)`, fromUserId, toUserId, req.Amount); err != nil {
			return logTxError("ошибка при записи транзакции", err)
		}

		return nil
	})
}


// func (t *TransactionPostgres) SendCoin(fromUserId int, req models.SendCoinRequest) error {
// 	tx, err := t.db.Begin()
// 	if err != nil {
// 		log.Println(err)
// 		log.Printf("Ошибка при получении монет отправителя: %v", err)
// 		return err
// 	}

// 	var senderCoins int
// 	err = tx.QueryRow("SELECT coins FROM wallet WHERE user_id = $1", fromUserId).Scan(&senderCoins)
// 	if err != nil {
// 		log.Println(err)
// 		log.Printf("Ошибка при получении монет отправителя: %v", err)
// 		tx.Rollback()
// 		return err
// 	}

// 	if senderCoins < req.Amount {
// 		tx.Rollback()
// 		log.Printf("Ошибка при получении монет отправителя: %v", err)
// 		return errors.New("недостаточно монет для отправки")
// 	}

// 	_, err = tx.Exec("UPDATE wallet SET coins = coins - $1 WHERE user_id = $2", req.Amount, fromUserId)
// 	if err != nil {
// 		log.Println(err)
// 		log.Printf("Ошибка при получении монет отправителя: %v", err)
// 		tx.Rollback()
// 		return err
// 	}

// 	var toUserId int
// 	err = tx.QueryRow(`SELECT user_id FROM wallet
// 						WHERE user_id = (SELECT id FROM "user" WHERE username = $1)`, req.ToUser).Scan(&toUserId)
// 	if err != nil {
// 		log.Println(err)
// 		log.Printf("Ошибка при получении ID получателя: %v", err)
// 		tx.Rollback()
// 		return err
// 	}

// 	_, err = tx.Exec("UPDATE wallet SET coins = coins + $1 WHERE user_id = $2", req.Amount, toUserId)
// 	if err != nil {
// 		log.Println(err)
// 		log.Printf("Ошибка при получении монет отправителя: %v", err)
// 		tx.Rollback()
// 		return err
// 	}

// 	_, err = tx.Exec("INSERT INTO transaction (from_user, to_user, amount) VALUES ($1, $2, $3)", fromUserId, toUserId, req.Amount)
// 	if err != nil {
// 		log.Println(err)
// 		log.Printf("Ошибка при получении монет отправителя: %v", err)
// 		tx.Rollback()
// 		return err
// 	}

// 	return tx.Commit()

// }

func (t *TransactionPostgres) BuyItem(userId int, name string) error {

	tx, err := t.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	var userCoins int
	err = tx.QueryRow("SELECT coins FROM wallet WHERE user_id = $1", userId).Scan(&userCoins)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	var priceItem int
	err = tx.QueryRow("SELECT price FROM item WHERE name = $1", name).Scan(&priceItem)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	if userCoins < priceItem {
		tx.Rollback()
		return errors.New("недостаточно монет для покупки товара")
	}

	_, err = tx.Exec("UPDATE wallet SET coins = coins - $1 WHERE user_id = $2", priceItem, userId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	var itemId int
	err = tx.QueryRow("SELECT id FROM item WHERE name = $1", name).Scan(&itemId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO purchase (user_id, item_id) VALUES ($1, $2)", userId, itemId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
