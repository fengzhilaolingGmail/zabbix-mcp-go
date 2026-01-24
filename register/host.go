/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-02 15:33:32
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-24 17:13:09
 * @FilePath: \zabbix-mcp-go\register\host.go
 * @Description: 文件解释
 * Copyright (c) 2026 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package register

import (
	"zabbixMcp/handler"
	"zabbixMcp/models"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerGetHost(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_host_by_name",
			mcp.WithDescription("获取实例Zabbix主机信息,支持通配符"),
			mcp.WithString("hostname", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithBoolean("searchWildcardsEnabled", mcp.Required(), mcp.Description("是否允许通配符搜索,默认: true"), mcp.DefaultBool(true)),
			mcp.WithNumber("limit", mcp.Required(), mcp.Description("返回主机数量,默认: 15"), mcp.DefaultNumber(15)),
		),
		handler.GetHostForNameHandler,
	)
	s.AddTool(
		mcp.NewTool("get_groups_by_hostname",
			mcp.WithDescription("根据主机名称获取主机的组信息"),
			mcp.WithString("hostname", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("selectHostGroups", mcp.Required(),
				mcp.Items([]string{"groupid", "name", "flags", "uuid"}),
				mcp.Description("要查询的组属性列表"), mcp.DefaultArray([]string{"extend"})),
			mcp.WithBoolean("searchWildcardsEnabled", mcp.Required(), mcp.Description("是否允许通配符搜索,默认: true"), mcp.DefaultBool(true)),
			mcp.WithNumber("limitSelects", mcp.Required(), mcp.Description("返回组数量,默认: 100"), mcp.DefaultNumber(100)),
		),
		handler.GetHostForNameHandler,
	)
	s.AddTool(
		mcp.NewTool("get_graphs_by_hostname",
			mcp.WithDescription("根据主机名称获取主机的图形信息,count使用时仅返回图形数量"),
			mcp.WithString("hostname", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("selectGraphs", mcp.Required(),
				mcp.Items([]string{"graphid", "name", "height", "width", "graphtype", "flags", "show_legend", "show_3d", "yaxismax", "yaxismin", "ymax_type", "ymin_type", "templateid", "uuid"}),
				mcp.Description("要查询的图形属性列表"), mcp.DefaultArray([]string{"extend"})),
			mcp.WithBoolean("count", mcp.Description("是否返回图形数量,默认: false"), mcp.DefaultBool(false)),
			mcp.WithBoolean("searchWildcardsEnabled", mcp.Required(), mcp.Description("是否允许通配符搜索,默认: true"), mcp.DefaultBool(true)),
			mcp.WithNumber("limitSelects", mcp.Required(), mcp.Description("返回图形数量,默认: 100"), mcp.DefaultNumber(100)),
		),
		handler.GetHostForNameHandler,
	)
	s.AddTool(
		mcp.NewTool("get_templates_by_hostname",
			mcp.WithDescription("根据主机名称获取主机的模板信息,count使用时仅返回模板数量"),
			mcp.WithString("hostname", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("selectParentTemplates", mcp.Required(),
				mcp.Items([]string{"templateid", "name", "host", "description", "vendor_version", "vendor_name", "uuid"}),
				mcp.Description("要查询的模板属性列表"), mcp.DefaultArray([]string{"extend"})),
			mcp.WithBoolean("count", mcp.Description("是否返回模板数量,默认: false"), mcp.DefaultBool(false)),
			mcp.WithBoolean("searchWildcardsEnabled", mcp.Required(), mcp.Description("是否允许通配符搜索,默认: true"), mcp.DefaultBool(true)),
			mcp.WithNumber("limitSelects", mcp.Required(), mcp.Description("返回模板数量,默认: 100"), mcp.DefaultNumber(100)),
		),
		handler.GetHostForNameHandler,
	)
	s.AddTool(
		mcp.NewTool("get_items_by_hostname",
			mcp.WithDescription("根据主机名称获取主机的监控项信息,count使用时仅返回监控项数量"),
			mcp.WithString("hostname", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("selectItems", mcp.Required(),
				mcp.Items([]string{"itemid", "delay", "hostid", "interfaceid", "key_", "name", "name_resolved", "type", "url",
					"value_type", "allow_traps", "authtype", "description", "error", "flags", "follow_redirects", "headers",
					"history", "http_proxy", "inventory_link", "ipmi_sensor", "jmx_endpoint", "lastclock", "lastns", "lastvalue",
					"logtimefmt", "master_itemid", "output_format", "params", "parameters", "password", "post_type", "posts",
					"prevvalue", "privatekey", "publickey", "query_fields", "request_method", "retrieve_mode", "snmp_oid",
					"ssl_cert_file", "ssl_key_file", "ssl_key_password", "state", "status", "status_codes", "templateid",
					"timeout", "trapper_hosts", "trends", "units", "username", "uuid", "valuemapid", "verify_host", "verify_peer"}),
				mcp.Description("要查询的监控项属性列表"), mcp.DefaultArray([]string{"extend"})),
			mcp.WithBoolean("count", mcp.Description("是否返回监控项数量,默认: false"), mcp.DefaultBool(false)),
			mcp.WithBoolean("searchWildcardsEnabled", mcp.Required(), mcp.Description("是否允许通配符搜索,默认: true"), mcp.DefaultBool(true)),
			mcp.WithNumber("limitSelects", mcp.Required(), mcp.Description("返回监控项数量,默认: 100"), mcp.DefaultNumber(100)),
		),
		handler.GetHostForNameHandler,
	)
	s.AddTool(
		mcp.NewTool("get_macros_by_hostname",
			mcp.WithDescription("根据主机名称获取主机的宏信息,count使用时仅返回宏数量"),
			mcp.WithString("hostname", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("selectMacros", mcp.Required(),
				mcp.Items([]string{"hostmacroid", "hostid", "macro", "value", "type", "description", "automatic"}),
				mcp.Description("要查询的宏属性列表"), mcp.DefaultArray([]string{"extend"})),
			mcp.WithBoolean("searchWildcardsEnabled", mcp.Required(), mcp.Description("是否允许通配符搜索,默认: true"), mcp.DefaultBool(true)),
			mcp.WithNumber("limitSelects", mcp.Required(), mcp.Description("返回宏数量,默认: 100"), mcp.DefaultNumber(100)),
		),
		handler.GetHostForNameHandler,
	)
	s.AddTool(
		mcp.NewTool("get_tags_by_hostname",
			mcp.WithDescription("根据主机名称获取主机的标签信息,count使用时仅返回标签数量"),
			mcp.WithString("hostname", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("selectTags", mcp.Required(),
				mcp.Items([]string{"tag", "value", "automatic"}),
				mcp.Description("要查询的标签属性列表"), mcp.DefaultArray([]string{"extend"})),
			mcp.WithBoolean("searchWildcardsEnabled", mcp.Required(), mcp.Description("是否允许通配符搜索,默认: true"), mcp.DefaultBool(true)),
			mcp.WithNumber("limitSelects", mcp.Required(), mcp.Description("返回标签数量,默认: 100"), mcp.DefaultNumber(100)),
		),
		handler.GetHostForNameHandler,
	)
	s.AddTool(
		mcp.NewTool("get_triggers_by_hostname",
			mcp.WithDescription("根据主机名称获取主机的触发器信息,count使用时仅返回触发器数量"),
			mcp.WithString("hostname", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("selectTriggers", mcp.Required(),
				mcp.Items([]string{"triggerid", "description", "expression", "event_name", "opdata",
					"comments", "error", "flags", "lastchange", "priority", "state", "status", "templateid",
					"type", "url", "value", "recovery_mode", "recovery_expression", "correlation_mode",
					"correlation_tag", "manual_close", "uuid"}),
				mcp.Description("要查询的触发器属性列表"), mcp.DefaultArray([]string{"extend"})),
			mcp.WithBoolean("count", mcp.Description("是否返回触发器数量,默认: false"), mcp.DefaultBool(false)),
			mcp.WithBoolean("searchWildcardsEnabled", mcp.Required(), mcp.Description("是否允许通配符搜索,默认: true"), mcp.DefaultBool(true)),
			mcp.WithNumber("limitSelects", mcp.Required(), mcp.Description("返回触发器数量,默认: 100"), mcp.DefaultNumber(100)),
		),
		handler.GetHostForNameHandler,
	)
	s.AddTool(
		mcp.NewTool("get_dashboards_by_hostname",
			mcp.WithDescription("根据主机名称获取主机的看板(仪表盘)信息,count使用时仅返回看板数量"),
			mcp.WithString("hostname", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("selectDashboards", mcp.Required(),
				mcp.Items([]string{"dashboardid", "name", "templateid", "display_period", "auto_start", "uuid"}),
				mcp.Description("要查询的看板属性列表"), mcp.DefaultArray([]string{"extend"})),
			mcp.WithBoolean("count", mcp.Description("是否返回看板(仪表盘)数量,默认: false"), mcp.DefaultBool(false)),
			mcp.WithBoolean("searchWildcardsEnabled", mcp.Required(), mcp.Description("是否允许通配符搜索,默认: true"), mcp.DefaultBool(true)),
			mcp.WithNumber("limitSelects", mcp.Required(), mcp.Description("返回看板(仪表盘)数量,默认: 100"), mcp.DefaultNumber(100)),
		),
		handler.GetHostForNameHandler,
	)
}

func registerUpdateHost(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("update_host_groups",
			mcp.WithDescription("更新主机的主机组"),
			mcp.WithString("instance", mcp.Description("Zabbix 实例名称")),
			mcp.WithString("hostid", mcp.Required(), mcp.Description("主机ID")),
			mcp.WithArray("groups", mcp.Required(), mcp.Description("主机组ID列表")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:groups 表示更新主机组")),
		),
		handler.UpdateNewHostHandler,
	)
	s.AddTool(
		mcp.NewTool("update_host_templates",
			mcp.WithDescription("更新主机的主机模板"),
			mcp.WithString("instance", mcp.Description("Zabbix 实例名称")),
			mcp.WithString("hostid", mcp.Required(), mcp.Description("主机ID")),
			mcp.WithArray("templates", mcp.Required(), mcp.Description("模板ID列表")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:templates 表示更新主机模板")),
		),
		handler.UpdateNewHostHandler,
	)
	s.AddTool(
		mcp.NewTool("update_host_tags",
			mcp.WithDescription("更新主机的标签"),
			mcp.WithString("instance", mcp.Description("Zabbix 实例名称")),
			mcp.WithString("hostid", mcp.Required(), mcp.Description("主机ID")),
			mcp.WithArray("tags", mcp.Required(), mcp.Description("标签列表, 示例: `[{ 'tag': 'tagname', 'value': 'tagvalue' }]`")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:tags 表示更新主机标签")),
		),
		handler.UpdateNewHostHandler,
	)
	s.AddTool(
		mcp.NewTool("update_host_interfaces",
			mcp.WithDescription("更新主机的接口"),
			mcp.WithString("instance", mcp.Description("Zabbix 实例名称")),
			mcp.WithString("hostid", mcp.Required(), mcp.Description("主机ID")),
			mcp.WithArray("interfaces", mcp.Items([]models.ZabbixInterface{}), mcp.Required(), mcp.Description("接口列表")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:interfaces 表示更新主机接口")),
		),
		handler.UpdateNewHostHandler,
	)
	s.AddTool(
		mcp.NewTool("clear_host_templates",
			mcp.WithDescription("清除主机的主机模板(取消模板并清除关联)"),
			mcp.WithString("instance", mcp.Description("Zabbix 实例名称")),
			mcp.WithString("hostid", mcp.Required(), mcp.Description("主机ID")),
			mcp.WithArray("templates", mcp.Required(), mcp.Description("模板ID列表")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:templates_clear 表示清除主机模板")),
		),
		handler.UpdateNewHostHandler,
	)
}

func registerHost(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("create_host",
			mcp.WithDescription("在指定实例中创建主机"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix 实例名称")),
			mcp.WithString("host", mcp.Required(), mcp.Description("主机技术名,host 字段")),
			mcp.WithString("name", mcp.Required(), mcp.Description("主机可见名称")),
			mcp.WithArray("groups", mcp.Items([]models.Groups{}), mcp.Required(), mcp.Description("主机组id列表, 示例: `[{ 'groupid': '12345' },{ 'groupid': '67890' }]`")),
			mcp.WithArray("interfaces", mcp.Items([]models.ZabbixInterface{}), mcp.Required(), mcp.Description("接口列表")),
			mcp.WithArray("templates", mcp.Items([]models.Templates{}), mcp.Required(), mcp.Description("模板ID列表, 示例: `[{ 'templateid': '12345' },{ 'templateid': '67890' }]`")),
			mcp.WithArray("tags", mcp.Required(), mcp.Items([]models.Tag{}), mcp.Description("主机标签数组, 示例: `[{ 'tag': 'tagname', 'value': 'tagvalue' },{ 'tag': 'tagname', 'value': 'tagvalue' }]`")),
			mcp.WithArray("macros", mcp.Required(), mcp.Description("用户宏数组, 示例: `[{ 'macro': 'macroname', 'value': 'macrovalue', 'description': 'macrodescription' },{ 'macro': 'macroname', 'value': 'macrovalue', 'description': 'macrodescription' }]`")),
			// mcp.WithObject("inventory", mcp.Required(), mcp.Description("主机清单对象, 示例: `{ 'macaddress_a': '01234', 'macaddress_b': '56768' }`")),
		),
		handler.CreateHostHandler,
	)
	s.AddTool(
		mcp.NewTool("delete_host",
			mcp.WithDescription("删除主机"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix 实例名称")),
			mcp.WithArray("hostids", mcp.Items([]string{}), mcp.Required(), mcp.Description("主机ID数组")),
		),
		handler.DeleteHostsHandler,
	)
}
