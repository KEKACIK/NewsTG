package parser

import "context"

type Parser interface {
	PoolNews(ctx context.Context)
}
