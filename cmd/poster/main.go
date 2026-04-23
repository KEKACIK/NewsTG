package main

import (
	"context"
	"fmt"
	"newtg/internal/config"
	"newtg/internal/service/poster"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"
)

func main() {
	fmt.Println("HELLO POSTER")

	ctx := context.TODO()

	cfg := config.NewConfig()
	logger := logging.NewLogger(cfg.Debug)
	client, err := postgresql.NewClient(ctx, 5, cfg.GetPostgresDsn())

	telegramPoster, err := poster.NewTelegramPoster(
		client,
		logger,

		cfg.TelegramBotToken,
		cfg.TelegramChatID,

		cfg.MaxNewsLength,
	)
	if err != nil {
		logger.Fatal(err.Error())
	}

	telegramPoster.CheckNews(ctx)
	// for {
	// 	now := time.Now()
	// 	// next := now.Truncate(time.Hour).Add(time.Hour)
	// 	next := now.Truncate(time.Hour).Add(10 * time.Second)

	// 	timer := time.NewTimer(time.Until(next))
	// 	fmt.Printf("Следующий запуск в: %v\n", next.Format("15:04:05"))

	// 	<-timer.C

	// 	go telegramPoster.CheckNews()
	// }
}
