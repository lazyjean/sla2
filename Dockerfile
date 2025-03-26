# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev make

# 设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main .

# 最终阶段
FROM alpine:3.19

# 设置工作目录
WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/config ./config

# 复制 swagger.json 文件
COPY --from=builder /app/api/swagger/swagger.json ./api/swagger/swagger.json

# 设置环境变量
ENV ACTIVE_PROFILE=prod

# 暴露端口
EXPOSE 9101 9102 

CMD ["./main"]
