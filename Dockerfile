FROM golang:1.17
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY app /app
RUN go mod download
RUN go env GOOS GOARCH
RUN pwd
RUN ls -la
RUN GOOS=linux GOARCH=arm64 go build -o ./grpc
RUN pwd
RUN ls -la
CMD ["./grpc -t='s' && ./grpc -t='c'"]