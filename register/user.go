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
	s.AddTool(
		mcp.NewTool("create_user", mcp.WithDescription("创建Zabbix用户"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithString("username", mcp.Required(), mcp.Description("Zabbix用户名")),
			mcp.WithString("name", mcp.Description("用户真实姓名")),
			mcp.WithString("userGroup", mcp.Required(), mcp.Description("用户组ID")),
			mcp.WithString("roleID", mcp.Description("角色ID")),
		),
		handler.CreateUsersHandler,
	)
	s.AddTool(
		mcp.NewTool("get_groups",
			mcp.WithDescription("获取所有Zabbix用户组信息"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithString("name", mcp.Description("用户组名称")),
			mcp.WithString("status", mcp.Description("用户组状态: 0启用 1禁用 默认: 0")),
			mcp.WithBoolean("selectUsers", mcp.Description("是否获取用户组下用户列表 默认: false")),
			mcp.WithBoolean("selectRights", mcp.Description("是否获取用户组权限列表 默认: false")),
			mcp.WithBoolean("selectTagFilters", mcp.Description("是否获取用户组标签过滤器列表 默认: false")),
		),
		handler.GetUserGroupsHandler,
	)
}
