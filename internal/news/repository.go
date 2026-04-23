package news

import (
	"context"
	"fmt"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"
	"strings"
	"time"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, dto *CreateDTO) error {
	q := `
		INSERT INTO news
			(title, link, content, source_id, likes, published_at)
		VALUES
			($1, $2, $3, $4, $5, $6)
		ON CONFLICT (link) DO UPDATE SET
			title = EXCLUDED.title, content = EXCLUDED.content, likes = EXCLUDED.likes
	`
	r.logger.DebugSQL(q, dto.Title, dto.Link, dto.Content, dto.SourceID, dto.Likes, dto.Published)

	_, err := r.client.Exec(ctx, q, dto.Title, dto.Link, dto.Content, dto.SourceID, dto.Likes, dto.Published)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetAll(ctx context.Context, dto *GetAllDTO) ([]News, error) {
	args := make([]any, 0)
	where_list := make([]string, 0)

	{
		zeroTime := time.Time{}

		if dto.Status != "" {
			args = append(args, dto.Status)
			where_list = append(where_list, fmt.Sprintf("status=$%d", len(args)))
		}

		if dto.FromDate != zeroTime {
			args = append(args, dto.FromDate)
			where_list = append(where_list, fmt.Sprintf("published_at >= $%d", len(args)))
		}

		if dto.ToDate != zeroTime {
			args = append(args, dto.ToDate)
			where_list = append(where_list, fmt.Sprintf("published_at < $%d", len(args)))
		}
	}

	if dto.Limit == 0 {
		dto.Limit = 10_000
	}
	args = append(args, dto.Limit)

	q := fmt.Sprintf(
		`SELECT %s FROM news %s ORDER BY likes DESC LIMIT $%d`,
		"id, title, link, content, source_id, likes, status, published_at, created_at",
		fmt.Sprintf("WHERE %s", strings.Join(where_list, " AND ")),
		len(args),
	)
	r.logger.DebugSQL(q, args...)

	rows, err := r.client.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	news := make([]News, 0)
	for rows.Next() {
		var new News

		err = rows.Scan(&new.ID, &new.Title, &new.Link, &new.Content, &new.Source.ID, &new.Likes, &new.Status, &new.Published, &new.Created)
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
	r.logger.DebugSQL(q, id)

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
	r.logger.DebugSQL(q, new.Title, new.Link, new.Content, new.Source.ID, new.Status, new.ID)

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
