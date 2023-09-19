package handler

import (
	"net/http"
	"net/http/httptest"
	"strava_bot/internals/models"
	"strava_bot/internals/service"
	"strava_bot/internals/telegram"
	"strava_bot/mocks"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {

	testTable := []struct {
		name               string
		url                string
		expCode            int
		excAuthFunc        func(*mocks.Strava, *models.AuthHandler)
		excSuccsesAuthFunc func(*mocks.BotInterface)
	}{
		{
			name:    "OK",
			url:     "/auth?state=123&code=authorization_code&scope=read_all",
			expCode: 301,
			excAuthFunc: func(mock *mocks.Strava, input *models.AuthHandler) {
				mock.
					On("Auth", input).
					Return(200, &models.StravaUser{Athlete: models.Athlete{Username: "name"}}, nil)
			},
			excSuccsesAuthFunc: func(mock *mocks.BotInterface) {
				mock.
					On("SuccsesAuth", int64(123), "name").
					Return(nil)
			},
		},
		{
			name:               "bad state",
			url:                "/auth?state=123b&code=authorization_code&scope=read_all",
			expCode:            400,
			excAuthFunc:        nil,
			excSuccsesAuthFunc: nil,
		},
		{
			name:               "empty state",
			url:                "/auth?state=&code=authorization_code&scope=read_all",
			expCode:            400,
			excAuthFunc:        nil,
			excSuccsesAuthFunc: nil,
		},
		{
			name:               "miss state",
			url:                "/auth?code=authorization_code&scope=read_all",
			expCode:            400,
			excAuthFunc:        nil,
			excSuccsesAuthFunc: nil,
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
			bot := mocks.NewBotInterface(t)
			tg_bot := &telegram.Bot{
				BotInterface: bot,
			}
			handler := NewHandler(service, tg_bot)

			stravaUser := models.StravaUser{}
			stravaUser.Athlete.Username = "name"
			r := chi.NewRouter()
			r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
				handler.authV2(w, r)
			})

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", test.url, nil)

			var input models.AuthHandler
			params := req.URL.Query()
			input.ID = params.Get("state")
			input.Code = params.Get("code")
			input.Scope = params.Get("scope")

			if test.excAuthFunc != nil {
				test.excAuthFunc(stravaService, &input)
			}
			if test.excSuccsesAuthFunc != nil {
				test.excSuccsesAuthFunc(bot)
			}

			// Make request
			r.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, w.Code, test.expCode)

		})

	}

}
