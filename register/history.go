package register

import (
	"zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerHistory(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_history_by_time",
			mcp.WithDescription("获取实例Zabbix监控历史数据"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("host_ids", mcp.Required(), mcp.Description("主机ID列表 ")),
			mcp.WithArray("item_ids", mcp.Required(), mcp.Description("监控项ID列表")),
			mcp.WithString("start_time", mcp.Required(), mcp.Description("开始时间")),
			mcp.WithString("end_time", mcp.Required(), mcp.Description("结束时间")),
			mcp.WithBoolean("summary", mcp.Description("是否汇总数据")),
			mcp.WithNumber("history", mcp.Description("历史数据类型: 0 (默认)- 数值型float; 1 - 字符型; 2 - 日志型; 3 - 无符号数值型; 4 - 文本型; 5 - 二进制型")),
		),
		handler.GetHistoryHandler,
	)
	s.AddTool(
		mcp.NewTool("get_history_by_range",
			mcp.WithDescription("获取实例Zabbix监控历史数据"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("host_ids", mcp.Required(), mcp.Description("主机ID列表 ")),
			mcp.WithArray("item_ids", mcp.Required(), mcp.Description("监控项ID列表")),
			mcp.WithString("time_range", mcp.Required(), mcp.Description("时间范围: 7d 15h 30m")),
			mcp.WithBoolean("summary", mcp.Description("是否汇总数据")),
			mcp.WithNumber("history", mcp.Description("历史数据类型: 0 (默认)- 数值型float; 1 - 字符型; 2 - 日志型; 3 - 无符号数值型; 4 - 文本型; 5 - 二进制型")),
		),
		handler.GetHistoryHandler,
	)
}
