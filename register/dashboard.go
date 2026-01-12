/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 19:05:52
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-09 15:40:05
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
		mcp.NewTool("get_dashboard",
			mcp.WithDescription("获取指定实例中的仪表盘信息"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix 实例名称")),
			mcp.WithString("dashboard_name", mcp.Required(), mcp.Description("仪表盘名称")),
			mcp.WithBoolean("select_users", mcp.Description("查看仪表盘共享用户列表")),
			mcp.WithBoolean("select_usergroups", mcp.Description("查看仪表盘共享用户组列表")),
			mcp.WithBoolean("select_pages", mcp.Description("查看仪表盘页面列表")),
		),
		handler.GetDashboardHandler,
	)

	s.AddTool(
		mcp.NewTool("create_dashboard",
			mcp.WithDescription("在指定实例中创建仪表盘"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix 实例名称")),
			mcp.WithString("name", mcp.Required(), mcp.Description("仪表盘名称")),
			mcp.WithArray("pages", mcp.Description("仪表盘页面数组，包含 widgets 配置。示例: [{\"name\":\"Default Page\",\"display_period\":0,\"widgets\":[{\"type\":\"graph\",\"height\":10,\"width\":36,\"x\":0,\"y\":0,\"fields\":[{\"name\":\"graphid\",\"value\":5658}]}] } ]")),
			mcp.WithString("userid", mcp.Description("仪表板所有者用户ID (可选)")),
			mcp.WithNumber("private", mcp.Description("仪表盘共享类型 (0 公共, 1 私有， 可选)")),
			mcp.WithNumber("display_period", mcp.Description("默认页面显示周期（秒，可选）")),
			mcp.WithNumber("auto_start", mcp.Description("是否自动启动幻灯片播放 (0/1 可选)")),
			mcp.WithArray("users", mcp.Description("仪表盘用户共享数组，包含 userid 和 permission")),
			mcp.WithArray("userGroups", mcp.Description("仪表盘用户组共享数组，包含 usrgrpid 和 permission")),
		),
		handler.CreateDashboardHandler,
	)

	s.AddTool(
		mcp.NewTool("create_graph_dashboard",
			mcp.WithDescription("自动创建图形仪表盘，支持两种模式：1) 多台主机相同图形名称聚合 2) 每台主机不同图形分列显示"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix 实例名称")),
			mcp.WithString("name", mcp.Required(), mcp.Description("仪表盘名称")),
			mcp.WithNumber("mode", mcp.Required(), mcp.Description("模式: 1-聚合模式(多台主机相同图形) 2-分列模式(每台主机不同图形)")),
			mcp.WithArray("hosts", mcp.Required(), mcp.Description("主机名称列表（字符串）")),
			mcp.WithArray("graph_names", mcp.Required(), mcp.Description("图形名称列表")),
			mcp.WithNumber("max_cols", mcp.Description("最大列数（可选，模式1默认2，模式2默认主机数量）")),
			mcp.WithNumber("max_rows", mcp.Description("最大行数（可选，默认10）")),
		),
		handler.CreateGraphDashboardHandler,
	)

	s.AddTool(
		mcp.NewTool("delete_dashboards",
			mcp.WithDescription("删除指定实例中的仪表盘"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix 实例名称")),
			mcp.WithArray("dashboard_ids", mcp.Required(), mcp.Description("仪表盘ID列表")),
		),
		handler.DeleteDashboardHandler,
	)
}
