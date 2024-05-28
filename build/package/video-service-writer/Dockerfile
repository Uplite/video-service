# Builder
FROM golang:latest AS builder

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apt-get update && apt-get install -y protobuf-compiler

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make clean

RUN make generate

RUN make build-writer

# Final
FROM scratch

COPY --from=builder /app/bin/video-service-writer /video-service-writer

EXPOSE 50051

CMD ["/video-service-writer"]
