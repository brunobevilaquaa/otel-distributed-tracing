FROM golang:1.22 as Build

ARG TARGET_SERVICE

WORKDIR /go/src/service

RUN apt update

COPY . .

RUN go get -d -v ./... && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o /go/bin/service ./cmd/${TARGET_SERVICE}

FROM alpine:3.18

COPY --from=Build /go/bin/service /usr/local/bin/service

ENTRYPOINT ["/usr/local/bin/service"]