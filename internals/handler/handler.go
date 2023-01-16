package handler

import (
	"log"
	"net/http"
	"strava_bot/internals/models"
	"strava_bot/internals/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
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
	log.Printf(strava_user.Athlete.Username)

	c.Redirect(http.StatusMovedPermanently, "https://t.me/strava_ru_bot")

}
