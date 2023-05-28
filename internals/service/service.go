package service

import (
	"strava_bot/internals/models"
	"strava_bot/internals/repository"
	"strava_bot/pkg/logger"
)

type Strava interface {
	Auth(input *models.AuthHandler) (int, *models.StravaUser, error)
	RefreshToken(id int64) error
	getNewToken(id int64) error
	UploadActivity(file string, id int64) error
}

type Telegram interface {
	GetFile(filename, fileid string) error
}

type Service struct {
	Strava
	Telegram
	Logger *logger.Logger
}

func NewService(rep *repository.Repository, logger *logger.Logger) *Service {
	return &Service{
		Strava:   NewStravaService(rep, logger),
		Telegram: NewTelegramService(logger),
		Logger:   logger,
	}
}
