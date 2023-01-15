package repository

import (
	"strava_bot/internals/models"
	"strava_bot/pkg/base"
)

type Repository struct {
}

func NewRepository(base base.Base) *Repository {
	return &Repository{}
}

func (r *Repository) Auth(userID, code string, chatID int64) (models.StravaUser, error) {
	m := new(models.StravaUser)
	return *m, nil
}
