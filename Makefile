# 变量定义
IMAGE_NAME := sla2
DOCKER_REGISTRY := registry.leeszi.cn/leeszi

# 使用函数而不是直接赋值，确保每次使用时都重新获取最新的tag
GET_LATEST_TAG = $(shell git tag -l | sort -V | tail -n1)
GET_FULL_IMAGE = $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(GET_LATEST_TAG)

BINARY_NAME := sla2

# 默认目标
.PHONY: all
all: build

# 安装 air 工具
.PHONY: install-air
install-air:
	@echo "Installing air..."
	@go install github.com/air-verse/air@latest
	@echo "Air installed successfully"

# 使用 air 运行服务（热重载）
.PHONY: run-hot
run-hot:
	@echo "Starting service with hot reload..."
	@air -c .air.toml

# 生成 protobuf 代码
.PHONY: proto
proto:
	@echo "生成 protobuf 代码..."
	@cd sla2-proto && buf generate
	@mkdir -p api/proto/v1
	@rm -f api/proto/v1/*.pb.go api/proto/v1/*.gw.go
	@cp -r sla2-proto/gen/go/proto/v1/* api/proto/v1/
	@echo "protobuf 代码生成完成"

# 本地编译
.PHONY: build
build: proto
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME) .

build-arm:
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY_NAME)-arm64 .

# 本地构建镜像
.PHONY: docker-build-local
docker-build-local: build
	docker build -t $(GET_FULL_IMAGE) -f Dockerfile.dev .
	docker tag $(GET_FULL_IMAGE) $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest

# 远程构建镜像（使用多阶段构建）
.PHONY: docker-build
docker-build:
	@echo "Building docker image: $(GET_FULL_IMAGE)"
	docker build -t $(GET_FULL_IMAGE) --platform linux/amd64 .
	docker tag $(GET_FULL_IMAGE) $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest

# 运行测试
.PHONY: test
test:
	go test -v ./...

# 推送 Docker 镜像
.PHONY: docker-push
docker-push:
	docker push $(GET_FULL_IMAGE)
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):latest

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
	rm -f api/proto/v1/*.pb.go api/proto/v1/*.gw.go
	@cd sla2-proto && rm -rf gen/
	docker-compose down -v

# 一键构建并推送
.PHONY: release
release: test docker-build docker-push

# 更新线上服务（使用 Helm）
.PHONY: deploy
deploy:
	@if [ -n "$(GET_LATEST_TAG)" ]; then \
		echo "Deploying version $(GET_LATEST_TAG)..."; \
		helm upgrade --install $(IMAGE_NAME) ./chart \
			--set image.repository=$(DOCKER_REGISTRY)/$(IMAGE_NAME) \
			--set image.tag=$(GET_LATEST_TAG) \
			--set appVersion=$(shell echo $(GET_LATEST_TAG) | sed 's/^v//') \
			-n default; \
	else \
		echo "No tags found, skipping deployment"; \
	fi

# 验证 Swagger 文档
.PHONY: swagger-validate
swagger-validate:
	@echo "验证 Swagger 文档..."
	@if [ ! -f ./api/swagger/swagger.json ]; then \
		echo "错误: swagger.json 文件不存在"; \
		exit 1; \
	fi
	@echo "Swagger 文档验证成功"

# 生成 Swagger 验证工具
.PHONY: build-swagger-validator
build-swagger-validator:
	@echo "构建 Swagger 验证工具..."
	@go build -o bin/swagger-validator cmd/swagger-validator/main.go
	@echo "Swagger 验证工具构建完成: bin/swagger-validator"

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
	@echo "  make swagger-validate - 验证 Swagger 文档"
	@echo "  make build-swagger-validator - 构建 Swagger 验证工具"
	@echo "  make helm-install  - 安装 Helm 应用"
	@echo "  make helm-uninstall - 卸载 Helm 应用"
	@echo "  make helm-template - 生成 Helm 模板"
	@echo "  make helm-lint     - 验证 Helm 模板"
	@echo "  make run-grpcui-local - 本地gRPC Web UI调试"
	@echo "  make run-grpcui-remote - 远程gRPC Web UI调试"

# Helm 相关命令
.PHONY: helm-install
helm-install:
	helm upgrade --install $(IMAGE_NAME) ./chart \
		--set image.repository=$(DOCKER_REGISTRY)/$(IMAGE_NAME) \
		--set image.tag=$(GET_LATEST_TAG) \
		-n default

.PHONY: helm-uninstall
helm-uninstall:
	helm uninstall $(IMAGE_NAME) -n default

.PHONY: helm-template
helm-template:
	helm template $(IMAGE_NAME) ./chart \
		--set image.repository=$(DOCKER_REGISTRY)/$(IMAGE_NAME) \
		--set image.tag=$(GET_LATEST_TAG)

.PHONY: helm-lint
helm-lint:
	helm lint ./chart 

# 本地运行服务
.PHONY: local-run
local-run:
	ACTIVE_PROFILE=local go run ./...

# 生成 wire 依赖注入代码
.PHONY: generate
generate:
	go generate ./...

# 获取最新tag并自增版本号
.PHONY: bump-version
bump-version:
	@if [ -z "$$(git tag)" ]; then \
		NEW_TAG="v1.0.0"; \
	else \
		LATEST_TAG=$$(git tag -l | sort -V | tail -n1); \
		if [ "$$(git rev-parse $$LATEST_TAG)" = "$$(git rev-parse HEAD)" ]; then \
			echo "Current commit already tagged as $$LATEST_TAG"; \
			exit 0; \
		fi; \
		MAJOR=$$(echo $$LATEST_TAG | cut -d. -f1 | sed 's/v//'); \
		MINOR=$$(echo $$LATEST_TAG | cut -d. -f2); \
		PATCH=$$(echo $$LATEST_TAG | cut -d. -f3); \
		NEW_PATCH=$$((PATCH + 1)); \
		NEW_TAG="v$$MAJOR.$$MINOR.$$NEW_PATCH"; \
	fi; \
	if git rev-parse "$$NEW_TAG" >/dev/null 2>&1; then \
		echo "Tag $$NEW_TAG already exists, skipping tag creation"; \
		exit 0; \
	else \
		echo "Creating new tag: $$NEW_TAG"; \
		git tag $$NEW_TAG && \
		git push origin $$NEW_TAG; \
	fi

# 更新 chart 版本
.PHONY: update-chart-version
update-chart-version:
	@echo "Updating chart version..."
	@LATEST_TAG=$$(git tag -l | sort -V | tail -n1); \
	if [ -n "$$LATEST_TAG" ]; then \
		sed -i '' "s/appVersion:.*/appVersion: \"$${LATEST_TAG#v}\"/" chart/Chart.yaml; \
		echo "Updated chart appVersion to $${LATEST_TAG#v}"; \
	else \
		echo "No tags found, skipping chart version update"; \
	fi

