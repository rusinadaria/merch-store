package repository

import (
	"merch-store/models"
	"database/sql"
	"log"
)

type InfoPostgres struct {
	db *sql.DB
}

func NewInfoPostgres(db *sql.DB) *InfoPostgres {
	return &InfoPostgres{db: db}
}

func (n *InfoPostgres) GetUserInfo(userId int) (models.InfoResponse, error) {
	var info models.InfoResponse
	var inventory models.Inventory
	var coinHistory models.CoinHistory
	var sent models.Sent
	var received models.Received

	rows, err := n.db.Query("SELECT to_user, amount FROM transaction WHERE from_user = $1", userId)
	if err != nil {
		log.Println("Пользователь не переводил коины другим сотрудникам:", err)
		return info, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&sent.ToUser, &sent.Amount); err != nil {
			log.Println("Ошибка при сканировании отправленных монет:", err)
			return info, err
		}
		coinHistory.Sent = append(coinHistory.Sent, sent)
	}

	rows, err = n.db.Query("SELECT from_user, amount FROM transaction WHERE to_user = $1", userId)
	if err != nil {
		log.Println("Этому сотруднику не переводили коины:", err)
		return info, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&received.FromUser, &received.Amount); err != nil {
			log.Println("Ошибка при сканировании полученных монет:", err)
			return info, err
		}
		coinHistory.Received = append(coinHistory.Received, received)
	}
	info.CoinHistory = append(info.CoinHistory, coinHistory)

	rows, err = n.db.Query(`SELECT i.name AS item_name, COUNT(p.item_id) AS quantity FROM purchase p
						JOIN item i ON p.item_id = i.id
						WHERE p.user_id = $1
						GROUP BY i.name;`, userId)
	if err != nil {
		log.Println("Этот сотрудник не покупал мерч:", err)
		return info, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&inventory.Type, &inventory.Quantity); err != nil {
			log.Println("Ошибка при сканировании инвентаря:", err)
			return info, err
		}
		info.Inventory = append(info.Inventory, inventory)
	}

	err = n.db.QueryRow("SELECT coins FROM wallet WHERE user_id = $1", userId).Scan(&info.Coins)
	if err != nil {
		log.Println("Этот сотрудник не покупал мерч:", err)
		return info, err
	}

	return info, nil
}
