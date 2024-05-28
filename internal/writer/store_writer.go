package writer

import (
	"bytes"
	"context"

	"github.com/uplite/video-service/internal/storage"
)

type storeWriter struct {
	store storage.Store
}

func NewStoreWriter(store storage.Store) *storeWriter {
	return &storeWriter{
		store: store,
	}
}

func (r *storeWriter) Write(ctx context.Context, key, contentType string, buf *bytes.Buffer) error {
	return r.store.Put(ctx, key, contentType, buf)
}

func (r *storeWriter) Delete(ctx context.Context, key string) error {
	return r.store.Delete(ctx, key)
}
