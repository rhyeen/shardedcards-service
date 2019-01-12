FROM golang:alpine

RUN mkdir -p /go/src/github.com/rhyeen/shardedcards-service

WORKDIR /go/src/github.com/rhyeen/shardedcards-service

CMD ["go", "run", "cmd/server/main.go"]