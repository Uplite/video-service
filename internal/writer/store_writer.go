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

func (r *storeWriter) Write(ctx context.Context, key string, buf *bytes.Buffer) error {
	return r.store.Put(ctx, key, buf)
}
