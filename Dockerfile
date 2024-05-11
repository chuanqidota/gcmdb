FROM golang:1.21.4 AS builder
RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env -w GOOS=linux \
    && go env -w GOPROXY=https://goproxy.cn,direct
RUN mkdir -p /root/data/
WORKDIR /root/data/
COPY . .
RUN go mod tidy
RUN go build -o gcmdb main.go

FROM alpine:3
RUN mkdir -p /root/data
WORKDIR /root/data
COPY --from=builder /root/data/rbac .
WORKDIR /root/data
RUN chmod +x gcmdb
EXPOSE 8080

ENTRYPOINT ["./gcmdb"]

# docker build -t gcmdb:v1.0.0 .
# docker run --name gcmdb -d -p 8080:8080 gcmdb:v1.0.0