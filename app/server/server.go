package server

import (
	"io"
	"log"
	"net"
	"os"

	proto "../proto"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func Env(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Load .env Error", err)
	}
	return os.Getenv(key)
}

type server struct{}

func main() {

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
			err := os.WriteFile(dir+"/files/"+filename, output, access)
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
