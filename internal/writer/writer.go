package writer

import (
	"bytes"
	"context"
)

type Writer interface {
	Write(ctx context.Context, key, contentType string, buf *bytes.Buffer) error
}

type Deleter interface {
	Delete(ctx context.Context, key string) error
}

type WriterDeleter interface {
	Writer
	Deleter
}
