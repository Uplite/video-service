package service

import (
	"log"
	"net"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/grpc"

	"github.com/uplite/video-service/api/pb"
	"github.com/uplite/video-service/internal/config"
	"github.com/uplite/video-service/internal/reader"
	"github.com/uplite/video-service/internal/server"
	"github.com/uplite/video-service/internal/storage"
)

type videoReaderService struct {
	grpcServer   *grpc.Server
	readerServer pb.VideoServiceReaderServer
}

func NewVideoReaderService() *videoReaderService {
	c := s3.NewFromConfig(config.GetAwsConfig())
	g := grpc.NewServer()
	s := storage.NewS3Store(c, config.GetS3BucketName())
	r := reader.NewStoreReader(s)

	readerServer := server.NewReaderServer(r)

	pb.RegisterVideoServiceReaderServer(g, readerServer)

	return &videoReaderService{
		grpcServer:   g,
		readerServer: readerServer,
	}
}

func (s *videoReaderService) Serve() error {
	lis, err := net.Listen("tcp", ":"+config.GetGrpcPort())
	if err != nil {
		log.Fatal(err)
	}

	return s.grpcServer.Serve(lis)
}

func (s *videoReaderService) Close() {
	s.grpcServer.GracefulStop()
}
