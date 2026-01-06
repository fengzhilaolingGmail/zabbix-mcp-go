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
	s.AddTool(
		mcp.NewTool("get_history_compare",
			mcp.WithDescription("获取实例Zabbix监控历史同比数据（current vs previous）"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("host_ids", mcp.Required(), mcp.Description("主机ID列表 ")),
			mcp.WithArray("item_ids", mcp.Required(), mcp.Description("监控项ID列表")),
			mcp.WithString("start_time", mcp.Description("开始时间，支持多种格式或与 time_range 二选一")),
			mcp.WithString("end_time", mcp.Description("结束时间，支持多种格式或与 time_range 二选一")),
			mcp.WithString("time_range", mcp.Description("相对时间范围，例如: 7d 15h 30m，可与 start_time/end_time 二选一")),
			mcp.WithString("period", mcp.Description("比较粒度: 'hour' 或 'day'，默认 'day'")),
			mcp.WithString("pct_format", mcp.Description("pct_change 格式: 'number' (默认浮点两位) 或 'string' (带% 的字符串)")),
			mcp.WithString("timezone", mcp.Description("时区，例如 Asia/Shanghai，默认 Asia/Shanghai")),
			mcp.WithNumber("history", mcp.Description("历史数据类型: 0 (默认)- 数值型float; 3 - 无符号数值型; 等")),
		),
		handler.GetHistoryCompareHandler,
	)
}
