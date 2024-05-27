package recorder

import (
	"bytes"
	"context"

	"github.com/uplite/video-service/internal/storage"
)

type storeRecorder struct {
	store storage.Store
}

func NewStoreRecorder(store storage.Store) *storeRecorder {
	return &storeRecorder{
		store: store,
	}
}

func (r *storeRecorder) Record(ctx context.Context, key string, buf *bytes.Buffer) error {
	return r.store.Put(ctx, key, buf)
}
