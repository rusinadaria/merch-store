package handlers

import (
	mock_services "merch-store/internal/services/mocks"
	"merch-store/internal/services"
	"bytes"
	"net/http/httptest"
	"testing"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"errors"
)

func TestHandler_BuyItemHandler(t *testing.T) {
	type mockBehavior func(s *mock_services.MockTransaction, userId int)

	testTable := []struct {
		name string
		inputBody string
		mockBehavior mockBehavior
		expectedStatusCode int
		expectedResponseBody string
	} {
		{
			name: "OK",
			inputBody: ``,
			mockBehavior: func(s *mock_services.MockTransaction, userId int) {
				s.EXPECT().BuyItem(userId, "hoody").Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: ``,
		},
		{
			name: "INCORRECT VALUE",
			inputBody: ``,
			mockBehavior: func(s *mock_services.MockTransaction, userId int) {
				s.EXPECT().BuyItem(userId, "incorect_value").Return(errors.New("Не возможно приобрести товар"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: `{"errors":"Не возможно приобрести товар"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			transaction := mock_services.NewMockTransaction(c)
			testCase.mockBehavior(transaction, 1)

			auth := mock_services.NewMockAuth(c)
			auth.EXPECT().ParseToken("some_valid_token").Return(1, nil)

			services := &services.Service{Transaction: transaction, Auth: auth}
			handler := NewHandler(services)

			r := chi.NewRouter()
			r.Get("/api/buy/{item}", handler.BuyItemHandler)

			w := httptest.NewRecorder()
			var req *http.Request

			if testCase.name == "INCORRECT VALUE" {
				req = httptest.NewRequest("GET", "/api/buy/incorect_value", bytes.NewBufferString(testCase.inputBody))
			} else {
				req = httptest.NewRequest("GET", "/api/buy/hoody", bytes.NewBufferString(testCase.inputBody))
			}

			req.AddCookie(&http.Cookie{Name: "auth_token", Value: "some_valid_token"})

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}