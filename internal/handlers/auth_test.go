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
)

func TestHandler_AddUserHandler(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAuth, user models.AuthRequest)

	testTable := []struct {
		name string
		inputBody string
		inputUser models.AuthRequest
		mockBehavior mockBehavior
		expectedStatusCode int
		expectedResponseBody string
	} {
		{
			name: "OK",
			inputBody: `{"username":"test", "password":"12345"}`,
			inputUser: models.AuthRequest {
				Username: "test",
				Password: "12345",
			},
			mockBehavior: func(s *mock_services.MockAuth, user models.AuthRequest) {
				s.EXPECT().FindUser(user.Username).Return(1, nil)
				s.EXPECT().SignIn(user.Username, user.Password).Return(1, nil)
				s.EXPECT().GenerateToken(1).Return("token", nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: ``,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_services.NewMockAuth(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &services.Service{Auth: auth}
			handler := NewHandler(services)

			r := chi.NewRouter()
			r.Post("/api/auth", handler.AddUserHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/auth", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}