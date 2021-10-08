module github.com/joshwi/go-test.git

go 1.16

replace github.com/joshwi/go-test/proto => ./proto

require (
	github.com/golang/protobuf v1.5.2 // indirect
	google.golang.org/grpc v1.41.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
