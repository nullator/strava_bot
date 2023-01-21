package models

type TelegramFileIdResp struct {
	Ok     bool               `json:"ok"`
	Result TelegramFileResult `json:"result"`
}

type TelegramFileResult struct {
	File_id        string `json:"file_id"`
	File_unique_id string `json:"file_unique_id"`
	File_size      int    `json:"file_size"`
	File_path      string `json:"file_path"`
}
