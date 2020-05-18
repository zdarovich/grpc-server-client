FROM golang:1.13-alpine as builder

RUN mkdir /build
ADD . /build/

WORKDIR /build

RUN go build -o grpc-server cmd/server/main.go

FROM alpine:latest

COPY --from=builder /build/grpc-server /app/
WORKDIR /app

CMD ["./grpc-server"]
