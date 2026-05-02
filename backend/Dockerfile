FROM golang:1.21.6 AS builder
RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env -w GOOS=linux \
    && go env -w GOARCH=amd64 \
    && go env -w GOPROXY=https://goproxy.cn,direct
RUN mkdir -p /opt/data/
WORKDIR /opt/data/
COPY . .
RUN go mod tidy
RUN go build -o gcmdb main.go

FROM alpine:3
RUN mkdir -p /opt/data
WORKDIR /opt/data
COPY --from=builder /opt/data/gcmdb .
WORKDIR /opt/data
RUN chmod +x gcmdb
EXPOSE 8080

ENTRYPOINT ["./gcmdb"]

# cd backend && docker build -t gcmdb:v1.0.0 .
# docker run --name gcmdb -d -p 8080:8080 gcmdb:v1.0.0