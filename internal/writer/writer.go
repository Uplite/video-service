package writer

import (
	"bytes"
	"context"
)

type Writer interface {
	Write(ctx context.Context, key string, buf *bytes.Buffer) error
}
