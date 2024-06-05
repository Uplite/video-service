package server

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/uplite/video-service/api/pb"
	"github.com/uplite/video-service/internal/videoutil"
	"github.com/uplite/video-service/internal/writer"
)

const (
	ErrNoContentType = "content_type cannot be empty"
	ErrNoKey         = "key cannot be empty"
)

type writerServer struct {
	pb.UnimplementedVideoServiceWriterServer

	writer writer.WriterDeleter
}

func NewWriterServer(writer writer.WriterDeleter) *writerServer {
	return &writerServer{writer: writer}
}

func newUploadError() *pb.UploadResponse {
	return &pb.UploadResponse{UploadStatus: pb.UploadStatus_UPLOAD_STATUS_ERROR}
}

func newUploadSuccess() *pb.UploadResponse {
	return &pb.UploadResponse{UploadStatus: pb.UploadStatus_UPLOAD_STATUS_SUCCESS}
}

func (s *writerServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := s.writer.Delete(ctx, req.GetKey()); err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{Ok: true}, nil
}

func (s *writerServer) Upload(stream pb.VideoServiceWriter_UploadServer) error {
	ctx := stream.Context()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	var buf bytes.Buffer
	var videoKey string
	var contentType string

	for {
		msg, recvErr := stream.Recv()
		if recvErr != nil {
			if recvErr == io.EOF {
				break
			}
			return recvErr
		}

		if videoKey == "" {
			videoKey = msg.GetKey()
		}

		if contentType == "" {
			contentType = videoutil.ContentTypeFrom(msg.GetContentType())
		}

		buf.Write(msg.GetData())
	}

	if videoKey == "" {
		return errors.New(ErrNoKey)
	}

	if contentType == "" {
		return errors.New(ErrNoContentType)
	}

	if err := s.writer.Write(ctx, videoKey, contentType, &buf); err != nil {
		return stream.SendAndClose(newUploadError())
	}

	return stream.SendAndClose(newUploadSuccess())
}
