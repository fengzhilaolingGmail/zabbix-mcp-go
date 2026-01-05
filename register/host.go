/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-02 15:33:32
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-05 17:58:21
 * @FilePath: \zabbix-mcp-go\register\host.go
 * @Description: 文件解释
 * Copyright (c) 2026 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
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
			mcp.WithString("active_available", mcp.Description("主机状态,0-未知,1-在线,2-离线,默认为1")),
			mcp.WithString("type", mcp.Description("主机类型,1 - Agent 2 - SNMP 3 - IPMI 4 - JMX,默认为1")),
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
	s.AddTool(
		mcp.NewTool("create_host",
			mcp.WithDescription("在指定实例中创建主机"),
			mcp.WithString("instance", mcp.Description("Zabbix 实例名称")),
			mcp.WithString("host", mcp.Required(), mcp.Required(), mcp.Description("主机技术名,host 字段")),
			mcp.WithString("name", mcp.Required(), mcp.Description("主机可见名称")),
			mcp.WithArray("groups", mcp.Required(), mcp.Description("主机组对象数组,指定 groupid 字段")),
			mcp.WithArray("interfaces", mcp.Description("主机接口数组")),
			mcp.WithArray("templateids", mcp.Description("模板ID数组")),
			mcp.WithArray("templates", mcp.Description("模板对象数组，包含 templateid 字段")),
			mcp.WithArray("tags", mcp.Description("主机标签数组")),
			mcp.WithArray("macros", mcp.Description("用户宏数组")),
			mcp.WithObject("inventory", mcp.Description("主机清单对象")),
		),
		handler.CreateHostHandler,
	)
}
