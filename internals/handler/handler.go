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
	var user *models.AuthHandler
	var strava_user *models.StravaUser

	if c.ShouldBind(&user) == nil {
		log.Println(user.ID)
		log.Println(user.Code)
		log.Println(user.Scope)
	}

	code, strava_user, err := h.services.Auth(user)
	if err != nil {
		c.JSON(code, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"login at": strava_user.Athlete.Username,
		"status":   "OK",
	})
}
