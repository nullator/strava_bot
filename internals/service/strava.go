package service

import (
	"strava_bot/internals/models"
	"strava_bot/internals/repository"
)

type StravaService struct {
	rep *repository.Repository
}

func NewStravaService(rep *repository.Repository) *StravaService {
	return &StravaService{rep: rep}
}

func (s *StravaService) Auth(userID, code string, chatID int64) (models.StravaUser, error) {
	return s.rep.Auth(userID, code, chatID)
}
