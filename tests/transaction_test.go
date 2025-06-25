package tests

import (
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"merch-store/models"
	"bytes"
	"github.com/stretchr/testify/assert"
)

func (s *APITestSuite) TestcreateTestUsers() {
	sender := models.AuthRequest{
		Username: "sender_username",
		Password: "12345",
	}

	requestBody, err := json.Marshal(sender)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(requestBody))
	assert.NoError(s.T(), err)

	req.AddCookie(&http.Cookie{Name: "auth_token", Value: s.token})


	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.AddUserHandler)

	handler.ServeHTTP(rr, req)
	assert.Equal(s.T(), http.StatusOK, rr.Code)
}

func (s *APITestSuite) TestcreateTestTwoUsers() {
	recipient := models.AuthRequest{
		Username: "recipient_username",
		Password: "12345",
	}

	requestBody, err := json.Marshal(recipient)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(requestBody))
	assert.NoError(s.T(), err)

	req.AddCookie(&http.Cookie{Name: "auth_token", Value: s.token_two_user})


	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.AddUserHandler)

	handler.ServeHTTP(rr, req)
	assert.Equal(s.T(), http.StatusOK, rr.Code)
	
}


func (s *APITestSuite) TestSendCoins() {
	sendCoinRequest := models.SendCoinRequest{
		ToUser: "recipient_username",
		Amount: 100,
	}

	requestBody, err := json.Marshal(sendCoinRequest)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(requestBody))
	assert.NoError(s.T(), err)

	req.AddCookie(&http.Cookie{Name: "auth_token", Value: s.token})


	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.SendHandler)


	handler.ServeHTTP(rr, req)


	assert.Equal(s.T(), http.StatusOK, rr.Code)

	s.checkCoinBalances("sender_username", "recipient_username", 100)
}


func (s *APITestSuite) checkCoinBalances(senderUsername, recipientUsername string, amount int) {
	senderBalance, err := s.getUserBalance(senderUsername)
	assert.NoError(s.T(), err, "Failed to get sender balance")

	recipientBalance, err := s.getUserBalance(recipientUsername)
	assert.NoError(s.T(), err, "Failed to get recipient balance")

	assert.Equal(s.T(), 1000-amount, senderBalance, "Sender balance is incorrect")
	assert.Equal(s.T(), 1000+amount, recipientBalance, "Recipient balance is incorrect")
}

func (s *APITestSuite) getUserBalance(username string) (int, error) {
	var balance int
	query := `SELECT coins FROM wallet WHERE user_id = (SELECT id FROM "user" WHERE username = $1)`
	err := s.db.QueryRow(query, username).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (s *APITestSuite) TestBuyItem() {
	req, err := http.NewRequest("GET", "/api/buy/powerbank", nil)
	assert.NoError(s.T(), err)

	req.AddCookie(&http.Cookie{Name: "auth_token", Value: s.token})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.BuyItemHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
}
