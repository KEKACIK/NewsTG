package source

import (
	"context"
	"newtg/pkg/client/postgresql"
	"newtg/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, source *Source) error {
	q := `
		INSERT INTO source (name)
		VALUES ($1)
		RETURNING id
	`
	r.logger.DebugSQL(q)

	err := r.client.QueryRow(ctx, q, source.Name).Scan(&source.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Get(ctx context.Context, id int) (Source, error) {
	panic("unimplemented")
}

func (r *repository) GetAll(ctx context.Context, source Source) ([]Source, error) {
	panic("unimplemented")
}

func (r *repository) Update(ctx context.Context, source *Source) error {
	panic("unimplemented")
}

func (r *repository) Delete(ctx context.Context, id int) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
