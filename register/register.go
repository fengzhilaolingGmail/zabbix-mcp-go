package register

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
)

func Registers(s *server.MCPServer) {
	// 使用handler包中的注册函数
	// handler.RegisterTools(s)
	fmt.Println("Register Tools Success")
	// 注册 ClientPool 相关工具
	registerClientPool(s)
}
