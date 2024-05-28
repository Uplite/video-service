package reader

import (
	"context"
)

type Reader interface {
	ReadOne(ctx context.Context, key string) (string, error)
	ReadMany(ctx context.Context, prefix string) ([]string, error)
}
