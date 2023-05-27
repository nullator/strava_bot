package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"strava_bot/internals/handler"
	"strava_bot/internals/models"
	"strava_bot/internals/repository"
	"strava_bot/internals/service"
	"strava_bot/internals/telegram"
	boltdb "strava_bot/pkg/base/boltdb"
	"strava_bot/pkg/logger"

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

	l := log.Default()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	l.SetOutput(wrt)
	logger := logger.New("strava", l)

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

	rep := repository.NewRepository(base)
	service := service.NewService(rep)
	tg_bot := telegram.NewBot(bot, service)
	handlers := handler.NewHandler(service, tg_bot, logger)
	srv := new(models.Server)
	go func() {
		err := srv.Run(os.Getenv("SERVER_PORT"), handlers.InitRouters())
		if err != nil {
			log.Fatalf("error running server: %v", err)
		}
	}()

	tg_bot.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	err = srv.Stop(context.Background())
	if err != nil {
		log.Fatalf("error stopping server: %v", err)
	}

	err = db.Close()
	if err != nil {
		log.Fatalf("error closing db: %v", err)
	}

}
