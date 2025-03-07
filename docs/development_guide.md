# 开发规范指南

本文档包含了项目开发过程中需要遵循的各项规范和最佳实践。

## 数据库相关规范

### PostgreSQL 数组类型处理规范

在使用 GORM 操作 PostgreSQL 数据库时，对于数组类型字段的处理有以下规则：

1. 对于简单数组类型（如 `text[]`, `varchar[]` 等），**必须**使用 `github.com/lib/pq` 提供的数组类型：
   ```go
   import pg "github.com/lib/pq"

   type Example struct {
       Names pg.StringArray `gorm:"type:text[]"`
       IDs   pg.Int64Array `gorm:"type:integer[]"`
   }
   ```

2. 对于需要存储更复杂的数据结构或需要 JSON 查询功能时，使用 `jsonb` 类型：
   ```go
   type Example struct {
       Tags []string `gorm:"type:jsonb;serializer:json"`
   }
   ```

3. **严禁**直接使用 Go 原生数组类型来映射 PostgreSQL 数组：
   ```go
   // ❌ 错误示例 - 会导致写入失败
   type Wrong struct {
       Names []string `gorm:"type:text[]"`  // 错误
   }

   // ✅ 正确示例
   type Correct struct {
       Names pg.StringArray `gorm:"type:text[]"`  // 正确
   }
   ```

原因说明：
- PostgreSQL 的数组类型需要特殊的序列化/反序列化处理
- 直接使用 Go 原生数组类型会导致数据写入失败
- `github.com/lib/pq` 包提供了必要的类型转换功能

选择建议：
1. 如果只需要简单的数组存储，优先使用 `pg.Array` 类型
2. 如果需要复杂的查询或数据结构，考虑使用 `jsonb` 类型 

## 代码组织与架构规范

### 日志处理规范

在处理日志时，必须遵循以下规则：

1. **从 Context 中获取 Logger**：不应在结构体中保存 logger 实例，而应从 context 中获取。

   ```go
   // ❌ 错误示例 - 在结构体中保存logger
   type QuestionTagService struct {
       pb.UnimplementedQuestionTagServiceServer
       tagRepo repository.QuestionTagRepository
       log     *zap.Logger  // 错误：不应在结构体中保存logger
   }
   
   func NewQuestionTagService(tagRepo repository.QuestionTagRepository) *QuestionTagService {
       return &QuestionTagService{
           tagRepo: tagRepo,
           log:     logger.GetLogger(context.Background()),  // 错误：创建时初始化logger
       }
   }
   
   func (s *QuestionTagService) ListTag(ctx context.Context, req *pb.QuestionTagServiceListTagRequest) (*pb.QuestionTagServiceListTagResponse, error) {
       s.log.Info("ListTag called")  // 错误：使用结构体中的logger
       // ...
   }
   ```

   ```go
   // ✅ 正确示例 - 从context中获取logger
   type QuestionTagService struct {
       pb.UnimplementedQuestionTagServiceServer
       tagRepo repository.QuestionTagRepository
   }
   
   func NewQuestionTagService(tagRepo repository.QuestionTagRepository) *QuestionTagService {
       return &QuestionTagService{
           tagRepo: tagRepo,
       }
   }
   
   func (s *QuestionTagService) ListTag(ctx context.Context, req *pb.QuestionTagServiceListTagRequest) (*pb.QuestionTagServiceListTagResponse, error) {
       log := logger.GetLogger(ctx)  // 正确：从context中获取logger
       log.Info("ListTag called", zap.Any("req", req))
       // ...
   }
   ```

原因说明：
- 保持日志上下文的一致性（请求ID、用户ID等）
- 避免在每个结构体中重复保存logger实例
- 便于在整个调用链中传递日志上下文

### 分层架构规范

项目采用领域驱动设计（DDD）的分层架构，接口层与其它层交互时必须遵循以下规则：

