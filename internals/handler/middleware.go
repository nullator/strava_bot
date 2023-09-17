package handler

import (
	"errors"
	"net/http"
	"strava_bot/internals/models"
	"strconv"
)

func validateAuthModel(user *models.AuthHandler) (int, error) {
	if user == nil {
		return http.StatusBadRequest, errors.New("empty response")
	}

	if user.ID == "" {
		return http.StatusBadRequest, errors.New("empty user ID")
	}

	_, err := strconv.ParseInt(user.ID, 10, 64)
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid user ID")
	}

	if user.Code == "" {
		return http.StatusBadRequest, errors.New("empty authorization code")
	}

	if user.Scope == "" {
		return http.StatusBadRequest, errors.New("empty scope")
	}

	return 200, nil

}
