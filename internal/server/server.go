package server

import (
	"log"
	"net"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/grpc"

	"github.com/uplite/video-service/internal/config"
	"github.com/uplite/video-service/internal/recorder"
	"github.com/uplite/video-service/internal/storage"
)

type server struct {
	grpcServer  *grpc.Server
	videoServer *videoServer
}

func New() *server {
	client := s3.NewFromConfig(config.GetAwsConfig())

	grpcServer := grpc.NewServer()

	videoServer := newVideoServer(recorder.NewStoreRecorder(storage.NewS3Store(client, config.GetS3BucketName())))
	videoServer.registerServer(grpcServer)

	return &server{
		grpcServer:  grpcServer,
		videoServer: videoServer,
	}
}

func (s *server) Serve() error {
	lis, err := net.Listen("tcp", ":"+config.GetGrpcPort())
	if err != nil {
		log.Fatal(err)
	}

	return s.grpcServer.Serve(lis)
}

func (s *server) Close() {
	s.grpcServer.GracefulStop()
}
