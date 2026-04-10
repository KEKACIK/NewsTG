package source

import "context"

type Repository interface {
	Create(ctx context.Context, source *Source) error
	GetAll(ctx context.Context, source Source) ([]Source, error)
	Get(ctx context.Context, id int) (Source, error)
	Update(ctx context.Context, source *Source) error
	Delete(ctx context.Context, id int) error
}
