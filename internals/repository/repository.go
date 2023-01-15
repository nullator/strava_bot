package repository

import (
	"fmt"
	"net/http"
	"strava_bot/internals/models"
	"strava_bot/pkg/base"
)

type Repository struct {
	db base.Base
}

func NewRepository(db base.Base) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Auth(user *models.AuthHandler,
	resp *models.StravaUser) (int, *models.StravaUser, error) {
	strava_id := fmt.Sprintf("%d", resp.Athlete.Id)
	expires := fmt.Sprintf("%d", resp.Expires_at)

	err := r.db.Save("id", strava_id, user.ID)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}

	err = r.db.Save("username", resp.Athlete.Username, strava_id)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}

	err = r.db.Save("acces_token", resp.Access_token, strava_id)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}

	err = r.db.Save("refresh_token", resp.Refresh_token, strava_id)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}

	err = r.db.Save("expies_at", expires, strava_id)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}

	return http.StatusOK, resp, nil
}
