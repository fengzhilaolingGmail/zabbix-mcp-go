package register

import (
	"zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerGetTemplate(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_templates", mcp.WithDescription("获取所有Zabbix模板信息"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithString("name", mcp.Description("模板名称")),
		),
		handler.GetTemplatesHandler,
	)
}
