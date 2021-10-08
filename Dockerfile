FROM 1.16-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY client/ ./client
COPY server/ ./server
COPY proto/ ./proto
COPY utils/ ./utils
COPY .env ./server/.env
COPY .env ./client/.env
RUN ls -la
RUN go mod download
RUN go env GOOS GOARCH
WORKDIR /client
RUN PWD
RUN GOOS=linux GOARCH=arm64 go build -o ./CLIENT
WORKDIR /server
RUN GOOS=linux GOARCH=arm64 go build -o ./SERVER
WORKDIR /app
CMD ["./client/CLIENT && ./server/SERVER"]