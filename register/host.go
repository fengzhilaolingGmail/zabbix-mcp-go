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
		mcp.NewTool("find_host_by_name", mcp.WithDescription("通过主机名称获取主机信息,模糊匹配需要启用search参数"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostnames", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithBoolean("select_discoveries", mcp.Description("是否查询主机低级发现规则")), // ok
			mcp.WithBoolean("select_discovery_data", mcp.Description("是否查询主机发现数据")),
			mcp.WithBoolean("select_discovery_rule", mcp.Description("是否查询主机低级发现规则的rule")),
			mcp.WithBoolean("select_discovery_rule_prototype", mcp.Description("是否查询主机低级发现规则原型")),
			mcp.WithBoolean("select_graphs", mcp.Description("是否查询主机图形")), // ok
			mcp.WithBoolean("select_host_discovery", mcp.Description("是否查询主机发现数据")),
			mcp.WithBoolean("select_host_groups", mcp.Description("是否查询主机组")),
			mcp.WithBoolean("select_http_tests", mcp.Description("是否查询主机Web检查")),
			mcp.WithBoolean("select_interfaces", mcp.Description("是否查询主机接口")),
			mcp.WithBoolean("select_inventory", mcp.Description("是否查询主机清单数据")),
			mcp.WithBoolean("select_items", mcp.Description("是否查询主机监控项")),
			mcp.WithBoolean("select_macros", mcp.Description("是否查询主机宏")),
			mcp.WithBoolean("select_parent_templates", mcp.Description("是否查询主机模板")), // ok
			mcp.WithBoolean("select_dashboards", mcp.Description("是否查询主机仪表盘")),
			mcp.WithBoolean("select_tags", mcp.Description("是否查询主机标签")),
			mcp.WithBoolean("select_inherited_tags", mcp.Description("是否查询主机继承标签")),
			mcp.WithBoolean("select_triggers", mcp.Description("是否查询主机触发器")),
			mcp.WithBoolean("select_value_maps", mcp.Description("是否查询主机值映射")),
			mcp.WithBoolean("search", mcp.Description("是否启用模糊搜索")),
		),
		handler.GetHostsHandler,
	)
}
