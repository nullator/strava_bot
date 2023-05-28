package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strava_bot/internals/models"
	"strava_bot/internals/repository"
	"strava_bot/pkg/logger"
	"strconv"
	"time"
)

const (
	strava_auth_url = "https://www.strava.com/api/v3/oauth/token"
	strava_api_url  = "https://www.strava.com/api/v3/"
)

type StravaService struct {
	rep    *repository.Repository
	logger *logger.Logger
}

func NewStravaService(rep *repository.Repository, logger *logger.Logger) *StravaService {
	return &StravaService{rep: rep, logger: logger}
}

func (s *StravaService) Auth(input *models.AuthHandler) (int, *models.StravaUser, error) {
	// validate auth
	code, err := validateAuthModel(input)
	if err != nil {
		s.logger.Errorf("error validate auth model: %v\n", err)
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
		s.logger.Errorf("error generate json to make strava auth: %v\n", err)
		return http.StatusInternalServerError, nil, err
	}

	req, err := http.NewRequest("POST", strava_auth_url, bytes.NewBuffer(json_request))
	if err != nil {
		s.logger.Errorf("error POST request strava auth: %v\n", err)
		return http.StatusInternalServerError, nil, err
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Errorf("error make request strava auth: %v\n", err)
		return http.StatusInternalServerError, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		s.logger.Errorf("autorization error: %v\n", err)
		return http.StatusInternalServerError, nil, errors.New("autorization error")
	}

	// get Athlete model
	body, _ := io.ReadAll(resp.Body)
	var res models.StravaUser
	err = json.Unmarshal(body, &res)
	if err != nil && err != io.EOF {
		s.logger.Errorf("error unmarshal athlete model: %v\n", err)
		return http.StatusInternalServerError, nil, err
	}

	// write data to repository
	s.logger.Infof("successful receipt of authorization data, writing data to the database "+
		"(strava_id=%d, strava_name=%s, scope=%s\n)",
		res.Athlete.Id, res.Athlete.Username, input.Scope)
	return s.rep.Auth(input, &res)
}

func (s *StravaService) UploadActivity(file string, id int64) error {
	data, err := os.Open(file)
	if err != nil {
		s.logger.Errorf("error opening file %s: %v\n", file, err)
		return err
	}
	defer data.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file))
	if err != nil {
		s.logger.Errorf("error create part: %v\n", err)
		return err
	}
	_, err = io.Copy(part, data)
	if err != nil {
		s.logger.Errorf("error copy file to part: %v\n", err)
		return err
	}

	l := len(file)
	file_format := string([]rune(file)[l-3 : l])
	err = writer.WriteField("data_type", file_format)
	if err != nil {
		s.logger.Errorf("error add 'data_type' field: %v\n", err)
		return err
	}
	err = writer.Close()
	if err != nil {
		s.logger.Errorf("error close writer: %v\n", err)
		return err
	}

	req, err := http.NewRequest("POST", strava_api_url+"uploads", body)
	if err != nil {
		s.logger.Errorf("error create upload request: %v\n", err)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = s.RefreshToken(id)
	if err != nil {
		s.logger.Errorf("error refresh token: %v\n", err)
		return err
	}

	token, err := s.rep.GetAccesToken(id)
	if err != nil {
		s.logger.Errorf("error get acces token from DB: %v\n", err)
		return err
	}
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Errorf("error do request (upload file to strava): %v\n", err)
		return err
	} else {
		if resp.StatusCode != 201 {
			return errors.New("incorrect upload activity file")
		}

		body, _ := io.ReadAll(resp.Body)
		var res models.RespondUploadActivity
		err = json.Unmarshal(body, &res)
		if err != nil && err != io.EOF {
			s.logger.Errorf("error response processing aftet upload file to strava: %v\n", err)
			return err
		}
		resp.Body.Close()
		s.logger.Infof("successful upload activity (id: %d, activity_id: %d, status: %s)\n",
			res.Id, res.Activity_id, res.Status)

		return s.rep.SaveActivityId(id, res.Id_str)
	}
}

func (s *StravaService) RefreshToken(id int64) error {

	et, err := s.rep.GetExpies(id)
	if err != nil {
		return err
	}

	delta, err := checkRefreshToken(et)
	if err != nil {
		return err
	}

	if delta <= 3600 {
		return s.getNewToken(id)
	} else {
		s.logger.Info("no token refresh needed")
		return nil
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
		s.logger.Errorf("error get refresh token from DB: %v\n", err)
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
		s.logger.Errorf("error creating request to strava (refresh token): %v\n", err)
		return err
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Errorf("error do request to strava (refresh token): %v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		s.logger.Errorf("error get new refresh token, get http.status %d\n", resp.StatusCode)
		return errors.New("autorization error")
	}

	body, _ := io.ReadAll(resp.Body)
	var res models.RespondRefreshToken
	err = json.Unmarshal(body, &res)
	if err != nil && err != io.EOF {
		s.logger.Errorf("error parse new refresh token: %v\n", err)
		return err
	}

	err = s.rep.RefreshToken(id, res)
	if err != nil {
		s.logger.Errorf("error write new refresh token to DB: %v\n", err)
		return err
	} else {
		s.logger.Info("successful get new refresh token")
		return nil
	}
}
