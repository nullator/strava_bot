package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"strava_bot/internals/models"
	"strava_bot/internals/service"
	"strava_bot/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {

	testTable := []struct {
		name         string
		inputModels  *models.AuthHandler
		inputBody    string
		expectedCode int
		expectedBody string
	}{
		{
			name: "OK",
			inputModels: &models.AuthHandler{
				ID:    "123",
				Code:  "authorization_code",
				Scope: "read_all",
			},
			inputBody:    `{"state": "123", "code": "autorization_code", "scope": "read_all"}`,
			expectedCode: 200,
			expectedBody: "OK",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			// Init
			stravaService := mocks.NewStrava(t)
			telegramService := mocks.NewTelegram(t)
			service := &service.Service{
				Strava:   stravaService,
				Telegram: telegramService,
			}
			handler := NewHandler(service, nil)

			stravaService.On("Auth", test.inputModels).Return(400, nil, errors.New("test"))

			// Test server
			r := gin.New()
			r.POST("/auth", handler.auth)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth", bytes.NewBufferString(test.inputBody))

			// Make request
			r.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, w.Code, test.expectedCode)
			assert.Equal(t, w.Body.String(), test.expectedBody)

		})

	}

}
