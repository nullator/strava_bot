package repository

import (
	"errors"
	"fmt"
	"log"
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

func (r *Repository) Auth(input *models.AuthHandler,
	resp *models.StravaUser) (int, *models.StravaUser, error) {
	tg_id := input.ID
	strava_id := fmt.Sprintf("%d", resp.Athlete.Id)
	expires := fmt.Sprintf("%d", resp.Expires_at)

	err := r.db.Save("id", strava_id, tg_id)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}

	err = r.db.Save("username", resp.Athlete.Username, tg_id)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}

	err = r.db.Save("acces_token", resp.Access_token, tg_id)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}

	err = r.db.Save("refresh_token", resp.Refresh_token, tg_id)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}

	err = r.db.Save("expies_at", expires, tg_id)
	if err != nil {
		return http.StatusInternalServerError, resp, err
	}

	return http.StatusOK, resp, nil
}

func (r *Repository) RefreshToken(id int64, input models.RespondRefreshToken) error {
	tg_id := fmt.Sprintf("%d", id)
	expires := fmt.Sprintf("%d", input.Expires_at)

	err := r.db.Save("acces_token", input.Access_token, tg_id)
	if err != nil {
		return err
	}
	log.Printf("acess token %s\n", input.Access_token)

	err = r.db.Save("expies_at", expires, tg_id)
	if err != nil {
		return err
	}

	err = r.db.Save("refresh_token", input.Refresh_token, tg_id)
	if err != nil {
		return err
	}
	log.Printf("new refresh token %s\n", input.Refresh_token)

	return nil
}

func (r *Repository) GetRefreshToken(id int64) (string, error) {
	var rt string
	tg_id := fmt.Sprintf("%d", id)
	rt, err := r.db.Get("refresh_token", tg_id)
	if err != nil {
		return "", err
	}
	if rt == "" {
		return "", errors.New("no refresh token in DB")
	}
	log.Printf("refresh token %s\n", rt)

	return rt, nil

}

func (r *Repository) GetAccesToken(id int64) (string, error) {
	var at string
	tg_id := fmt.Sprintf("%d", id)
	at, err := r.db.Get("acces_token", tg_id)
	if err != nil {
		return "", err
	}
	if at == "" {
		return "", errors.New("no acces token in DB")
	}

	return at, nil

}

func (r *Repository) GetExpies(id int64) (string, error) {
	var exp string
	tg_id := fmt.Sprintf("%d", id)
	exp, err := r.db.Get("expies_at", tg_id)
	if err != nil {
		return "", err
	}
	if exp == "" {
		return "", errors.New("no expies_at in DB")
	}

	return exp, nil

}

func (r *Repository) SaveActivityId(id int64, activity_id string) error {
	tg_id := fmt.Sprintf("%d", id)
	err := r.db.Save("last_activity_id", activity_id, tg_id)
	return err
}

func (r *Repository) GetActivityId(id int64) (string, error) {
	tg_id := fmt.Sprintf("%d", id)

	activity_id, err := r.db.Get("last_activity_id", tg_id)
	if err != nil {
		return "", err
	}
	if activity_id == "" {
		return "", errors.New("no activity_id in DB")
	}

	return activity_id, nil
}
