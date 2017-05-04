FROM scratch
MAINTAINER Octoblu, Inc. <docker@octoblu.com>
EXPOSE 80

HEALTHCHECK CMD ./healthchecker --url http://localhost:80/healthcheck || exit 1

ADD https://github.com/octoblu/go-http-healthcheck/releases/download/v1.0.2/http-healthcheck-linux-amd64 healthchecker
RUN chmod +x healthchecker

ADD entrypoint entrypoint
ENTRYPOINT ["./entrypoint"]
