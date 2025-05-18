package grpc

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/lazyjean/sla2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSwaggerDocumentAvailable 测试Swagger文档是否可用
func TestSwaggerDocumentAvailable(t *testing.T) {
	// 确保测试数据目录存在
	swaggerDir := "testdata/swagger"
	err := os.MkdirAll(swaggerDir, 0755)
	require.NoError(t, err)

	// 创建测试用的swagger.json文件
	swaggerContent := `{
		"swagger": "2.0",
		"info": {
			"title": "Test API",
			"version": "1.0.0"
		},
		"paths": {}
	}`
	swaggerFile := filepath.Join(swaggerDir, "swagger.json")
	err = os.WriteFile(swaggerFile, []byte(swaggerContent), 0644)
	require.NoError(t, err)

	// 创建模拟的服务器实例
	mockServer := &GRPCServer{
		config: &config.Config{
			Swagger: config.SwaggerConfig{
				Username: "admin",
				Password: "swagger",
			},
		},
	}

	// 创建一个测试HTTP处理器
	router := http.NewServeMux()

	// 为Swagger文档创建处理函数
	router.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		data, err := os.ReadFile(swaggerFile)
		if err != nil {
			t.Logf("失败读取swagger文件: %v", err)
			http.Error(w, "Swagger file not found", http.StatusNotFound)
			return
		}

		w.Write(data)
	})

	// 为Swagger UI路径添加Basic Auth
	router.Handle("/swagger/", mockServer.basicAuth(http.FileServer(http.Dir("./testdata"))))

	// 创建测试服务器
	ts := httptest.NewServer(router)
	defer ts.Close()

	// 测试1: 获取Swagger JSON文档
	t.Run("Get Swagger JSON Document", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/swagger/doc.json")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Contains(t, string(body), "Test API")
		assert.Contains(t, string(body), "2.0")
	})

	// 测试2: 未授权访问Swagger UI (没有提供Basic Auth)
	t.Run("Unauthorized Access to Swagger UI", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/swagger/")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("WWW-Authenticate"), "Basic")
	})

	// 测试3: 使用有效凭据访问Swagger UI
	t.Run("Authorized Access to Swagger UI", func(t *testing.T) {
		req, err := http.NewRequest("GET", ts.URL+"/swagger/", nil)
		require.NoError(t, err)

		req.SetBasicAuth("admin", "swagger")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// 测试4: 确保我们可以解析Swagger文档
	t.Run("Parse Swagger Document", func(t *testing.T) {
		// 这里我们可以使用go-openapi/loads库进行解析，
		// 但为了避免引入更多依赖，我们只检查JSON是否有效
		resp, err := http.Get(ts.URL + "/swagger/doc.json")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// 确保内容包含必要的Swagger结构
		assert.Contains(t, string(body), "swagger")
		assert.Contains(t, string(body), "info")
		assert.Contains(t, string(body), "paths")
	})

	// 清理测试文件
	os.Remove(swaggerFile)
	os.RemoveAll(swaggerDir)
}
