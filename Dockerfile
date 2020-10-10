FROM golang:1.14.7-buster

ENV GO111MODULE=on \
    CONFIG_DIR=/go/src/app

WORKDIR /go/src/app

CMD ["./app"]

COPY . .

RUN go vet && go build -o app

