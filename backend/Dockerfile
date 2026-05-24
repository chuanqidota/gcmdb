FROM docker.1ms.run/library/golang:1.25.10 AS builder
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

FROM docker.1ms.run/library/alpine:3.22.4
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /opt/data
COPY --from=builder /opt/data/gcmdb .
COPY --from=builder /opt/data/config/ ./config/
RUN chmod +x gcmdb
EXPOSE 8080

ENTRYPOINT ["./gcmdb"]
