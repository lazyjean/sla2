# 生词本服务 (Vocabulary Service)

这是一个基于 Gin 框架开发的生词本后端服务，提供单词管理的相关功能。

## 项目结构

```
.
├── api/              # API 相关定义
│   ├── dto/         # 数据传输对象
│   └── swagger/     # Swagger 文档配置
├── config/          # 配置文件目录
│   ├── config.go    # 配置定义
│   └── dev.yaml     # 开发环境配置
├── internal/        # 内部包
│   ├── models/      # 数据模型
│   │   └── word.go  # 单词模型
│   ├── handlers/    # HTTP 处理器
│   ├── repositories/# 数据访问层
│   └── services/    # 业务逻辑层
├── pkg/             # 公共工具包
│   ├── logger/      # 日志工具
│   ├── redis/       # Redis 客户端
│   ├── database/    # 数据库工具
│   └── errors/      # 错误处理
├── scripts/         # 脚本文件
│   ├── migrations/  # 数据库迁移
│   └── deploy/      # 部署脚本
├── test/            # 测试文件
│   ├── integration/ # 集成测试
│   └── mock/        # 测试模拟
├── deployments/     # 部署配置
│   ├── docker/      # Docker 相关
│   └── kubernetes/  # K8s 配置
├── docs/            # 项目文档
├── main.go          # 应用入口
├── Makefile         # 构建脚本
├── go.mod           # Go 模块文件
└── README.md        # 项目说明
```

目录说明：

- `api/`: 存放 API 相关定义，包括 DTO 和 Swagger 配置
- `internal/`: 存放项目内部代码，不对外暴露
- `pkg/`: 可被外部项目引用的公共工具包
- `scripts/`: 各类脚本文件，如数据库迁移和部署脚本
- `test/`: 测试相关文件
- `deployments/`: 部署相关的配置文件
- `docs/`: 项目详细文档

## 技术栈

- Gin: Web 框架
- Redis: 缓存服务
- MySQL: 数据存储
- Zap: 日志记录

## 主要功能

- 单词管理（增删改查）
- 单词本管理
- 学习记录
- 复习提醒

## 配置说明

服务配置包括：

- 服务器配置（端口、运行模式等）
- 数据库配置（MySQL 连接信息）
- Redis 配置（缓存服务器信息）
- 日志配置（日志级别）

## 开发环境设置

1. 安装依赖

```bash
go mod download
```

2. 构建和运行

```bash
# 构建二进制文件
make build

# 本地运行服务
make run

# 停止服务
make stop

# 查看服务日志
make logs
```

## 构建和部署

### Docker 镜像构建

镜像版本号自动从 Git tag 获取，请确保在构建前已创建相应的 Git tag：

```bash
# 创建 Git tag（例如：v1.0.0）
git tag v1.0.0
git push origin v1.0.0

# 本地构建镜像(需要先构建二进制文件)
make docker-build-local

# 远程构建镜像(使用多阶段构建)
make docker-build
```

### 版本管理规范

镜像版本号遵循语义化版本规范 (Semantic Versioning)，通过 Git tag 进行管理：

- 主版本号：不兼容的 API 修改（MAJOR）
- 次版本号：向下兼容的功能性新增（MINOR）
- 修订号：向下兼容的问题修正（PATCH）

示例：v1.2.3

- v1：主版本号
- 2：次版本号
- 3：修订号

### 镜像发布

```bash
# 推送镜像到仓库
make docker-push

# 一键构建并推送
make release
```

### Kubernetes 部署

项目支持使用 Helm 在 Kubernetes 集群中部署:

```bash
# 安装/更新 Helm 应用
make helm-install

# 卸载 Helm 应用
make helm-uninstall

# 生成 Helm 模板
make helm-template

# 验证 Helm 模板
make helm-lint

# 更新已部署的服务
make deploy
```

## API 文档

### Swagger 文档

项目集成了 Swagger 文档，可通过以下方式访问：

```bash
# 本地开发环境
http://localhost:9000/swagger/index.html

# 生产环境
https://sla2.leeszi.cn/swagger/index.html
```

Swagger 文档提供：

- API 接口列表
- 请求/响应参数说明
- 在线接口测试功能
- OpenAPI 规范文档下载

## 部署

项目支持 Docker 部署，使用 K8s 进行容器编排。

## 日志

服务使用结构化日志，以 JSON 格式输出到标准输出（stdout），方便在 K8s 环境中收集和分析。

日志字段说明：

- time: 时间戳
- level: 日志级别
- caller: 调用位置
- msg: 日志消息
- 其他字段: 根据具体日志内容添加

## 开发规范

1. 代码风格遵循 Go 标准规范
2. 提交信息需要清晰描述改动内容
3. 重要功能需要添加单元测试
4. 接口需要添加文档注释
