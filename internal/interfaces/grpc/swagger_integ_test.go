//go:build integration
// +build integration

package grpc

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/lazyjean/sla2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSwaggerIntegration 测试实际环境中Swagger的集成
// 要运行此测试: go test -tags=integration ./internal/interfaces/grpc/
func TestSwaggerIntegration(t *testing.T) {
	// 确保测试环境有配置文件
	err := config.InitConfig()
	if err != nil {
		t.Skip("无法加载配置，跳过集成测试:", err)
	}
	cfg := config.GetConfig()

	// 首先尝试启动服务器
	server := createTestServer(t)
	if server == nil {
		t.Skip("无法创建测试服务器，跳过集成测试")
		return
	}

	// 启动服务器
	err = server.Start()
	require.NoError(t, err)

	// 确保在测试结束时关闭服务器
	defer server.Stop()

	// 给服务器一些启动时间
	time.Sleep(2 * time.Second)

	// 构建Swagger文档URL
	swaggerURL := "http://localhost:" + cfg.GRPC.GatewayPort + "/swagger/doc.json"

	// 测试1: 获取Swagger文档不需要认证
	t.Run("Get Swagger JSON Without Auth", func(t *testing.T) {
		resp, err := http.Get(swaggerURL)
		if err != nil {
			t.Fatalf("无法连接到Swagger文档: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	})

	// 测试2: 解析Swagger文档，确保它是有效的JSON
	t.Run("Parse Swagger Document", func(t *testing.T) {
		resp, err := http.Get(swaggerURL)
		if err != nil {
			t.Fatalf("无法连接到Swagger文档: %v", err)
		}
		defer resp.Body.Close()

		// 尝试解析为JSON
		var swaggerDoc map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&swaggerDoc)
		require.NoError(t, err, "Swagger文档不是有效的JSON")

		// 基本验证
		require.Contains(t, swaggerDoc, "swagger", "文档不包含swagger版本字段")
		require.Contains(t, swaggerDoc, "info", "文档不包含info字段")
		require.Contains(t, swaggerDoc, "paths", "文档不包含paths字段")

		// 验证版本
		assert.Equal(t, "2.0", swaggerDoc["swagger"], "不是Swagger 2.0文档")

		// 验证paths是否包含API端点
		paths, ok := swaggerDoc["paths"].(map[string]interface{})
		require.True(t, ok, "paths不是对象类型")
		assert.NotEmpty(t, paths, "API路径为空")
	})

	// 测试3: 访问Swagger UI需要认证
	t.Run("Swagger UI Requires Auth", func(t *testing.T) {
		swaggerUIURL := "http://localhost:" + cfg.GRPC.GatewayPort + "/swagger/"
		resp, err := http.Get(swaggerUIURL)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("WWW-Authenticate"), "Basic")
	})

	// 测试4: 使用正确的凭据访问Swagger UI
	t.Run("Access Swagger UI With Auth", func(t *testing.T) {
		swaggerUIURL := "http://localhost:" + cfg.GRPC.GatewayPort + "/swagger/"
		req, err := http.NewRequest("GET", swaggerUIURL, nil)
		require.NoError(t, err)

		req.SetBasicAuth(cfg.Swagger.Username, cfg.Swagger.Password)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// createTestServer 创建用于测试的服务器实例
func createTestServer(t *testing.T) *Server {
	// 这里需要根据实际情况构建服务器
	// 此示例仅为最简化版本
	cfg := config.GetConfig()

	// 如果需要完整的测试服务器，您需要提供所有依赖项
	// 在这里，我们只是返回一个空壳，仅供集成测试使用
	// 实际应用中，您应该考虑使用依赖注入或工厂方法

	return &Server{
		config: cfg,
		// 其他依赖项...
	}
}
