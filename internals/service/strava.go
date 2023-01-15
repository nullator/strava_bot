package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strava_bot/internals/models"
	"strava_bot/internals/repository"
)

const (
	strava_url = "https://www.strava.com/api/v3/oauth/token"
)

type StravaService struct {
	rep *repository.Repository
}

func NewStravaService(rep *repository.Repository) *StravaService {
	return &StravaService{rep: rep}
}

func (s *StravaService) Auth(user *models.AuthHandler) (int, *models.StravaUser, error) {
	// validate auth
	code, err := validateAuthModel(user)
	if err != nil {
		return code, nil, err
	}

	// make request to STRAVA
	request := map[string]string{
		"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
		"client_secret": os.Getenv("STRAVA_SECRET"),
		"code":          user.Code,
		"grant_type":    "authorization_code",
	}

	json_request, err := json.Marshal(request)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	req, err := http.NewRequest("POST", strava_url, bytes.NewBuffer(json_request))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer resp.Body.Close()

	// get Athlete model
	body, _ := io.ReadAll(resp.Body)
	var res models.StravaUser
	err = json.Unmarshal(body, &res)
	if err != nil && err != io.EOF {
		return http.StatusInternalServerError, nil, err
	}

	// write data to repository
	return s.rep.Auth(user, &res)
}