# 修改 ci target
.PHONY: ci
ci: fmt proto generate test bump-version docker-build docker-push deploy

# gRPC 接口调试
.PHONY: run-grpcui-local run-grpcui-remote
run-grpcui-local:
	@echo "Starting local gRPC Web UI..."
	@grpcui -plaintext localhost:9000
	@echo "gRPC Web UI session ended"

run-grpcui-remote:
	@echo "Starting remote gRPC Web UI..."
	@grpcui sla2-grpc.leeszi.cn:443
	@echo "gRPC Web UI session ended"

# 更新 proto 子模块
.PHONY: update-proto
update-proto:
	git submodule update --remote --merge
	git submodule foreach git checkout main

# 生成代码
.PHONY: docs
docs:
	@echo "生成 protobuf 和 swagger 文档..."
	@cd sla2-proto && buf generate
	@mkdir -p api/swagger
	@cp -r sla2-proto/gen/openapiv2/api.swagger.json api/swagger/swagger.json
	@echo "文档生成完成"

# 发布前检查
.PHONY: fmt
fmt:
	go fmt ./...

changelog:
	@echo "## [$(GET_LATEST_TAG)] - $$(date +%Y-%m-%d)" >> CHANGELOG.md
	@echo "" >> CHANGELOG.md
	@echo "### Added" >> CHANGELOG.md
	@echo "### Changed" >> CHANGELOG.md
	@echo "### Deprecated" >> CHANGELOG.md
	@echo "### Removed" >> CHANGELOG.md
	@echo "### Fixed" >> CHANGELOG.md
	@echo "### Security" >> CHANGELOG.md
	@echo "" >> CHANGELOG.md
