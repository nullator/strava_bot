package main

import (
	"io"
	"log"
	"os"

	"strava_bot/internals/telegram"
	boltdb "strava_bot/pkg/base/boltDb"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := os.MkdirAll("log", os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.OpenFile("log/all.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	bot.Debug = false

	db, err := bolt.Open(os.Getenv("DB_FILE"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	base := boltdb.NewBase(db)
	tg_bot := telegram.NewBot(bot, base)
	tg_bot.Start()

}
