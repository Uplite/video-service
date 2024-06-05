package writer

import (
	"io"

	"github.com/uplite/video-service/api/pb"
)

type Client interface {
	pb.VideoServiceWriterClient
	io.Closer
}

var _ Client = (*writerClient)(nil)
