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
			id, title, link, content, source_id, status, published, created
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

		err = rows.Scan(&new.ID, &new.Title, &new.Link, &new.Content, &new.Source.ID, &new.Status, &new.Published, &new.Created)
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

func (r *repository) GetAllByStatus(ctx context.Context, status NewStatus) ([]News, error) {
	q := `
		SELECT
			id, title, link, content, source_id, status, published, created
		FROM news
		WHERE status = $1
	`
	r.logger.DebugSQL(q)

	rows, err := r.client.Query(ctx, q, status)
	if err != nil {
		return nil, err
	}

	news := make([]News, 0)

	for rows.Next() {
		var new News

		err = rows.Scan(&new.ID, &new.Title, &new.Link, &new.Content, &new.Source.ID, &new.Status, &new.Published, &new.Created)
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
			id, title, link, content, source_id, status, published, created
		FROM news
		WHERE id = $1
	`
	r.logger.DebugSQL(q)

	var new News
	err := r.client.QueryRow(ctx, q, id).Scan(&new.ID, &new.Title, &new.Link, &new.Content, &new.Source, new.Status, new.Published, new.Created)
	if err != nil {
		return News{}, err
	}

	return new, nil
}

// TODO
func (r *repository) Update(ctx context.Context, new *News) (err error) {
	q := `
		UPDATE news SET
			title=$1, link=$2, content=$3, source_id=$4, status=$5
		WHERE id = $6
	`
	r.logger.DebugSQL(q)

	_, err = r.client.Exec(ctx, q, new.Title, new.Link, new.Content, new.Source.ID, new.Status, new.ID)

	return err
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
