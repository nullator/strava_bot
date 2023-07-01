package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strava_bot/internals/models"
	"strava_bot/pkg/logger"
)

type TelegramService struct {
	logger *logger.Logger
}

func NewTelegramService(logger *logger.Logger) *TelegramService {
	return &TelegramService{logger: logger}
}

func (tg *TelegramService) GetFile(filename, fileid string) error {
	// Get the telegram file path from file_id
	querry := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s",
		os.Getenv("TG_TOKEN"), fileid)
	req, err := http.NewRequest("GET", querry, nil)
	if err != nil {
		tg.logger.Error("error create GET file request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tg.logger.Error("error do request (get telegram file): %v", err)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res models.TelegramFileIdResp
	err = json.Unmarshal(body, &res)
	if err != nil && err != io.EOF {
		tg.logger.Error("error parse telegram file ID: %v", err)
		return err
	}

	// download file from telegram server
	querry = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv("TG_TOKEN"),
		res.Result.File_path)
	req, err = http.NewRequest("GET", querry, nil)
	if err != nil {
		tg.logger.Error("error create request (download telegram file by ID): %v", err)
		return err
	}

	resp, err = client.Do(req)
	if err != nil {
		tg.logger.Error("error do request (download telegram file by ID): %v", err)
		return err
	}
	defer resp.Body.Close()

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		tg.logger.Error("error read file: %v", err)
		return err
	}
	err = os.MkdirAll("activity", os.ModePerm)
	if err != nil {
		tg.logger.Error("error create 'activity' directory: %v", err)
		return err
	}
	f, err := os.OpenFile("activity/"+filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		tg.logger.Error("error create new file in 'activity' directory: %v", err)
		return err
	}
	_, err = f.Write(file)
	if err != nil {
		tg.logger.Error("error write new file in 'activity' directory: %v", err)
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			tg.logger.Error("error save and close new activity file: %v", err)
		}
	}()

	tg.logger.Info("successful download and save new activity file: %s", filename)
	return nil

}
