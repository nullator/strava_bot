package service

import (
	"strava_bot/internals/models"
	"strava_bot/internals/repository"
)

type Strava interface {
	Auth(input *models.AuthHandler) (int, *models.StravaUser, error)
	RefreshToken(id int64) error
	UploadActivity(file string, id int64) error
}

type Telegram interface {
	GetFile(filename, fileid string) error
	GetStatus(id int64) (string, error)
	SaveStatus(id int64, status string) error
}

type Service struct {
	Strava
	Telegram
}

func NewService(rep *repository.Repository) *Service {
	return &Service{
		Strava:   NewStravaService(rep),
		Telegram: NewTelegramService(rep),
	}
}
