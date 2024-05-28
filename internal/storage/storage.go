package storage

import (
	"context"
	"io"
)

type Store interface {
	Put(ctx context.Context, key, contentType string, data io.Reader) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Head(ctx context.Context, key string) error
	Delete(ctx context.Context, key string) error
	List(ctx context.Context, prefix string) ([]string, error)
}
