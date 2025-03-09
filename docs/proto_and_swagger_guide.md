# Protocol Buffers 和 Swagger 文档指南

## 目录结构

```
sla2/
├── sla2-proto/              # Protocol Buffers 定义目录
│   ├── proto/              # proto 文件目录
│   │   └── v1/            # API 版本目录
│   │       ├── ai_chat.proto
│   │       ├── admin.proto
│   │       └── ...
│   ├── buf.yaml           # buf 配置文件
│   └── buf.gen.yaml       # buf 生成器配置文件
└── api/                   # 生成的 API 文档目录
    └── swagger/           # Swagger 文档目录
        └── swagger.json   # 生成的 Swagger 文档
```

## 定义新服务的流程

### 1. 创建 Proto 文件

在 `sla2-proto/proto/v1/` 目录下创建新的 `.proto` 文件，例如 `example_service.proto`：

```protobuf
syntax = "proto3";

package proto.v1;

option go_package = "github.com/lazyjean/sla2/api/proto/v1;pb";

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";
import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

// 定义服务
service ExampleService {
  // 定义 API 方法
  rpc CreateExample(CreateExampleRequest) returns (ExampleResponse) {
    // HTTP 路由配置
    option (google.api.http) = {
      post: "/v1/examples"
      body: "*"
    };
    // Swagger 文档配置
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      summary: "Create Example";
      description: "Create a new example";
      tags: ["ExampleService"];  // 使用服务名作为标签
      security: {
        security_requirement: {
          key: "Bearer";
          value: {};
        }
      };
      responses: {
        key: "200";
        value: {
          description: "Successfully created";
          schema: {
            json_schema: {
              ref: "#/definitions/proto.v1.ExampleResponse";
            }
          }
        }
      };
    };
  }
}

// 定义消息
message CreateExampleRequest {
  // Swagger 模式定义
  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: {
      title: "Create Example Request";
      description: "Request parameters for creating an example";
      required: ["name"];  // 必填字段
    }
  };
  
  string name = 1 [(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
    description: "Example name";
    example: '{"value": "example name"}';  // 示例值
  }];
}
```

### 2. 生成代码和文档

在 `sla2-proto` 目录下运行：

```bash
buf generate
```

这将生成：
- Go 代码
- gRPC Gateway 代码
- Swagger 文档

### 3. 实现服务

在相应的服务目录下实现 gRPC 服务接口。

## Swagger 文档设计

### 1. 全局配置

在 proto 文件中添加全局 Swagger 配置：

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    title: "SLA2 API";
    version: "1.0";
    description: "API Documentation";
    contact: {
      name: "Your Name";
      email: "your.email@example.com";
    };
  };
  host: "localhost:9102";
  base_path: "/v1";
  schemes: HTTP;
  schemes: HTTPS;
  security_definitions: {
    security: {
      key: "Bearer";
      value: {
        type: TYPE_API_KEY;
        in: IN_HEADER;
        name: "Authorization";
        description: "Bearer token for authentication";
      }
    };
  };
};
```

### 2. API 文档注解

#### 2.1 方法注解

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
  summary: "简短描述";
  description: "详细描述";
  tags: ["服务名"];  // 使用服务名作为标签
  security: { ... }  // 安全配置
  responses: { ... }  // 响应定义
};
```

#### 2.2 消息注解

```protobuf
message ExampleRequest {
  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_schema) = {
    json_schema: {
      title: "标题";
      description: "描述";
      required: ["必填字段"];
    }
  };
}
```

#### 2.3 字段注解

```protobuf
string field = 1 [(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
  description: "字段描述";
  example: '{"value": "示例值"}';  // 注意使用 JSON 格式
}];
```

## 注意事项

### 1. 命名规范

- 服务名使用 PascalCase，以 "Service" 结尾（例如：UserService）
- 方法名使用 PascalCase（例如：CreateUser）
- 消息名使用 PascalCase（例如：CreateUserRequest）
- 字段名使用 snake_case（例如：user_id）

### 2. API 路径规范

- 使用 RESTful 风格
- 路径使用小写，单词用下划线分隔
- 版本号放在路径开头（例如：/v1/users）

### 3. Swagger 文档注意事项

- 标签（tags）使用服务名，保持一致性
- 必须提供清晰的描述和示例
- 示例值必须使用 JSON 格式：`'{"value": "example"}'`
- 响应定义要包含成功和失败的情况
- 安全要求要明确指定

### 4. 生成和部署

- 修改 proto 文件后必须重新生成代码和文档
- 确保生成的 swagger.json 文件被正确部署
- 检查 Swagger UI 是否能正确访问文档

### 5. 常见问题

1. Swagger UI 无法显示文档
   - 检查 swagger.json 文件路径是否正确
   - 确认 Content-Type 是否为 application/json
   - 验证 swagger.json 文件格式是否正确

2. API 分组重复
   - 确保使用统一的 tags 名称
   - 推荐使用服务名作为 tag

3. 示例值格式错误
   - 必须使用 JSON 格式：`'{"value": "example"}'`
   - 不要直接使用字符串：`"example"`（这会导致生成失败）

## 相关命令

```bash
# 编译 proto 文件（主项目）
make proto

# 更新 grpc-gateway 生成的 swagger 文档
make docs

# 手动生成代码和文档（如果不使用 Makefile）
cd sla2-proto
buf generate

# 验证 proto 文件
buf lint

# 启动服务
cd ..
ACTIVE_PROFILE=local go run ./...

# 访问 Swagger UI
http://localhost:9102/swagger/
```

## 编译流程说明

### 1. 编译 Proto 文件

使用 `make proto` 命令会：
1. 进入 sla2-proto 目录
2. 执行 buf generate 生成 Go 代码
3. 生成 gRPC 和 gRPC-Gateway 相关代码
4. 更新 protobuf 生成的 Go 代码

### 2. 更新 Swagger 文档

使用 `make docs` 命令会：
1. 重新生成 swagger.json 文件
2. 更新 API 文档
3. 确保 Swagger UI 能够正确显示最新的 API 文档

### 3. 开发流程

1. 修改 proto 文件后：
   ```bash
   make proto  # 更新 Go 代码
   make docs   # 更新 API 文档
   ```

2. 验证更改：
   - 检查生成的 Go 代码是否正确
   - 访问 Swagger UI 确认文档更新
   - 测试新的 API 端点

3. 常见问题处理：
   - 如果 `make proto` 失败，检查 proto 文件语法
   - 如果 `make docs` 失败，检查 Swagger 注解格式
   - 如果文档未更新，尝试重启服务 