package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/planetscale/vtprotobuf/codec/grpc"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"

	"github.com/uplite/video-service/internal/server"
)

func init() {
	encoding.RegisterCodec(grpc.Codec{})
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	s := server.New()

	go func() {
		if err := s.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	s.Close()
}
