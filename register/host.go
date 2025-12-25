package register

import (
	"zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerHost(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_hosts",
			mcp.WithDescription("获取实例Zabbix主机信息,支持所有或模糊匹配"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_for_hostname", mcp.WithDescription("通过主机名称获取主机信息"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostname", mcp.Required(), mcp.Description("主机名称")),
		),
		handler.GetHostsHandler,
	)
}
