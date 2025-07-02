FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go mod init fetcher && go mod tidy
RUN go build -o fetcher src/main.go

CMD ["./fetcher"]
