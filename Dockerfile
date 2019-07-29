FROM golang:1.12-alpine3.10

WORKDIR /go/src/route-dumper
COPY . .
RUN go build main.go

ENTRYPOINT ["tail", "-f", "/dev/null"]