FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go mod init fetcher-go && go mod tidy
RUN go build -o target src/main.go

CMD ["./target"]
