package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/uplite/video-service/api/pb"
)

const (
	videoKey = "user1/video_key2"
	userKey  = "user1"
)

func main() {
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewVideoServiceReaderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.GetOne(ctx, &pb.GetOneRequest{Key: videoKey})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)

	resp2, err := client.GetMany(ctx, &pb.GetManyRequest{UserPrefix: userKey})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp2)
}
