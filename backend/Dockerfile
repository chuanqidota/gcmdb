FROM golang:1.21-alpine AS builder
RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env -w GOOS=linux \
    && go env -w GOARCH=amd64 \
    && go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /opt/data/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gcmdb main.go

FROM alpine:3
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /opt/data
COPY --from=builder /opt/data/gcmdb .
RUN chmod +x gcmdb
EXPOSE 8080

ENTRYPOINT ["./gcmdb"]
