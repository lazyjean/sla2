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

配置文件位于 `config` 目录下，支持多环境配置。配置文件命名规则如下：

- `config.yaml`: 默认配置文件
- `config-{env}.yaml`: 特定环境配置文件，其中 `{env}` 为环境名称

环境配置通过环境变量 `ACTIVE_PROFILE` 指定，例如：

- `ACTIVE_PROFILE=dev` 将加载 `config-dev.yaml`
- `ACTIVE_PROFILE=prod` 将加载 `config-prod.yaml`
- 未设置 `ACTIVE_PROFILE` 时将加载 `config.yaml`

配置文件支持两种加载方式：

1. 从项目内嵌的配置文件加载（使用 `go:embed` 特性）
2. 从本地文件系统加载

此外，如果启用了 Apollo 配置中心，还支持从 Apollo 加载配置：

```yaml
apollo:
  enabled: true # 是否启用 Apollo 配置中心
  app_id: "your-app" # Apollo 应用 ID
  cluster: "default" # Apollo 集群名称
  ip: "localhost:8080" # Apollo 服务地址
  namespace: "application" # Apollo 命名空间
  secret: "" # Apollo 访问密钥
```

### 配置项说明

```yaml
server:
  port: "9000" # 服务端口
  mode: "debug" # 运行模式：debug/release
  version: "v1.0.0" # 服务版本号

grpc:
  port: 9001 # gRPC 服务端口

database:
  host: localhost # 数据库主机
  port: 5432 # 数据库端口
  user: "user" # 数据库用户名
  password: "****" # 数据库密码
  dbname: "sla2" # 数据库名
  sslmode: disable # SSL 模式
  debug: true # 是否开启调试模式
  max_open_conns: 100 # 最大打开连接数
  max_idle_conns: 10 # 最大空闲连接数
  conn_max_lifetime: "30m" # 连接最大生命周期
  conn_max_idle_time: "10m" # 连接最大空闲时间

redis:
  host: localhost # Redis 主机
  port: 6379 # Redis 端口
  password: "" # Redis 密码
  db: 0 # 使用的数据库编号
  max_retries: 3 # 最大重试次数
  min_retry_backoff: "100ms" # 最小重试间隔
  max_retry_backoff: "2s" # 最大重试间隔
  pool_size: 100 # 连接池大小
  min_idle_conns: 10 # 最小空闲连接数
  max_conn_age: "30m" # 连接最大生命周期

log:
  level: "debug" # 日志级别：debug/info/warn/error
  file_path: "./logs/app.log" # 日志文件路径

jwt:
  token_secret_key: "****" # JWT Token 密钥（已隐藏）
  refresh_secret_key: "****" # JWT 刷新密钥（已隐藏）

apple:
  client_id: "your-apple-client-id-here" # Apple 登录 Client ID
```

### 配置加载优先级

配置项的加载遵循以下优先级（从高到低）：

1. 环境变量（使用 `APP_` 前缀）
2. Apollo 配置中心（如果启用）
3. 本地配置文件
   - 先尝试加载内嵌的配置文件
   - 如果内嵌配置不存在，则尝试从文件系统加载
4. 默认值

### 环境变量配置

所有配置项都可以通过环境变量覆盖，环境变量名称规则：

- 使用大写字母
- 使用下划线连接
- 添加 `APP_` 前缀（注意：不是之前文档中提到的 `SLA2_` 前缀）

例如：

```bash
export APP_SERVER_PORT=9000
export APP_DATABASE_HOST=db.example.com
export APP_REDIS_PASSWORD=secret
```

### 配置热重载

服务支持配置热重载，可以在不重启服务的情况下更新配置：

1. 修改配置文件
2. 发送 SIGHUP 信号给服务进程

```bash
kill -SIGHUP $(pgrep sla2)
```

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

## 环境搭建

### microk8s 安装

```bash
sudo snap install microk8s --classic
```

### 组件安装 (cert-manager, ingress, registry)

```bash
microk8s enable cert-manager ingress registry
```

### 配置远程访问集群

todo

### 配置代理，解决镜像拉取问题

编辑 `/var/snap/microk8s/current/args/containerd-env` 文件，添加以下内容：

```bash
HTTP_PROXY=http://your-proxy-server:port
HTTPS_PROXY=http://your-proxy-server:port
NO_PROXY=localhost,127.0.0.1
```

然后重启 MicroK8s：

```bash
sudo systemctl restart snap.microk8s.daemon-containerd
```

### 备案域名

### cert-manager 使用

```bash
microk8s kubectl get cert # 查看证书
microk8s kubectl delete cert <xxx> # 重新生成证书
```
