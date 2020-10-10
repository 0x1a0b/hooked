FROM golang:1.14.7-buster

ENV GO111MODULE="on" \
    CONFIG_DIR="/opt/config/config.yml"

WORKDIR /go/src/app

CMD ["./app"]

COPY . .
COPY config.yml /opt/config/config.yml

RUN go vet && go build -o app

