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

func (tg *TelegramService) GetFile(filename, fileid string) error {
	// Get the telegram file path from file_id
	querry := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s",
		os.Getenv("TG_TOKEN"), fileid)
	req, err := http.NewRequest("GET", querry, nil)
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
	var res models.TelegramFileIdResp
	err = json.Unmarshal(body, &res)
	if err != nil && err != io.EOF {
		return err
	}

	// download file from telegram server
	querry = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv("TG_TOKEN"),
		res.Result.File_path)
	req, err = http.NewRequest("GET", querry, nil)
	if err != nil {
		return err
	}

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.MkdirAll("activity", os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.OpenFile("activity/"+filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write(file)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Println("error save and close activity file")
		}
	}()

	return nil

}
