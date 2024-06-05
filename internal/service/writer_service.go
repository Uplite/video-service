package service

import (
	"log"
	"net"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/grpc"

	"github.com/uplite/video-service/api/pb"
	"github.com/uplite/video-service/internal/config"
	"github.com/uplite/video-service/internal/server"
	"github.com/uplite/video-service/internal/storage"
	"github.com/uplite/video-service/internal/writer"
)

type videoWriterService struct {
	grpcServer   *grpc.Server
	writerServer pb.VideoServiceWriterServer
}

func NewVideoWriterService() *videoWriterService {
	c := s3.NewFromConfig(config.GetAwsConfig())
	g := grpc.NewServer()
	s := storage.NewS3Store(c, config.GetS3BucketName())
	w := writer.NewStoreWriter(s)

	writerServer := server.NewWriterServer(w)

	pb.RegisterVideoServiceWriterServer(g, writerServer)

	return &videoWriterService{
		grpcServer:   g,
		writerServer: writerServer,
	}
}

func (s *videoWriterService) Serve() error {
	lis, err := net.Listen("tcp", ":"+config.GetGrpcPort())
	if err != nil {
		log.Fatal(err)
	}

	return s.grpcServer.Serve(lis)
}

func (s *videoWriterService) Close() {
	s.grpcServer.GracefulStop()
}
