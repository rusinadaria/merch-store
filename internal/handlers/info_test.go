package handlers

import (
	mock_services "merch-store/internal/services/mocks"
	"merch-store/internal/services"
	"merch-store/models"
	"bytes"
	"net/http/httptest"
	"testing"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
)

func TestHandler_InfoHandler(t *testing.T) {
	type mockBehavior func(s *mock_services.MockInfo, userId int) models.InfoResponse

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
			mockBehavior: func(s *mock_services.MockInfo, userId int) models.InfoResponse {
				response := models.InfoResponse{
					Coins: 1000,
					Inventory: nil,
					CoinHistory: []models.CoinHistory{
						{
							Received: nil,
							Sent:     nil,
						},
					},
				}
				s.EXPECT().UserInfo(userId).Return(response, nil)
				return response
			},
			
			expectedStatusCode: 200,
			expectedResponseBody: `{"coins":1000,"inventory":null,"coinHistory":[{"received":null,"sent":null}]}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			info := mock_services.NewMockInfo(c)
			testCase.mockBehavior(info, 1)

			auth := mock_services.NewMockAuth(c)
			auth.EXPECT().ParseToken("some_valid_token").Return(1, nil)

			services := &services.Service{Info: info, Auth: auth}
			handler := NewHandler(services)

			r := chi.NewRouter()
			r.Get("/api/info", handler.InfoHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/info", bytes.NewBufferString(testCase.inputBody))

			req.AddCookie(&http.Cookie{Name: "auth_token", Value: "some_valid_token"})

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}