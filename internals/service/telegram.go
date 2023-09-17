package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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
	const path = "internal.service.telegram.GetFile"

	// Get the telegram file path from file_id
	querry := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s",
		os.Getenv("TG_TOKEN"), fileid)
	req, err := http.NewRequest("GET", querry, nil)
	if err != nil {
		slog.Error("error create GET file request", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error do request (get telegram file)",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res models.TelegramFileIdResp
	err = json.Unmarshal(body, &res)
	if err != nil && err != io.EOF {
		slog.Error("error parse telegram file ID", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}

	// download file from telegram server
	querry = fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", os.Getenv("TG_TOKEN"),
		res.Result.File_path)
	req, err = http.NewRequest("GET", querry, nil)
	if err != nil {
		slog.Error("error create request (download telegram file by ID)",
			slog.String("file_path", res.Result.File_path),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("%s: %w", path, err)
	}

	resp, err = client.Do(req)
	if err != nil {
		slog.Error("error do request (download telegram file by ID)",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}
	defer resp.Body.Close()

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("error read file", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}
	err = os.MkdirAll("activity", os.ModePerm)
	if err != nil {
		slog.Error("error create 'activity' directory",
			slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", path, err)
	}
	f, err := os.OpenFile("activity/"+filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("error create new file in 'activity' directory",
			slog.String("filename", filename),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("%s: %w", path, err)
	}
	_, err = f.Write(file)
	if err != nil {
		slog.Error("error write new file in 'activity' directory",
			slog.String("filename", filename),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("%s: %w", path, err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			slog.Error("error save and close new activity file:",
				slog.String("filename", filename),
				slog.String("error", err.Error()),
			)
		}
	}()

	slog.Info("successful download and save new activity file",
		slog.String("filename", filename))
	return nil
}
