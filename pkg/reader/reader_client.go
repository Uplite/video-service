package reader

import (
	"context"

	"google.golang.org/grpc"

	"github.com/uplite/video-service/api/pb"
)

type readerClient struct {
	conn   *grpc.ClientConn
	client pb.VideoServiceReaderClient
}

func New(conn *grpc.ClientConn) *readerClient {
	return &readerClient{
		conn:   conn,
		client: pb.NewVideoServiceReaderClient(conn),
	}
}

func (c *readerClient) GetOne(ctx context.Context, req *pb.GetOneRequest, opts ...grpc.CallOption) (*pb.GetOneResponse, error) {
	return c.client.GetOne(ctx, req, opts...)
}

func (c *readerClient) GetMany(ctx context.Context, req *pb.GetManyRequest, opts ...grpc.CallOption) (*pb.GetManyResponse, error) {
	return c.client.GetMany(ctx, req, opts...)
}

func (c *readerClient) Close() error {
	return c.conn.Close()
}
