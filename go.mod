module github.com/joshwi/go-test

go 1.16

replace github.com/joshwi/go-test/app/proto => ./proto

replace github.com/joshwi/go-test/app/utils => ./utils

require (
	github.com/joho/godotenv v1.4.0 // indirect
	google.golang.org/grpc v1.41.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
