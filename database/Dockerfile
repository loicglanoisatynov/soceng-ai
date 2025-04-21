# front/Dockerfile
FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o server server.go

EXPOSE 8081

CMD ["./server"]
