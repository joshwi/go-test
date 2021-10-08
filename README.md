protoc --go_out=plugins=grpc:. *.proto

sudo docker build --progress=plain --no-cache -t node .