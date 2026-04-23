package postgresql

import (
	"context"
	"newtg/pkg/logging"
	"newtg/pkg/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func Migration(ctx context.Context, logger *logging.Logger, pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatal(err.Error())
	}

	if err := goose.Up(db, "migrations"); err != nil {
		logger.Fatal(err.Error())
	}

	return nil
}

func NewClient(
	ctx context.Context,
	logger *logging.Logger,
	maxAttempts int,
	dsn string,
) (pool *pgxpool.Pool, err error) {
	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		logger.Debug("database trying to connect...")
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		logger.Fatal("error do with tries postgresql")
	}
	logger.Info("database connect OK")

	logger.Debug("migrations...")
	err = Migration(ctx, logger, pool)
	logger.Info("migrations OK")

	return
}