1. **接口层只能调用应用服务层**：接口层（如gRPC服务实现）不能直接调用领域层或基础设施层。

   ```go
   // ❌ 错误示例 - 接口层直接调用领域层
   func (s *QuestionService) GetQuestion(ctx context.Context, req *pb.QuestionServiceGetRequest) (*pb.QuestionServiceGetResponse, error) {
       log := logger.GetLogger(ctx)
       
       // 错误：直接使用仓储接口
       questionRepo := s.questionRepo
       id := strconv.FormatUint(uint64(req.GetIds()[0]), 10)
       
       // 错误：直接调用领域层
       question, err := questionRepo.Get(ctx, id)
       if err != nil {
           log.Error("GetQuestion failed", zap.Error(err))
           return nil, status.Error(codes.Internal, "failed to get question")
       }
       
       return &pb.QuestionServiceGetResponse{
           Questions: []*pb.Question{question.ToProto()},
       }, nil
   }
   ```

   ```go
   // ✅ 正确示例 - 接口层只调用应用服务层
   func (s *QuestionService) GetQuestion(ctx context.Context, req *pb.QuestionServiceGetRequest) (*pb.QuestionServiceGetResponse, error) {
       log := logger.GetLogger(ctx)
       log.Info("GetQuestion called", zap.Any("req", req))
   
       // 正确：调用应用服务层
       id := strconv.FormatUint(uint64(req.GetIds()[0]), 10)
       question, err := s.questionService.Get(ctx, id)
       if err != nil {
           log.Error("GetQuestion failed", zap.Error(err))
           return nil, status.Error(codes.Internal, "failed to get question")
       }
   
       return &pb.QuestionServiceGetResponse{
           Questions: []*pb.Question{question.ToProto()},
       }, nil
   }
   ```

原因说明：
- 保持清晰的分层架构
- 确保业务逻辑集中在应用服务层
- 避免接口层与领域实现细节的耦合
- 便于测试和维护

分层职责说明：
- **接口层**：处理请求解析、参数验证、错误转换和响应格式化
- **应用服务层**：编排领域对象，实现用例逻辑，管理事务
- **领域层**：包含业务规则和核心逻辑
- **基础设施层**：提供技术实现，如数据库访问、外部服务集成等

## 单元测试规范

### 应用层单元测试规范

应用层单元测试旨在验证应用层服务的业务逻辑正确性，确保它们能够正确处理输入和输出、调用领域服务和基础设施服务，以及正确处理各种边界条件和错误场景。

#### 测试结构规范

1. **文件命名与位置**
   - 测试文件应放置在与被测试服务相同的包中
   - 测试文件命名应遵循 `[服务名]_test.go` 的格式，例如 `admin_service_test.go`

2. **测试函数命名**
   - 测试函数应采用以下命名约定：`Test[服务名]_[方法名]`
   - 例如：`TestAdminService_Login`、`TestUserService_Register`

3. **测试用例分组**
   - 使用 `t.Run()` 组织同一方法的多个测试场景：

   ```go
   func TestUserService_Register(t *testing.T) {
       t.Run("注册成功", func(t *testing.T) {
           // 测试注册成功的场景
       })
       
       t.Run("用户名已存在", func(t *testing.T) {
           // 测试用户名已存在的场景
       })
   }
   ```

#### Mock 对象规范

1. **Mock 对象定义**
   - 为依赖的接口创建 Mock 实现，使用 `testify/mock` 包：

   ```go
   type MockUserRepository struct {
       mock.Mock
   }

   func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
       args := m.Called(ctx, username)
       if args.Get(0) == nil {
           return nil, args.Error(1)
       }
       return args.Get(0).(*entity.User), args.Error(1)
   }
   ```

2. **Mock 对象命名**
   - Mock 对象应以 `Mock` 前缀命名，例如：`MockUserRepository`、`MockPasswordService`

#### 测试用例编写规范

