# 变量定义
IMAGE_NAME := sla2
IMAGE_TAG := prod-$(shell git describe --tags --always)
DOCKER_REGISTRY := registry.leeszi.cn/leeszi
FULL_IMAGE_NAME := $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
LATEST_IMAGE_NAME := $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest
BINARY_NAME := sla2

# 默认目标
.PHONY: all
all: build

# 本地编译
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME) .

build-arm:
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY_NAME)-arm64 .

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
	ACTIVE_PROFILE=local go run ./...

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
	rm -f api/proto/v1/*.pb.go
	docker-compose down -v

# 一键构建并推送
.PHONY: release
release: test docker-build docker-push

# 更新线上服务（需要配置 kubectl）
.PHONY: deploy
deploy:
	kubectl set image deployment/$(IMAGE_NAME) $(IMAGE_NAME)=$(FULL_IMAGE_NAME) -n default

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
	@echo "  make helm-install  - 安装 Helm 应用"
	@echo "  make helm-uninstall - 卸载 Helm 应用"
	@echo "  make helm-template - 生成 Helm 模板"
	@echo "  make helm-lint     - 验证 Helm 模板"

# Helm 相关命令
.PHONY: helm-install
helm-install:
	helm upgrade --install $(IMAGE_NAME) ./chart \
		--set image.repository=$(DOCKER_REGISTRY)/$(IMAGE_NAME) \
		--set image.tag=$(IMAGE_TAG) \
		-n default

.PHONY: helm-uninstall
helm-uninstall:
	helm uninstall $(IMAGE_NAME) -n default

.PHONY: helm-template
helm-template:
	helm template $(IMAGE_NAME) ./chart \
		--set image.repository=$(DOCKER_REGISTRY)/$(IMAGE_NAME) \
		--set image.tag=$(IMAGE_TAG)

.PHONY: helm-lint
helm-lint:
	helm lint ./chart 

# 本地运行服务
.PHONY: local-run
local-run:
	ACTIVE_PROFILE=local go run ./...

# 生成 protobuf 代码
.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative,require_unimplemented_servers=true \
		api/proto/v1/*.proto
