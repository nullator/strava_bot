package handler

import (
	"log/slog"
	"net/http"
	"strava_bot/internals/models"
	"strava_bot/internals/service"
	"strava_bot/internals/telegram"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi"
)

type Handler struct {
	services *service.Service
	bot      *telegram.Bot
}

func NewHandler(services *service.Service, bot *telegram.Bot) *Handler {
	return &Handler{services, bot}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()
	router.GET("/auth", h.auth)
	return router
}

func (h *Handler) InitRoutersV2() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		h.authV2(w, r)
	})

	return router
}

func (h *Handler) auth(c *gin.Context) {
	slog.Info("на auth поступил запрос авторизации")
	var input *models.AuthHandler
	var strava_user *models.StravaUser

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		slog.Error("не удалось распарсить полученный запрос в JSON",
			slog.String("input", c.FullPath()),
			slog.String("error", err.Error()),
		)
		return
	}
	slog.Info("запрос распарсен")

	code, err := validateAuthModel(input)
	if err != nil {
		c.JSON(code, err.Error())
		slog.Error("input data validation error",
			slog.String("input.Code", input.Code),
			slog.String("input.ID", input.ID),
			slog.String("input.Scope", input.Scope),
			slog.String("error", err.Error()),
		)
		return
	}

	code, strava_user, err = h.services.Auth(input)
	if err != nil {
		c.JSON(code, err.Error())
		slog.Error("auth error", slog.String("error", err.Error()))
		return
	}

	c.Redirect(http.StatusMovedPermanently, "https://t.me/strava_ru_bot")

	tg_id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		slog.Error("при выполнении авторизации не удалось распарсить ID в TelegramId",
			slog.String("error", err.Error()))
	}

	h.bot.SuccsesAuth(tg_id, strava_user.Athlete.Username)

}

func (h *Handler) authV2(w http.ResponseWriter, resp *http.Request) {
	slog.Info("на auth поступил запрос авторизации")
	var input models.AuthHandler
	var strava_user *models.StravaUser

	params := resp.URL.Query()
	input.ID = params.Get("state")
	input.Code = params.Get("code")
	input.Scope = params.Get("scope")
	slog.Debug("запрос распарсен",
		slog.String("ID", input.ID),
		slog.String("code", input.Code),
		slog.String("scope", input.Scope),
	)

	code, err := validateAuthModel(&input)
	if err != nil {
		slog.Error("input data validation error",
			slog.Int("code", code),
			slog.String("input.Code", input.Code),
			slog.String("input.ID", input.ID),
			slog.String("input.Scope", input.Scope),
			slog.String("error", err.Error()),
		)
		w.WriteHeader(code)
		return
	}

	code, strava_user, err = h.services.Auth(&input)
	if err != nil {
		slog.Error("auth error",
			slog.Int("code", code),
			slog.String("error", err.Error()),
		)
		w.WriteHeader(code)
		return
	}

	http.Redirect(w, resp, "https://t.me/strava_ru_bot", http.StatusMovedPermanently)

	tg_id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		slog.Error("при выполнении авторизации не удалось распарсить ID в TelegramId",
			slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.bot.SuccsesAuth(tg_id, strava_user.Athlete.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