1. **测试用例结构**
   - 每个测试用例应包含四部分：准备(Arrange)、执行(Act)、断言(Assert)、验证(Verify)

   ```go
   // 准备测试数据
   ctx := context.Background()
   req := &dto.UserRegisterRequest{
       Username: "testuser",
       Password: "password123",
       Email:    "test@example.com",
   }

   // 创建并配置 Mock 对象
   mockRepo := new(MockUserRepository)
   mockRepo.On("FindByUsername", ctx, req.Username).Return(nil, errors.New("not found"))
   mockRepo.On("Create", ctx, mock.MatchedBy(func(u *entity.User) bool {
       return u.Username == req.Username
   })).Return(nil)

   // 创建服务实例
   service := NewUserService(mockRepo, mockPasswordService, mockTokenService)

   // 执行测试
   resp, err := service.Register(ctx, req)

   // 断言结果
   assert.NoError(t, err)
   assert.NotNil(t, resp)
   assert.Equal(t, req.Username, resp.User.Username)

   // 验证 Mock 是否按预期被调用
   mockRepo.AssertExpectations(t)
   ```

2. **测试数据准备**
   - 优先使用显式声明的测试数据，避免依赖外部资源
   - 对于复杂对象，使用工厂函数或测试辅助函数创建

   ```go
   // 创建测试用的管理员实体
   admin := &entity.Admin{
       ID:        entity.UID(1),
       Username:  "admin",
       Password:  "hashed_password",
       Nickname:  "Admin",
       Roles:     []string{"admin"},
       CreatedAt: time.Now(),
       UpdatedAt: time.Now(),
   }
   ```

3. **Mock 配置规范**
   - 为每个预期的方法调用配置 Mock 对象
   - 使用 `mock.MatchedBy` 验证复杂参数
   - 确保 Mock 的返回值符合真实场景

   ```go
   // 配置 Mock 对象的行为
   mockTokenService.On("ValidateTokenFromContext", ctx).Return(admin.ID, admin.Roles, nil)
   mockRepo.On("FindByID", ctx, admin.ID).Return(admin, nil)
   ```

4. **断言规范**
   - 使用 `testify/assert` 或 `testify/require` 包进行断言
   - 对关键字段进行显式断言，不仅仅断言整个对象
   - 正确处理错误场景的断言

   ```go
   // 结果断言
   assert.NoError(t, err)
   assert.NotNil(t, resp)
   assert.Equal(t, admin.ID, resp.ID)
   assert.Equal(t, admin.Username, resp.Username)
   ```

5. **Mock 验证规范**
   - 完成测试后，验证所有 Mock 对象是否按预期被调用：

   ```go
   // 验证 Mock 对象是否按预期被调用
   mockRepo.AssertExpectations(t)
   mockTokenService.AssertExpectations(t)
   ```

#### 错误处理测试规范

1. **错误场景测试**
   - 为每个方法测试可能的错误场景，包括：输入验证错误、数据库操作错误、依赖服务错误等

   ```go
   t.Run("用户名已存在", func(t *testing.T) {
       // 准备测试数据
       existingUser := &entity.User{
           ID:       entity.UID(1),
           Username: "existinguser",
       }
       
       // 配置 Mock 对象
       mockRepo.On("FindByUsername", ctx, req.Username).Return(existingUser, nil)
       
       // 执行测试
       resp, err := service.Register(ctx, req)
       
       // 断言结果
       assert.Error(t, err)
       assert.Nil(t, resp)
       assert.Equal(t, "username already exists", err.Error())
       
       // 验证 Mock 对象
       mockRepo.AssertExpectations(t)
   })
   ```

2. **边界条件测试**
   - 测试空值、零值、最大值、最小值、特殊字符等边界条件

#### 集成测试与单元测试的区分

1. **单元测试范围**
   - 单元测试应该只测试单一组件（应用服务）的功能
   - Mock 所有外部依赖
   - 不访问数据库、文件系统、网络等
   - 避免依赖环境变量或配置

2. **集成测试**
   - 集成测试应单独编写，放在 `tests` 目录下
   - 测试多个组件的交互
   - 使用真实的依赖（如数据库、缓存等）
   - 测试端到端功能

