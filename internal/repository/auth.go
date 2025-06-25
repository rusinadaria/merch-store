package repository

import (
	// "github.com/jmoiron/sqlx"
	"merch-store/models"
	"database/sql"
	"log"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) GetUserPassword(username string) (string, error) {
	var password string
	err := r.db.QueryRow(`SELECT password FROM "user" WHERE username = $1`, username).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (r *AuthPostgres) GetUserID(username string) (int, error) {
	var id int
	err := r.db.QueryRow(`SELECT id FROM "user" WHERE username = $1`, username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) CreateUser(user models.AuthRequest) (int, error) {
	query := `INSERT INTO "user" (username, password) VALUES ($1, $2) RETURNING id`
	var id int
	err := r.db.QueryRow(query, user.Username, user.Password).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) AddCoins(user_wallet models.Wallet) error {
	coins := 1000
	query := `INSERT INTO wallet (user_id, coins) VALUES ($1, $2)`
	_, err := r.db.Exec(query, user_wallet.User_id, coins)
	if err != nil {
		return err
	}
	return nil
}