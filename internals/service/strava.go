package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strava_bot/internals/models"
	"strava_bot/internals/repository"
	"strconv"
	"time"
)

const (
	strava_auth_url = "https://www.strava.com/api/v3/oauth/token"
	strava_api_url  = "https://www.strava.com/api/v3/"
)

type StravaService struct {
	rep *repository.Repository
}

func NewStravaService(rep *repository.Repository) *StravaService {
	return &StravaService{rep: rep}
}

func (s *StravaService) Auth(input *models.AuthHandler) (int, *models.StravaUser, error) {
	const path = "internal.service.strava.Auth"

	// make request to STRAVA
	request := map[string]string{
		"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
		"client_secret": os.Getenv("STRAVA_SECRET"),
		"code":          input.Code,
		"grant_type":    "authorization_code",
	}

	json_request, err := json.Marshal(request)
	if err != nil {
		slog.Error("error generate json to make strava auth",
			slog.String("error", err.Error()))
		return http.StatusInternalServerError, nil, fmt.Errorf("%s: %w", path, err)
	}

	req, err := http.NewRequest("POST", strava_auth_url, bytes.NewBuffer(json_request))
	if err != nil {
		slog.Error("error POST request strava auth",
			slog.String("error", err.Error()))
		return http.StatusInternalServerError, nil, fmt.Errorf("%s: %w", path, err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error make request strava auth",
			slog.String("error", err.Error()))
		return http.StatusInternalServerError, nil, fmt.Errorf("%s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		slog.Error("autorization error (resp.StatusCode != 200)",
			slog.Int("get code", resp.StatusCode))
		return http.StatusInternalServerError, nil, fmt.Errorf("%s: auth error", path)
	}

	// get Athlete model
	body, _ := io.ReadAll(resp.Body)
	var res models.StravaUser
	err = json.Unmarshal(body, &res)
	if err != nil && err != io.EOF {
		slog.Error("error unmarshal athlete model",
			slog.String("error", err.Error()))
		return http.StatusInternalServerError, nil, fmt.Errorf("%s: %w", path, err)
	}

	// write data to repository
	slog.Info("successful receipt of authorization data,"+
		" ready to writing data to the database",
		slog.Int64("strava_id", res.Athlete.Id),
		slog.String("strava_name", res.Athlete.Username),
		slog.String("scope", input.Scope),
		slog.String("strava_city", res.Athlete.City),
	)
	code, user, err := s.rep.Auth(input, &res)
	if err != nil {
		return code, nil, fmt.Errorf("%s: %w", path, err)
	}
	return code, user, nil
}

func (s *StravaService) UploadActivity(file string, id int64) error {
	const path = "internal.service.strava.UploadActivity"

	data, err := os.Open(file)
	if err != nil {
		slog.Error("error opening file",
			slog.String("file", file),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("%s: %w", path, err)
	}
	defer data.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file))
	if err != nil {
		slog.Error("error create part", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}
	_, err = io.Copy(part, data)
	if err != nil {
		slog.Error("error copy file to part", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}

	l := len(file)
	file_format := string([]rune(file)[l-3 : l])
	err = writer.WriteField("data_type", file_format)
	if err != nil {
		slog.Error("error add 'data_type' field",
			slog.String("data_type", file_format),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("%s: %w", path, err)
	}
	err = writer.Close()
	if err != nil {
		slog.Error("error close writer", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}

	req, err := http.NewRequest("POST", strava_api_url+"uploads", body)
	if err != nil {
		slog.Error("error create upload request", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = s.RefreshToken(id)
	if err != nil {
		slog.Error("error refresh token", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}

	token, err := s.rep.GetAccesToken(id)
	if err != nil {
		slog.Error("error get acces token from DB", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error do request (upload file to strava)",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	} else {
		if resp.StatusCode != 201 {
			slog.Error("error do request (resp.Code != 201)",
				slog.Int("get code", resp.StatusCode))
			return fmt.Errorf("%s: incorrect upload activity file", path)
		}

		body, _ := io.ReadAll(resp.Body)
		var res models.RespondUploadActivity
		err = json.Unmarshal(body, &res)
		if err != nil && err != io.EOF {
			slog.Error("error response processing aftet upload file to strava",
				slog.String("error", err.Error()))
			return fmt.Errorf("%s: %w", path, err)
		}
		resp.Body.Close()
		slog.Info("successful upload activity",
			slog.Int64("id", res.Id),
			slog.String("id_str", res.Id_str),
			slog.Int64("activity_id", res.Activity_id),
			slog.Any("activity_id_ext", res.External_id),
			slog.String("status", res.Status),
		)

		err = s.rep.SaveActivityId(id, res.Id_str)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}

		return nil
	}
}

func (s *StravaService) RefreshToken(id int64) error {
	const path = "internal.service.strava.RefreshToken"

	et, err := s.rep.GetExpies(id)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	delta, err := checkRefreshToken(et)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	if delta <= 3600 {
		return s.getNewToken(id)
	} else {
		slog.Info("no token refresh needed")
		return nil
	}

}

func checkRefreshToken(exp_time string) (int64, error) {
	const path = "internal.service.strava.checkRefreshToken"

	et_int, err := strconv.ParseInt(exp_time, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", path, err)
	}

	now_time := time.Now().Unix()
	delta := et_int - now_time

	return delta, nil

}

func (s *StravaService) getNewToken(id int64) error {
	const path = "internal.service.strava.getNewToken"

	old_refresh_token, err := s.rep.GetRefreshToken(id)
	if err != nil {
		slog.Error("error get refresh token from DB", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}

	request := map[string]string{
		"client_id":     os.Getenv("STRAVA_CLIENT_ID"),
		"client_secret": os.Getenv("STRAVA_SECRET"),
		"grant_type":    "refresh_token",
		"refresh_token": old_refresh_token,
	}

	json_request, err := json.Marshal(request)
	if err != nil {
		slog.Error("error make json_request", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}

	req, err := http.NewRequest("POST", strava_auth_url, bytes.NewBuffer(json_request))
	if err != nil {
		slog.Error("error creating request to strava (refresh token)",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error do request to strava (refresh token)",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		slog.Error("error get new refresh token (status code != 200)",
			slog.Int("get code", resp.StatusCode))
		return fmt.Errorf("%s: autorization error", path)
	}

	body, _ := io.ReadAll(resp.Body)
	var res models.RespondRefreshToken
	err = json.Unmarshal(body, &res)
	if err != nil && err != io.EOF {
		slog.Error("error parse new refresh token", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}

	err = s.rep.RefreshToken(id, res)
	if err != nil {
		slog.Error("error write new refresh token to DB",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	} else {
		slog.Info("successful get new refresh token")
		return nil
	}
}