#### 测试覆盖率要求

1. **覆盖率目标**
   - 应用服务层代码的测试覆盖率目标：
     - 行覆盖率：≥ 80%
     - 分支覆盖率：≥ 75%
     - 方法覆盖率：100%（所有公开方法必须有测试）

2. **覆盖率检查**
   - 使用以下命令检查测试覆盖率：

   ```bash
   go test -coverprofile=coverage.out ./internal/application/service/...
   go tool cover -html=coverage.out -o coverage.html
   ```

#### 完整的应用服务测试示例

```go
package service

import (
	"context"
	"testing"
	"time"
    "errors"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository 用户仓库的 Mock 实现
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id entity.UID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

// MockPasswordService 密码服务的 Mock 实现
type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) VerifyPassword(password, hashedPassword string) bool {
	args := m.Called(password, hashedPassword)
	return args.Bool(0)
}

// MockTokenService 令牌服务的 Mock 实现
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(userID entity.UID, roles []string) (string, error) {
	args := m.Called(userID, roles)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateTokenFromContext(ctx context.Context) (entity.UID, []string, error) {
	args := m.Called(ctx)
	return args.Get(0).(entity.UID), args.Get(1).([]string), args.Error(2)
}

func TestUserService_Login(t *testing.T) {
    // 准备基本测试数据
    ctx := context.Background()
    
    t.Run("登录成功", func(t *testing.T) {
        // 准备测试数据
        req := &dto.UserLoginRequest{
            Username: "testuser",
            Password: "password123",
        }
        
        user := &entity.User{
            ID:        entity.UID(1),
            Username:  "testuser",
            Password:  "hashed_password",
            Nickname:  "Test User",
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
        
        // 创建 Mock 对象
        mockRepo := new(MockUserRepository)
        mockPasswordService := new(MockPasswordService)
        mockTokenService := new(MockTokenService)
        
        // 配置 Mock 行为
        mockRepo.On("FindByUsername", ctx, req.Username).Return(user, nil)
        mockPasswordService.On("VerifyPassword", req.Password, user.Password).Return(true)
        mockTokenService.On("GenerateToken", user.ID, []string{"user"}).Return("access_token", nil)
        mockTokenService.On("GenerateRefreshToken", user.ID, []string{"user"}).Return("refresh_token", nil)
        
        // 创建服务实例
        service := NewUserService(mockRepo, mockTokenService, mockPasswordService, nil)
        
        // 执行测试
        resp, err := service.Login(ctx, req)
        
        // 断言结果
        assert.NoError(t, err)
        assert.NotNil(t, resp)
        assert.Equal(t, user.ID, resp.User.ID)
        assert.Equal(t, user.Username, resp.User.Username)
        assert.Equal(t, "access_token", resp.AccessToken)
        assert.Equal(t, "refresh_token", resp.RefreshToken)
        
        // 验证 Mock 对象
        mockRepo.AssertExpectations(t)
        mockPasswordService.AssertExpectations(t)
        mockTokenService.AssertExpectations(t)
    })
    
    t.Run("用户不存在", func(t *testing.T) {
        // 准备测试数据
        req := &dto.UserLoginRequest{
            Username: "nonexistent",
            Password: "password123",
        }
        
        // 创建 Mock 对象
        mockRepo := new(MockUserRepository)
        mockPasswordService := new(MockPasswordService)
        mockTokenService := new(MockTokenService)
        
        // 配置 Mock 行为
        mockRepo.On("FindByUsername", ctx, req.Username).Return(nil, errors.New("user not found"))
        
        // 创建服务实例
        service := NewUserService(mockRepo, mockTokenService, mockPasswordService, nil)
        
        // 执行测试
        resp, err := service.Login(ctx, req)
        
        // 断言结果
        assert.Error(t, err)
        assert.Nil(t, resp)
        assert.Contains(t, err.Error(), "user not found")
        
        // 验证 Mock 对象
        mockRepo.AssertExpectations(t)
    })
}
``` 