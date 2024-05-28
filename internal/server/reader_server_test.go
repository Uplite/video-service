package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/uplite/video-service/api/pb"
	"github.com/uplite/video-service/internal/config"
)

type mockReader struct{}

func (m *mockReader) ReadOne(ctx context.Context, key string) (string, error) {
	return config.GetCloudFrontURL() + "/" + "video_1", nil
}

func (m *mockReader) ReadMany(ctx context.Context, prefix string) ([]string, error) {
	v1 := config.GetCloudFrontURL() + "/" + "video_1"
	v2 := config.GetCloudFrontURL() + "/" + "video_2"

	return []string{v1, v2}, nil
}

func TestReaderServer(t *testing.T) {
	grpcServer := grpc.NewServer()

	readerServer := newReaderServer(&mockReader{})
	readerServer.registerServer(grpcServer)

	lis, err := net.Listen("tcp", ":"+config.GetGrpcPort())
	if err != nil {
		t.Fatal(err)
	}

	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	conn, err := grpc.NewClient(":"+config.GetGrpcPort(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewVideoServiceReaderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	t.Run("should get one", func(t *testing.T) {
		res, err := client.GetOne(ctx, &pb.GetOneRequest{Key: "video_1"})
		assert.NoError(t, err)
		assert.Equal(t, res.GetUrl(), config.GetCloudFrontURL()+"/"+"video_1")
	})

	t.Run("should get many", func(t *testing.T) {
		res, err := client.GetMany(ctx, &pb.GetManyRequest{UserPrefix: "user1"})
		assert.NoError(t, err)
		assert.Len(t, res.GetUrls(), 2)
		assert.Equal(t, res.GetUrls()[0], config.GetCloudFrontURL()+"/"+"video_1")
		assert.Equal(t, res.GetUrls()[1], config.GetCloudFrontURL()+"/"+"video_2")
	})
}
