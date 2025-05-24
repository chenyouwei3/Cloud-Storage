FROM golang:1.20 AS builder


ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app/cloud-storage

COPY go.mod go.sum ./

RUN go mod download
COPY . .


RUN go build -o web ./cmd/main.go

# 第二阶段：运行阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /app/cloud-storage

COPY --from=builder /app/cloud-storage/web .

# 声明服务端口
EXPOSE 8080
EXPOSE 6066

# 运行应用
ENTRYPOINT ["./web"]
