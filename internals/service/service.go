package service

import "strava_bot/internals/models"

type StravaAuth interface {
	CreateAuthLink(iserID int64) (string, error)
	Auth(userID, code string, chatID int64) (models.StravaUser, error)
}

type Service struct {
	StravaAuth
}

func NewService() *Service {
	return &Service{}
}
