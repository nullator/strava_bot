package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
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

func (s *StravaService) UploadActivity(file string, id int64) error {
	data, err := os.Open(file)
	if err != nil {
		log.Println("Ошибка открытия файлв")
		return err
	}
	defer data.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file))
	if err != nil {
		log.Println("Ошибка создания part")
		return err
	}
	_, err = io.Copy(part, data)
	if err != nil {
		log.Println("Ошибка копирования файла в part")
		return err
	}

	l := len(file)
	file_format := string([]rune(file)[l-3 : l])
	err = writer.WriteField("data_type", file_format)
	if err != nil {
		log.Println("Ошибка добавления параметра data_type")
		return err
	}
	err = writer.Close()
	if err != nil {
		log.Println("Ошибка закрытия writer")
		return err
	}

	req, err := http.NewRequest("POST", strava_api_url+"uploads", body)
	if err != nil {
		log.Println("Ошибка создания req")
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = s.RefreshToken(id)
	if err != nil {
		log.Println("Ошибка обновления refresh token")
		return err
	}

	token, err := s.rep.GetAccesToken(id)
	if err != nil {
		log.Println("Ошибка получения токена")
		return err
	}
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка выполнения req")
		return err
	} else {
		if resp.StatusCode != 201 {
			return errors.New("incorrect upload activity file")
		}

		body, _ := io.ReadAll(resp.Body)
		var res models.RespondUploadActivity
		err = json.Unmarshal(body, &res)
		if err != nil && err != io.EOF {
			log.Println("Ошибка Обработки respond")
			return err
		}
		resp.Body.Close()
		log.Println(res.Id_str)
		log.Println(res)

		return s.rep.SaveActivityId(id, res.Id_str)
	}
}

func (s *StravaService) GetActivity(id int64) (string, error) {
	activity_id, err := s.rep.GetActivityId(id)
	if err != nil {
		log.Println("Ошибка получения activity_id")
		return "", err
	}

	// activity_id = "8435151828"
	log.Println(strava_api_url + "/activities/" + activity_id)
	req, err := http.NewRequest("GET", strava_api_url+"activities/"+activity_id, nil)
	if err != nil {
		log.Println("Ошибка создания req")
		return "", err
	}
	token, err := s.rep.GetAccesToken(id)
	if err != nil {
		log.Println("Ошибка получения токена")
		return "", err
	}
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка выполнения req")
		return "", err
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		log.Println(resp.StatusCode)
		log.Println(resp.Header)
		log.Println(body)
	}
	defer resp.Body.Close()

	// body, _ := io.ReadAll(resp.Body)
	// log.Println(resp.StatusCode)
	// log.Println(body)

	return "", nil

}
