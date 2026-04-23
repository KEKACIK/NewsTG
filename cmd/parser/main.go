package main

import (
	"context"
	"fmt"
	"newtg/internal/config"
	"newtg/internal/service/parser"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"
)

func main() {
	fmt.Println("HELLO PARSER")

	cfg := config.NewConfig()
	logger := logging.NewLogger(cfg.Debug)

	client, err := postgresql.NewClient(context.Background(), 5, cfg.GetPostgresDsn())
	if err != nil {
		logger.Fatal(err.Error())
	}

	rc := parser.NewRiaClient(client, logger, "РИА Новости", cfg.MaxNewsPerHourRia)
	rc.PoolNews(context.Background())
}
