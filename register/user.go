package register

import (
	"zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerClientPool 注册与 Zabbix 客户端池相关的 MCP 工具。
// 目前仅注册工具元信息；具体处理器在 handler 包未实现时暂留为 nil，以免影响构建。
func registerUser(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_users",
			mcp.WithDescription("获取所有Zabbix用户信息"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithString("username", mcp.Description("Zabbix用户名,留空表示获取所有用户")),
		),
		handler.GetUsersHandler,
	)
}
