# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev

# 设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 生成 swagger 文档
RUN make docs

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o sla2 .

# 运行阶段
FROM alpine:3.19

WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai

# 从构建阶段复制二进制文件和 swagger 文件
COPY --from=builder /app/sla2 .
COPY --from=builder /app/api/swagger /app/api/swagger

EXPOSE 9000 8080

CMD ["./sla2"]