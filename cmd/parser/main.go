package main

import (
	"context"
	"newtg/config"
	"newtg/internal/service/parser"
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

	rc := parser.NewRiaClient(client, logger, "РИА Новости", cfg.MaxNewsPerHourRia)
	rc.PoolNews(context.Background())
}
