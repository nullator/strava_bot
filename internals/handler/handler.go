package handler

import (
	"net/http"
	"strava_bot/internals/models"
	"strava_bot/internals/service"
	"strava_bot/internals/telegram"
	"strava_bot/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	bot      *telegram.Bot
	logger   *logger.Logger
}

func NewHandler(services *service.Service, bot *telegram.Bot, l *logger.Logger) *Handler {
	return &Handler{services, bot, l}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()
	router.GET("/auth", h.auth)
	return router
}

func (h *Handler) auth(c *gin.Context) {
	h.logger.Info("на auth поступил запрос авторизации")
	var input *models.AuthHandler
	var strava_user *models.StravaUser

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		h.logger.ErrorF("не удалось распарсить полученный запрос в JSON - %s", err.Error())
		return
	}

	code, strava_user, err := h.services.Auth(input)
	if err != nil {
		c.JSON(code, err.Error())
		return
	}

	c.Redirect(http.StatusMovedPermanently, "https://t.me/strava_ru_bot")

	tg_id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		h.logger.ErrorF("при выполнении авторизации не удалось распарсить ID в Telegram id - %s",
			err.Error())
	}

	h.bot.SuccsesAuth(tg_id, strava_user.Athlete.Username)

}
