package embed

import (
	"embed"
)

//go:embed model.conf
var FS embed.FS

// GetRBACModelBytes 获取RBAC模型配置文件的字节数据
func GetRBACModelBytes() ([]byte, error) {
	return FS.ReadFile("model.conf")
}
