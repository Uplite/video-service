package main

import (
	"context"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/uplite/video-service/api/pb"
)

const (
	chunkSize = 3 * 1024 * 1024 // 3 MB
	videoKey  = "user1/video_key2"
	videoPath = "../../test/mov_bbb.mp4"
)

func main() {
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewVideoServiceWriterClient(conn)

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalf("failed to open stream: %v", err)
	}
	defer stream.CloseSend()

	file, err := os.Open(videoPath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, chunkSize)

	var chunkNum int
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatalf("failed to read from video file: %v", err)
		}
		if n == 0 {
			break
		}
		chunk := &pb.UploadRequest{
			Key:         videoKey,
			Data:        buffer[:n],
			ContentType: pb.VideoContentType_VIDEO_CONTENT_TYPE_MP4,
		}
		if err := stream.Send(chunk); err != nil {
			log.Fatalf("error sending chunk: %v", err)
		}
		chunkNum++
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receiving response: %v", err)
	}

	log.Printf("Server response: %s", resp.GetUploadStatus().String())

	resp2, err := client.Delete(context.Background(), &pb.DeleteRequest{Key: videoKey})
	if err != nil {
		log.Fatalf("error deleting vid: %v", err)
	}

	log.Printf("Server response: %v", resp2.GetOk())
}
