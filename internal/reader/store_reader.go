package reader

import (
	"context"
	"strings"

	"github.com/uplite/video-service/internal/config"
	"github.com/uplite/video-service/internal/storage"
)

type storeReader struct {
	store storage.Store
}

func NewStoreReader(store storage.Store) *storeReader {
	return &storeReader{
		store: store,
	}
}

func (r *storeReader) ReadOne(ctx context.Context, key string) (string, error) {
	if err := r.store.Head(ctx, key); err != nil {
		return "", err
	}
	return appendDelimiter(config.GetCloudFrontURL()) + key, nil
}

func (r *storeReader) ReadMany(ctx context.Context, prefix string) ([]string, error) {
	var us []string

	keys, err := r.store.List(ctx, appendDelimiter(prefix))
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		us = append(us, appendDelimiter(config.GetCloudFrontURL())+key)
	}

	return us, nil
}

func appendDelimiter(input string) string {
	if !strings.HasSuffix(input, "/") {
		return input + "/"
	}
	return input
}
