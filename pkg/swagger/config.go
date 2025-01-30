package swagger

import (
	"os"

	"github.com/lazyjean/sla2/docs"
)

// InitSwagger 初始化 Swagger 配置
func InitSwagger() {
	// 动态设置 Host
	docs.SwaggerInfo.Host = getSwaggerHost()
	// 动态设置 Schemes
	docs.SwaggerInfo.Schemes = getSwaggerSchemes()
}

// getSwaggerHost 根据环境返回适当的 host
func getSwaggerHost() string {
	// 优先使用环境变量
	if host := os.Getenv("SWAGGER_HOST"); host != "" {
		return host
	}

	if env := os.Getenv("GIN_MODE"); env != "release" {
		return "localhost:9000"
	}
	return "sla2.leeszi.cn"
}

// getSwaggerSchemes 根据环境返回适当的 schemes
func getSwaggerSchemes() []string {
	// 优先使用环境变量
	if scheme := os.Getenv("SWAGGER_SCHEME"); scheme != "" {
		return []string{scheme}
	}

	if env := os.Getenv("GIN_MODE"); env != "release" {
		return []string{"http"}
	}
	return []string{"https"}
}
