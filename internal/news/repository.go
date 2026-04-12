package news

import (
	"context"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"
	"time"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, new *News) error {
	q := `
		INSERT INTO news
			(title, link, content, source_id, published)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id, created
	`
	r.logger.DebugSQL(q)

	err := r.client.QueryRow(ctx, q, new.Title, new.Link, new.Content, new.Source.ID, time.Now()).Scan(&new.ID, &new.Created)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetAll(ctx context.Context) ([]News, error) {
	q := `
		SELECT
			id, title, link, content, source_id, posted, published, created
		FROM news
	`
	r.logger.DebugSQL(q)

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	news := make([]News, 0)

	for rows.Next() {
		var new News

		err = rows.Scan(&new.ID, &new.Title, &new.Link, &new.Content, &new.Source.ID, &new.Posted, &new.Published, &new.Created)
		if err != nil {
			return nil, err
		}

		news = append(news, new)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return news, nil
}

func (r *repository) GetAllByPost(ctx context.Context, posted bool) ([]News, error) {
	q := `
		SELECT
			id, title, link, content, source_id, posted, published, created
		FROM news
		WHERE posted = $1
	`
	r.logger.DebugSQL(q)

	rows, err := r.client.Query(ctx, q, posted)
	if err != nil {
		return nil, err
	}

	news := make([]News, 0)

	for rows.Next() {
		var new News

		err = rows.Scan(&new.ID, &new.Title, &new.Link, &new.Content, &new.Source.ID, &new.Posted, &new.Published, &new.Created)
		if err != nil {
			return nil, err
		}

		news = append(news, new)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return news, nil
}

func (r *repository) Get(ctx context.Context, id int) (News, error) {
	q := `
		SELECT
			id, title, link, content, source_id, posted, published, created
		FROM news
		WHERE id = $1
	`
	r.logger.DebugSQL(q)

	var new News
	err := r.client.QueryRow(ctx, q, id).Scan(&new.ID, &new.Title, &new.Link, &new.Content, &new.Source, new.Posted, new.Published, new.Created)
	if err != nil {
		return News{}, err
	}

	return new, nil
}

// TODO
func (r *repository) Update(ctx context.Context, new *News) error {
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
