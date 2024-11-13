FROM golang:1.23-alpine as builder

COPY src /go/src/
ENV GOOS=linux GOARCH=amd64 GOPATH=/go

WORKDIR /go/src
RUN go build -a -v -o highload-sn-backend cmd/main.go

FROM alpine:3.20

EXPOSE 8080
COPY --from=builder /go/src/highload-sn-backend /usr/bin/

WORKDIR /usr/src/app/

ENTRYPOINT ["/usr/bin/highload-sn-backend"]
