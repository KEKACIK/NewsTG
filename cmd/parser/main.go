package main

import (
	"context"
	"newtg/config"
	"newtg/internal/service/parser"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"

	"github.com/robfig/cron/v3"
)

func main() {
	cfg := config.NewConfig()
	logger := logging.NewLogger(cfg.Debug)

	client, err := postgresql.NewClient(context.Background(), logger, 5, cfg.GetPostgresDsn())
	if err != nil {
		logger.Fatal(err.Error())
	}

	riaClient := parser.NewRiaClient(client, logger, "РИА Новости", cfg.MaxNewsPerHourRia)
	parsers := []parser.Parser{riaClient}

	c := cron.New()
	c.AddFunc("*/15 * * * *", func() {
		logger.Info("Start parsing")
		for _, p := range parsers {
			go p.PoolNews(context.Background())
		}
	})
	c.Start()

	logger.Info("Parser start pooling..")
	select {}
}
