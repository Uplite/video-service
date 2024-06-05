package reader

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/uplite/video-service/api/pb"
	"github.com/uplite/video-service/internal/server"
)

type mockReader struct{}

func (r *mockReader) ReadOne(ctx context.Context, key string) (string, error) {
	return "video_url1", nil
}

func (r *mockReader) ReadMany(ctx context.Context, prefix string) ([]string, error) {
	return []string{"user1/video_url1", "user1/video_url2"}, nil
}

func TestSnowflakeClient(t *testing.T) {
	srv := server.NewReaderServer(new(mockReader))

	grpcServer := grpc.NewServer()

	pb.RegisterVideoServiceReaderServer(grpcServer, srv)

	lis, err := net.Listen("tcp", ":50053")
	assert.NoError(t, err)

	go grpcServer.Serve(lis)
	defer grpcServer.Stop()

	conn, err := grpc.NewClient(":50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	c := New(conn)

	t.Run("should get one video url", func(t *testing.T) {
		res, err := c.GetOne(context.Background(), &pb.GetOneRequest{Key: "video_url1"})
		assert.NoError(t, err)
		assert.Equal(t, "video_url1", res.GetUrl())
	})

	t.Run("should get many video urls", func(t *testing.T) {
		res, err := c.GetMany(context.Background(), &pb.GetManyRequest{UserPrefix: "user1"})
		assert.NoError(t, err)
		assert.Equal(t, []string{"user1/video_url1", "user1/video_url2"}, res.GetUrls())
	})
}
