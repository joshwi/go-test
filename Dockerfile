FROM golang:1.17
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY app /app
COPY .env /app/server/.env
COPY .env /app/client/.env
RUN go mod download
RUN go env GOOS GOARCH
RUN pwd
RUN ls -la
WORKDIR /client
RUN pwd
RUN ls -la
RUN GOOS=linux GOARCH=arm64 go build client.go
WORKDIR /server
RUN GOOS=linux GOARCH=arm64 go build server.go
WORKDIR /app
CMD ["./client/client && ./server/server"]