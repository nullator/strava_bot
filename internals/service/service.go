package service

import (
	"log/slog"
	"strava_bot/internals/models"
	"strava_bot/internals/repository"
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
	Logger *slog.Logger
}

func NewService(rep *repository.Repository, log *slog.Logger) *Service {
	return &Service{
		Strava:   NewStravaService(rep, log),
		Telegram: NewTelegramService(log),
		Logger:   log,
	}
}
