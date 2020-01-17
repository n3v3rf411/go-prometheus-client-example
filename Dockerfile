FROM golang:1.13.3-buster as build-env

ADD . /go/src/app
RUN cd /go/src/app && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod=vendor ./cmd/server

FROM alpine

WORKDIR /app

COPY --from=build-env /go/src/app/server /app/
COPY --from=build-env /go/src/app/public /app/public

CMD ["/app/server"]
