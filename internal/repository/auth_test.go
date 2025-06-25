package repository

import (
	"merch-store/models"
	"log"
	"testing"
	"fmt"
	"database/sql"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
)


func TestAuth_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	user := models.AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}

	type mockBehavior func(m sqlmock.Sqlmock, user models.AuthRequest)

	testTable := []struct {
		name        string
		mockBehavior mockBehavior
		user       models.AuthRequest
		wantID     int
		wantErr    bool
	}{
		{
			name: "OK",
			user: user,
			wantID: 2,
			mockBehavior: func(m sqlmock.Sqlmock, user models.AuthRequest) {
				m.ExpectQuery(`INSERT INTO "user"`).
					WithArgs(user.Username, user.Password).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
			},
		},
		{
			name: "Error on insert",
			user: user,
			wantID: 0,
			wantErr: true,
			mockBehavior: func(m sqlmock.Sqlmock, user models.AuthRequest) {
				m.ExpectQuery(`INSERT INTO "user"`).
					WithArgs(user.Username, user.Password).
					WillReturnError(fmt.Errorf("some error"))
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mock, tt.user)

			gotID, err := r.CreateUser(tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("CreateUser() gotID = %v, want %v", gotID, tt.wantID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAuth_AddCoins(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	userWallet := models.Wallet{
		User_id: 1,
	}

	type mockBehavior func(m sqlmock.Sqlmock, userWallet models.Wallet)

	testTable := []struct {
		name        string
		mockBehavior mockBehavior
		wallet       models.Wallet
		wantErr    bool
	}{
		{
			name: "OK",
			wallet: userWallet,
			mockBehavior: func(m sqlmock.Sqlmock, userWallet models.Wallet) {
				m.ExpectExec(`INSERT INTO wallet`).
					WithArgs(userWallet.User_id, 1000).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "Empty User ID",
			wallet: models.Wallet{
				User_id: 0,
			},
			mockBehavior: func(m sqlmock.Sqlmock, userWallet models.Wallet) {
			},
			wantErr: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mock, tt.wallet)

			err := r.AddCoins(tt.wallet)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddCoins() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAuth_GetUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	username := "testuser"

	type mockBehavior func(m sqlmock.Sqlmock, username string)

	testTable := []struct {
		name        string
		mockBehavior mockBehavior
		wantID int
		wantErr    bool
	}{
		{
			name: "OK",
			mockBehavior: func(m sqlmock.Sqlmock, username string) {
				m.ExpectQuery(`SELECT id FROM "user"`).
					WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantID: 1,
			wantErr: false,
		},
		{
			name: "USER NOT FOUND",
			mockBehavior: func(m sqlmock.Sqlmock, username string) {
				m.ExpectQuery(`SELECT id FROM "user"`).
					WithArgs(username).
					WillReturnError(sql.ErrNoRows)
			},
			wantID: 0,
			wantErr: true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mock, username)

			id, err := r.GetUserID(username)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if id != tt.wantID {
				t.Errorf("GetUserID() gotID = %v, wantID %v", id, tt.wantID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAuth_GetUserPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	r := NewAuthPostgres(db)

	username := "testuser"

	type mockBehavior func(m sqlmock.Sqlmock, username string)

	testTable := []struct {
		name        string
		mockBehavior mockBehavior
		wantPasswd string
		wantErr    bool
	}{
		{
			name: "OK",
			mockBehavior: func(m sqlmock.Sqlmock, username string) {
				m.ExpectQuery(`SELECT password FROM "user"`).
					WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("qwerty"))
			},
			wantPasswd: "qwerty",
			wantErr: false,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mock, username)

			passwd, err := r.GetUserPassword(username)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if passwd != tt.wantPasswd {
				t.Errorf("GetUserPassword() gotID = %v, wantID %v", passwd, tt.wantPasswd)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}