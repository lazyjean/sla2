# 生词本服务 (Vocabulary Service)

这是一个基于 Gin 框架开发的生词本后端服务，提供单词管理的相关功能。

## 项目结构

```
.
├── config/             # 配置文件目录
│   └── config.go      # 配置定义
├── models/            # 数据模型目录
│   └── word.go       # 单词模型定义
├── handlers/          # HTTP 处理器目录
│   └── word_handler.go
├── repositories/      # 数据仓库层目录
│   └── word_repository.go
├── routes/           # 路由配置目录
│   └── word_routes.go
├── services/         # 业务逻辑层目录
│   └── word_service.go
├── pkg/              # 公共包目录
│   ├── logger/       # 日志工具
│   └── redis/        # Redis 客户端
├── main.go           # 应用入口
├── go.mod            # Go 模块文件
└── README.md         # 项目说明文档
```

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

2. 运行服务

```bash
go run main.go
```

## API 文档

### 基础接口

- `GET /ping`: 健康检查接口

### 单词接口

- `POST /api/words`: 创建新单词
- `GET /api/words`: 获取单词列表
- (更多接口文档待补充)

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

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交改动
4. 创建 Pull Request

## 许可证

[MIT License](LICENSE)
