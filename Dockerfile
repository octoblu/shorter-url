FROM golang:1.6
MAINTAINER Octoblu, Inc. <docker@octoblu.com>

WORKDIR /go/src/github.com/octoblu/shorter-url
COPY . /go/src/github.com/octoblu/shorter-url

RUN env CGO_ENABLED=0 go build -o shorter-url -a -ldflags '-s' .

CMD ["./shorter-url"]
