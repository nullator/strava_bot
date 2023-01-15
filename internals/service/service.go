package service

import (
	"strava_bot/internals/models"
	"strava_bot/internals/repository"
)

type Strava interface {
	//CreateAuthLink(iserID int64) (string, error)
	Auth(userID, code string, chatID int64) (models.StravaUser, error)
}

type Service struct {
	Strava
}

func NewService(rep *repository.Repository) *Service {
	return &Service{
		Strava: NewStravaService(rep),
	}
}
