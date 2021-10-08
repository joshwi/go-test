package client

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/joshwi/go-test/app/proto"
	"github.com/joshwi/go-test/app/utils"
	"google.golang.org/grpc"
)

func main() {

	URI := utils.Env("URI")
	PORT := utils.Env("PORT")
	filename := utils.Env("FILENAME")

	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	conn, err := grpc.Dial(URI+":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := proto.NewStreamServiceClient(conn)

	stream, err := c.StreamFile(context.Background())
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}

	index := 0

	for {
		chunk := make([]byte, 64*1024)
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		stream.Send(&proto.Request{File: filename, Index: int64(index), Data: chunk[:n]})
		index++
	}

	// close the stream and recieve result
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to recieve response: %v", err)
	}

	log.Println(res)
}
