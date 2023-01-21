package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strava_bot/internals/models"
)

type TelegramService struct {
}

func NewTelegramService() *TelegramService {
	return &TelegramService{}
}

func (tg *TelegramService) GetFile(file_id string) (string, error) {
	querry := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", os.Getenv("TG_TOKEN"), file_id)
	req, err := http.NewRequest("GET", querry, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res models.TelegramFileIdResp
	err = json.Unmarshal(body, &res)
	if err != nil && err != io.EOF {
		return "", err
	}

	querry = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv("TG_TOKEN"), res.Result.File_path)
	req, err = http.NewRequest("GET", querry, nil)
	if err != nil {
		return "", err
	}

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile("test.pdf", os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	_, err = f.Write(file)
	if err != nil {
		return "", err
	}

	body, _ = io.ReadAll(resp.Body)
	log.Println(string(body))

	return "OK", nil

}
