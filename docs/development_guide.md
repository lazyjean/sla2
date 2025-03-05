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