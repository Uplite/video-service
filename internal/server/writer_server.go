package server

import (
	"bytes"
	"context"
	"io"

	"google.golang.org/grpc"

	"github.com/uplite/video-service/api/pb"
	"github.com/uplite/video-service/internal/videoutil"
	"github.com/uplite/video-service/internal/writer"
)

type writerServer struct {
	pb.UnimplementedVideoServiceWriterServer
	writer writer.WriterDeleter
}

func newWriterServer(writer writer.WriterDeleter) *writerServer {
	return &writerServer{writer: writer}
}

func newUploadError() *pb.UploadResponse {
	return &pb.UploadResponse{UploadStatus: pb.UploadStatus_UPLOAD_STATUS_ERROR}
}

func newUploadSuccess() *pb.UploadResponse {
	return &pb.UploadResponse{UploadStatus: pb.UploadStatus_UPLOAD_STATUS_SUCCESS}
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
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if videoKey == "" {
			videoKey = msg.GetKey()
		}

		if contentType == "" {
			contentType = videoutil.ContentTypeFrom(msg.GetContentType())
		}

		buf.Write(msg.GetData())
	}

	if err := s.writer.Write(ctx, videoKey, contentType, &buf); err != nil {
		if sendErr := stream.SendAndClose(newUploadError()); sendErr != nil {
			return sendErr
		}
		return err
	}

	if err := stream.SendAndClose(newUploadSuccess()); err != nil {
		return err
	}

	return nil
}

func (s *writerServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if err := s.writer.Delete(ctx, req.GetKey()); err != nil {
		return nil, err
	}
	return &pb.DeleteResponse{Ok: true}, nil
}

func (s *writerServer) registerServer(g *grpc.Server) {
	pb.RegisterVideoServiceWriterServer(g, s)
}
