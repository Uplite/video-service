package server

import (
	"bytes"
	"io"

	"google.golang.org/grpc"

	"github.com/uplite/video-service/api/pb"
	"github.com/uplite/video-service/internal/recorder"
)

type videoServer struct {
	pb.UnimplementedVideoServiceWriterServer
	recorder recorder.Recorder
}

func newVideoServer(recorder recorder.Recorder) *videoServer {
	return &videoServer{recorder: recorder}
}

func newUploadError() *pb.UploadResponse {
	return &pb.UploadResponse{UploadStatus: pb.UploadStatus_UPLOAD_STATUS_ERROR}
}

func newUploadSuccess() *pb.UploadResponse {
	return &pb.UploadResponse{UploadStatus: pb.UploadStatus_UPLOAD_STATUS_SUCCESS}
}

func (s *videoServer) Upload(stream pb.VideoServiceWriter_UploadServer) error {
	ctx := stream.Context()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	var buf bytes.Buffer
	var videoKey string

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

		buf.Write(msg.GetData())
	}

	if err := s.recorder.Record(ctx, videoKey, &buf); err != nil {
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

func (s *videoServer) registerServer(g *grpc.Server) {
	pb.RegisterVideoServiceWriterServer(g, s)
}
