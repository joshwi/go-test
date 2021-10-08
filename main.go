package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/joshwi/go-test/proto"
	"google.golang.org/grpc"
)

func Env(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Load .env Error", err)
	}
	return os.Getenv(key)
}

func InitClient() {

	URI := Env("URI")
	PORT := Env("PORT")
	filename := Env("FILENAME")

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

type server struct{}

func InitServer() {

	URI := Env("URI")
	PORT := Env("PORT")

	lis, err := net.Listen("tcp", URI+":"+PORT)
	if err != nil {
		log.Fatalf(" [ URI: %v ] [ Message: TCP listen failed ] [ Error: %v ]", URI+":"+PORT, err)
	}

	s := grpc.NewServer()
	proto.RegisterStreamServiceServer(s, &server{})

	log.Println("Server starting on port:", PORT)
	s.Serve(lis)
	if err := s.Serve(lis); err != nil {
		log.Fatalf(" [ URI: %v ] [ Message: GRPC server failed to start ] [ Error: %v ]", URI+":"+PORT, err)
	}

}

func (*server) StreamFile(stream proto.StreamService_StreamFileServer) error {

	var output []byte
	filename := ""
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Printf("[ File: %v ] [ Size: %v bytes ]", filename, len(output))
			if _, err := os.Stat(dir + "/files"); os.IsNotExist(err) {
				os.Mkdir(dir+"/files", 0777)
			}
			access := os.FileMode(0644)
			err := ioutil.WriteFile(dir+"/files/"+filename, output, access)
			if err != nil {
				log.Fatalf("[ File: %v ] [ Error: %v ]", filename, err)
			}
			return stream.SendAndClose(&proto.Response{File: filename, Size: int64(len(output)), Completed: true})
		}
		if err != nil {
			log.Fatalf("Failed to receive stream: %v", err)
		}
		filename = msg.File
		output = append(output, msg.Data...)
	}
}

func main() {

	var name string

	flag.StringVar(&name, `n`, `c`, `Specify server or client`)
	flag.Parse()

	if name == "c" {
		InitClient()
	} else {
		InitServer()
	}

}
