package reader

import (
	"io"

	"github.com/uplite/video-service/api/pb"
)

type Client interface {
	pb.VideoServiceReaderClient
	io.Closer
}

var _ Client = (*readerClient)(nil)
