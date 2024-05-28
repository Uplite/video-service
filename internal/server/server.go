package server

import (
	"log"
	"net"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/grpc"

	"github.com/uplite/video-service/internal/config"
	"github.com/uplite/video-service/internal/reader"
	"github.com/uplite/video-service/internal/storage"
	"github.com/uplite/video-service/internal/writer"
)

type videoWriterServer struct {
	grpcServer   *grpc.Server
	writerServer *writerServer
}

type videoReaderServer struct {
	grpcServer   *grpc.Server
	readerServer *readerServer
}

func NewWriter() *videoWriterServer { return newVideoWriterServer() }

func NewReader() *videoReaderServer { return newVideoReaderServer() }

func newVideoWriterServer() *videoWriterServer {
	client := s3.NewFromConfig(config.GetAwsConfig())

	grpcServer := grpc.NewServer()

	writerServer := newWriterServer(writer.NewStoreWriter(storage.NewS3Store(client, config.GetS3BucketName())))
	writerServer.registerServer(grpcServer)

	return &videoWriterServer{
		grpcServer:   grpcServer,
		writerServer: writerServer,
	}
}

func newVideoReaderServer() *videoReaderServer {
	client := s3.NewFromConfig(config.GetAwsConfig())

	grpcServer := grpc.NewServer()

	readerServer := newReaderServer(reader.NewStoreReader(storage.NewS3Store(client, config.GetS3BucketName())))
	readerServer.registerServer(grpcServer)

	return &videoReaderServer{
		grpcServer:   grpcServer,
		readerServer: readerServer,
	}
}

func (s *videoWriterServer) Serve() error {
	lis, err := net.Listen("tcp", ":"+config.GetGrpcPort())
	if err != nil {
		log.Fatal(err)
	}

	return s.grpcServer.Serve(lis)
}

func (s *videoWriterServer) Close() {
	s.grpcServer.GracefulStop()
}

func (s *videoReaderServer) Serve() error {
	lis, err := net.Listen("tcp", ":"+config.GetGrpcPort())
	if err != nil {
		log.Fatal(err)
	}

	return s.grpcServer.Serve(lis)
}

func (s *videoReaderServer) Close() {
	s.grpcServer.GracefulStop()
}
