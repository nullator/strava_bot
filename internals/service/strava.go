package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strava_bot/internals/models"
	"strava_bot/internals/repository"
	"strconv"
	"time"
)

const (
	strava_auth_url = "https://www.strava.com/api/v3/oauth/token"
)

type StravaService struct {
	rep *repository.Repository
}

func NewStravaService(rep *repository.Repository) *StravaService {
	return &StravaService{rep: rep}
}

func (s *StravaService) Auth(input *models.AuthHandler) (int, *models.StravaUser, error) {
	// validate auth
	code, err := validateAuthModel(input)
	if err != nil {
		return code, nil, err
	}

	// make request to STRAVA
	request := map[string]string{
		"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
		"client_secret": os.Getenv("STRAVA_SECRET"),
		"code":          input.Code,
		"grant_type":    "authorization_code",
	}

	json_request, err := json.Marshal(request)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	req, err := http.NewRequest("POST", strava_auth_url, bytes.NewBuffer(json_request))
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
	return s.rep.Auth(input, &res)
}

func (s *StravaService) RefreshToken(id int64) (string, error) {

	et, err := s.rep.GetExpies(id)
	if err != nil {
		return "", err
	}
	delta, err := checkRefreshToken(et)
	if err != nil {
		return "", err
	}

	if delta <= 3600 {
		return "Token was refreshed", s.getNewToken(id)
	} else {
		return "Token was not refreshed", nil
	}

}

func checkRefreshToken(exp_time string) (int64, error) {
	et_int, err := strconv.ParseInt(exp_time, 10, 32)
	if err != nil {
		return 0, err
	}

	now_time := time.Now().Unix()
	delta := et_int - now_time

	return delta, nil

}

func (s *StravaService) getNewToken(id int64) error {
	old_refresh_token, err := s.rep.GetRefreshToken(id)
	if err != nil {
		return err
	}

	request := map[string]string{
		"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
		"client_secret": os.Getenv("STRAVA_SECRET"),
		"grant_type":    "refresh_token",
		"refresh_token": old_refresh_token,
	}

	json_request, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", strava_auth_url, bytes.NewBuffer(json_request))
	if err != nil {
		return err
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res models.RespondRefreshToken
	err = json.Unmarshal(body, &res)
	if err != nil && err != io.EOF {
		return err
	}

	err = s.rep.RefreshToken(id, res)
	return err
}
