package repository

import (
	"log"
	"testing"
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestTransaction_BuyItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := NewTransactionPostgres(db)

	userId := 1
	nameItem := "hoody"

	type mockBehavior func(m sqlmock.Sqlmock, userId int, nameItem string)

	testTable := []struct {
		name        string
		mockBehavior mockBehavior
		userCoins int
		priceItem int
		wantErr    bool
	}{
		{
			name: "OK",
			userCoins: 1000,
			mockBehavior: func(m sqlmock.Sqlmock, userId int, nameItem string) {
				m.ExpectBegin()
				m.ExpectQuery(`SELECT coins FROM wallet`).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000))
				m.ExpectQuery(`SELECT price FROM item`).
					WithArgs(nameItem).
					WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(300))
				m.ExpectExec(`UPDATE wallet SET coins = coins - \$1 WHERE user_id = \$2`).
					WithArgs(300, userId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectQuery(`SELECT id FROM item`).
					WithArgs(nameItem).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				m.ExpectExec(`INSERT INTO purchase`).
					WithArgs(userId, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				m.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Not Enough Coins",
			userCoins: 200,
			priceItem: 300,
			mockBehavior: func(m sqlmock.Sqlmock, userId int, nameItem string) {
				m.ExpectBegin()
				m.ExpectQuery(`SELECT coins FROM wallet`).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(200))
				m.ExpectQuery(`SELECT price FROM item`).
					WithArgs(nameItem).
					WillReturnRows(sqlmock.NewRows([]string{"price"}).AddRow(300))
				m.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Error Getting Item Price",
			userCoins: 1000,
			priceItem: 500,
			mockBehavior: func(m sqlmock.Sqlmock, userId int, nameItem string) {
				m.ExpectBegin()
				m.ExpectQuery(`SELECT coins FROM wallet`).
					WithArgs(userId).
					WillReturnRows(sqlmock.NewRows([]string{"coins"}).AddRow(1000))
				m.ExpectQuery(`SELECT price FROM item`).
					WithArgs(nameItem).
					WillReturnError(errors.New("item not found"))
				m.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mock, userId, nameItem)

			err := r.BuyItem(userId, nameItem)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuyItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}