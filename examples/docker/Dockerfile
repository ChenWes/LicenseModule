FROM golang:1.19-alpine AS builder

WORKDIR /app

# 复制go mod文件
COPY go.* ./
RUN go mod download

# 复制源代码
COPY . .

# 编译生成工具
RUN CGO_ENABLED=0 GOOS=linux go build -o /license-generator ./cmd/license/generate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /example-app ./examples/app/main.go

# 使用更小的基础镜像
FROM alpine:latest

WORKDIR /app

# 复制可执行文件
COPY --from=builder /license-generator .
COPY --from=builder /example-app .

# 创建配置目录
RUN mkdir -p ./config

# 设置环境变量
ENV LICENSE_PATH=/app/config/license.dat
ENV TIMESTAMP_PATH=/app/config/timestamp.dat

# 设置入口点
ENTRYPOINT ["/app/example-app"] 