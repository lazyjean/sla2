# 变量定义
IMAGE_NAME := sla2
IMAGE_TAG := $(shell git describe --tags --always)
DOCKER_REGISTRY := leeszi
FULL_IMAGE_NAME := $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
LATEST_IMAGE_NAME := $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest
BINARY_NAME := sla2

# 默认目标
.PHONY: all
all: build

# 本地编译
.PHONY: build
build:
	go build -o bin/$(BINARY_NAME) .

# 本地构建镜像
.PHONY: docker-build-local
docker-build-local: build
	docker build -t $(FULL_IMAGE_NAME) -f Dockerfile.dev .
	docker tag $(FULL_IMAGE_NAME) $(LATEST_IMAGE_NAME)

# 远程构建镜像（使用多阶段构建）
.PHONY: docker-build
docker-build:
	docker build -t $(FULL_IMAGE_NAME) -f Dockerfile .
	docker tag $(FULL_IMAGE_NAME) $(LATEST_IMAGE_NAME)

# 运行测试
.PHONY: test
test:
	go test -v ./...

# 推送 Docker 镜像
.PHONY: docker-push
docker-push:
	docker push $(FULL_IMAGE_NAME)
	docker push $(LATEST_IMAGE_NAME)

# 本地运行服务
.PHONY: run
run:
	docker-compose up -d

# 停止服务
.PHONY: stop
stop:
	docker-compose down

# 查看日志
.PHONY: logs
logs:
	docker-compose logs -f app

# 清理构建产物
.PHONY: clean
clean:
	rm -f $(IMAGE_NAME)
	docker-compose down -v

# 一键构建并推送
.PHONY: release
release: test docker-build docker-push

# 更新线上服务（需要配置 kubectl）
.PHONY: deploy
deploy:
	kubectl set image deployment/$(IMAGE_NAME) $(IMAGE_NAME)=$(FULL_IMAGE_NAME) -n your-namespace

# 帮助信息
.PHONY: help
help:
	@echo "可用的 make 命令："
	@echo "  make build         - 构建二进制文件"
	@echo "  make test          - 运行测试"
	@echo "  make docker-build  - 构建 Docker 镜像"
	@echo "  make docker-push   - 推送 Docker 镜像到仓库"
	@echo "  make run           - 本地运行服务"
	@echo "  make stop          - 停止服务"
	@echo "  make logs          - 查看服务日志"
	@echo "  make clean         - 清理构建产物"
	@echo "  make release       - 构建并推送镜像"
	@echo "  make deploy        - 更新线上服务" 