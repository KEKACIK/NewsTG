package source

import (
	"context"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, source *Source) error {
	q := `
		INSERT INTO source
			(name)
		VALUES
			($1)
	`
	r.logger.DebugSQL(q, source.Name)

	_, err := r.client.Exec(ctx, q, source.Name)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetAll(ctx context.Context) ([]Source, error) {
	q := `
		SELECT
			id, name
		FROM source
	`
	r.logger.DebugSQL(q)

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	sources := make([]Source, 0)

	for rows.Next() {
		var source Source

		err = rows.Scan(&source.ID, &source.Name)
		if err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}

func (r *repository) Get(ctx context.Context, id int) (Source, error) {
	q := `
		SELECT
			id, name
		FROM source WHERE
			id = $1
	`
	r.logger.DebugSQL(q, id)

	var source Source
	err := r.client.QueryRow(ctx, q, id).Scan(&source.ID, &source.Name)
	if err != nil {
		return Source{}, err
	}

	return source, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (Source, error) {
	q := `
		SELECT
			id, name
		FROM source WHERE
			name = $1
	`
	r.logger.DebugSQL(q, name)

	var source Source
	err := r.client.QueryRow(ctx, q, name).Scan(&source.ID, &source.Name)
	if err != nil {
		return Source{}, err
	}

	return source, nil
}

// TODO
func (r *repository) Update(ctx context.Context, source *Source) error {
	panic("unimplemented")
}

// TODO
func (r *repository) Delete(ctx context.Context, id int) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
