package banner

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const banner = `
███████╗██╗      █████╗ ██████╗
██╔════╝██║     ██╔══██╗╚════██╗
███████╗██║     ███████║ █████╔╝
╚════██║██║     ██╔══██║██╔═══╝ 
███████║███████╗██║  ██║███████╗
╚══════╝╚══════╝╚═╝  ╚═╝╚══════╝
`

// PrintBanner 打印启动 banner
func PrintBanner(version string) {
	// 获取终端宽度
	width := 80
	if w, _, err := term.GetSize(0); err == nil {
		width = w
	}

	// 计算 banner 的宽度
	bannerLines := strings.Split(banner, "\n")
	maxWidth := 0
	for _, line := range bannerLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// 计算居中所需的空格
	padding := (width - maxWidth) / 2

	// 打印 banner
	fmt.Println("\033[36m") // 设置颜色为青色
	for _, line := range bannerLines {
		if line != "" {
			fmt.Printf("%s%s\n", strings.Repeat(" ", padding), line)
		}
	}
	fmt.Println("\033[0m") // 重置颜色

	// 获取当前 profile
	profile := os.Getenv("ACTIVE_PROFILE")
	if profile == "" {
		profile = "local"
	}

	// 打印 profile 信息
	infoLine := fmt.Sprintf("Profile: %s", profile)
	infoPadding := (width - len(infoLine)) / 2
	fmt.Printf("%s%s\n", strings.Repeat(" ", infoPadding), infoLine)

	// 打印分隔线
	fmt.Println(strings.Repeat("=", width))
}
