package main

import (
	"context"
	"newtg/config"
	"newtg/internal/service/poster"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"
)

func main() {
	cfg := config.NewConfig()
	logger := logging.NewLogger(cfg.Debug)

	client, err := postgresql.NewClient(context.Background(), logger, 5, cfg.GetPostgresDsn())
	if err != nil {
		logger.Fatal(err.Error())
	}

	telegramPoster, err := poster.NewTelegramPoster(
		client,
		logger,
		// TELEGRAM
		cfg.TelegramBotToken,
		cfg.TelegramChatID,
		cfg.TelegramMaxMsgLength,
		// RIA NEWS
		"РИА Новости",
		cfg.MaxNewsPerHourRia,
	)
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = telegramPoster.StartPool(context.Background())
	if err != nil {
		logger.Fatal(err.Error())
	}
}
