package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"tg_botfather/config"
	"tg_botfather/internal/controller"
	rep "tg_botfather/internal/repository"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	if err := config.InitENV(dir); err != nil {
		return err
	}

	cfg := config.GetConfig()

	repository, err := rep.InitConnect(*cfg)
	if err != nil {
		return err
	}

	if err := repository.AutoMigrate(); err != nil {
		return err
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	handlers := controller.NewHandler(repository, bot)

	router := controller.InitRouter(*handlers)

	srv := http.Server{
		Addr:              ":" + cfg.Port,
		ReadHeaderTimeout: time.Millisecond * 500,
		ReadTimeout:       time.Millisecond * 500,
		Handler:           router,
	}

	idleConnClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Print(err)
		}
		close(idleConnClosed)
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Print(err)
		return err
	}

	<-idleConnClosed

	return nil
}
