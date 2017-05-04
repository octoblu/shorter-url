FROM golang:1.8
MAINTAINER Octoblu, Inc. <docker@octoblu.com>
WORKDIR /go/src/github.com/octoblu/shorter-url

EXPOSE 80
HEALTHCHECK CMD ./healthchecker --url http://localhost:80/healthcheck || exit 1

ADD https://github.com/octoblu/go-http-healthcheck/releases/download/v1.0.2/http-healthcheck-linux-amd64 healthchecker
RUN chmod +x healthchecker

COPY . /go/src/github.com/octoblu/shorter-url

RUN env CGO_ENABLED=0 go build -o shorter-url -a -ldflags '-s' .

CMD ["./shorter-url"]
