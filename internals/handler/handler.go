package handler

import (
	"log"
	"net/http"
	"strava_bot/internals/models"
	"strava_bot/internals/service"
	"strava_bot/internals/telegram"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *Handler) auth(c *gin.Context) {
	var input *models.AuthHandler
	var strava_user *models.StravaUser

	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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
		log.Println(err.Error())
	}

	h.bot.SuccsesAuth(tg_id, strava_user.Athlete.Username)

}
