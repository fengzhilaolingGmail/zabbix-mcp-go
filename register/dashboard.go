/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 19:05:52
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-06 20:55:33
 * @FilePath: \zabbix-mcp-go\register\dashboard.go
 * @Description: 仪表盘注册
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package register

import (
	"zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerDashboard(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("create_dashboard",
			mcp.WithDescription("在指定实例中创建仪表盘"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix 实例名称")),
			mcp.WithString("name", mcp.Required(), mcp.Description("仪表盘名称")),
			mcp.WithArray("pages", mcp.Description("仪表盘页面数组，包含 widgets 配置。示例: [{\"name\":\"Default Page\",\"display_period\":0,\"widgets\":[{\"type\":\"graph\",\"height\":10,\"width\":36,\"x\":0,\"y\":0,\"fields\":[{\"name\":\"graphid\",\"value\":5658}]}] } ]")),
			mcp.WithArray("users", mcp.Description("仪表盘用户共享数组，包含 userid 和 permission")),
			mcp.WithArray("userGroups", mcp.Description("仪表盘用户组共享数组，包含 usrgrpid 和 permission")),
		),
		handler.CreateDashboardHandler,
	)

	s.AddTool(
		mcp.NewTool("create_graph_dashboard",
			mcp.WithDescription("使用 hosts(主机名列表) + graphids(图形ID列表) 自动创建布局良好的仪表盘（隐藏底层 JSON）"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix 实例名称")),
			mcp.WithString("name", mcp.Required(), mcp.Description("仪表盘名称")),
			mcp.WithArray("hosts", mcp.Required(), mcp.Description("主机列表（字符串）")),
			mcp.WithArray("graphids", mcp.Required(), mcp.Description("图形ID列表（数值或字符串）")),
			mcp.WithNumber("rows", mcp.Description("期望的行数（可选），将据此计算每行 widget 个数")),
			mcp.WithNumber("widgetWidth", mcp.Description("每个 widget 的宽度（可选，默认36）")),
			mcp.WithNumber("widgetHeight", mcp.Description("每个 widget 的高度（可选，默认5）")),
		),
		handler.CreateGraphDashboardHandler,
	)
}
