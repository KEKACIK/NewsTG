package news

import "context"

type Repository interface {
	Create(ctx context.Context, new *News) error
	GetAll(ctx context.Context) ([]News, error)
	Get(ctx context.Context, id int) (News, error)
	Update(ctx context.Context, new *News) error
	Delete(ctx context.Context, id int) error
}
