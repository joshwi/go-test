package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/joshwi/go-test/proto"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterStreamServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// SumAll sums up all the numbers recieved from the client stream
func (*server) SumAll(stream proto.StreamService_StreamFileServer) error {
	var sum string

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.Response{Success: true})
		}
		if err != nil {
			log.Fatalf("could not recieve stream: %v", err)
		}
		sum = sum + fmt.Sprintf("%v", msg.Data)
	}

	log.Println(sum)

	return nil
}
