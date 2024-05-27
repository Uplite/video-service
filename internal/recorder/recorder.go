package recorder

import (
	"bytes"
	"context"
)

type Recorder interface {
	Record(ctx context.Context, key string, buf *bytes.Buffer) error
}
