package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
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
	// load environment
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env: %v", err)
	}

	// setup logging
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
	l.SetOutput(wrt)

	log, err := setupLogger(os.Getenv("ENV"))
	if err != nil {
		l.Fatal(err)
	}

	log.Info("start", slog.String("env", os.Getenv("ENV")))
	log.Debug("debug level is enabled")

	// creare TG bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		l.Fatal(err.Error())
	}
	bot.Debug = false

	// open DB
	db, err := bolt.Open(os.Getenv("DB_FILE"), 0600, nil)
	if err != nil {
		l.Fatal(err.Error())
	}
	defer func() {
		err := db.Close()
		if err != nil {
			l.Fatal(err.Error())
		}
	}()
	base := boltdb.NewBase(db)

	// init
	rep := repository.NewRepository(base)
	service := service.NewService(rep, log)
	tg_bot := telegram.NewBot(bot, service)
	handlers := handler.NewHandler(service, tg_bot)
	srv := new(models.Server)
	go func() {
		err := srv.Run(os.Getenv("SERVER_PORT"), handlers.InitRouters())
		if err != nil {
			l.Fatalf("error running server: %v", err)
		}
	}()

	// run bot
	tg_bot.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	err = srv.Stop(context.Background())
	if err != nil {
		l.Fatalf("error stopping server: %v", err)
	}

	err = db.Close()
	if err != nil {
		l.Fatalf("error closing db: %v", err)
	}

}

func setupLogger(env string) (*slog.Logger, error) {
	var log *slog.Logger

	switch env {
	case "local":
		h := logger.NewCustomSlogHandler(slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: false,
			}))
		log = slog.New(h)
	case "prod":
		h := logger.NewCustomSlogHandler(slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: false,
			}))
		log = slog.New(h)
	default:
		return nil, fmt.Errorf("incorrect error level: %s", env)
	}

	return log, nil
}
