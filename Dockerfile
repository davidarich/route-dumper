FROM golang:1.12-alpine3.10

WORKDIR /go/src/route-dumper

RUN apk add git

COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on go mod download
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 go build -o route-dumper

ENTRYPOINT ["tail", "-f", "/dev/null"]