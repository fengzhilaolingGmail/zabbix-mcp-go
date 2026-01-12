/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-02 15:33:32
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-12 13:32:13
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

func registerGetHost(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_host_by_name",
			mcp.WithDescription("获取实例Zabbix主机信息,支持所有或模糊匹配"),
			mcp.WithString("hostnames", mcp.Required(), mcp.Description("主机名称列表")),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithString("active_available", mcp.Description("主机状态,0-未知,1-在线,2-离线,默认为1")),
			mcp.WithString("type", mcp.Description("主机类型,1 - Agent 2 - SNMP 3 - IPMI 4 - JMX,默认为1")),
		),
		handler.GetHostForNameHandler,
	)
	// 获取主机列表的图形
	s.AddTool(
		mcp.NewTool("get_host_graphs", mcp.WithDescription("获取主机列表的图形"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:graphs 表示获取图形")),
		),
		handler.GetHostsHandler,
	)
	// 获取主机列表的图形
	s.AddTool(
		mcp.NewTool("get_host_templates", mcp.WithDescription("获取主机列表的模板"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:templates 表示获取模板")),
		),
		handler.GetHostsHandler,
	)
	// 获取主机列表的监控项
	s.AddTool(
		mcp.NewTool("get_host_items", mcp.WithDescription("获取主机列表的监控项"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:items 表示获取监控项")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_triggers", mcp.WithDescription("获取主机列表的触发器"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:triggers 表示获取触发器")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_macros", mcp.WithDescription("获取主机列表的宏"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:macros 表示获取宏")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_value_maps", mcp.WithDescription("获取主机列表的数值映射"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:value_maps 表示获取数值映射")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_groups", mcp.WithDescription("获取主机列表的主机组"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:groups 表示获取主机组")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_tags", mcp.WithDescription("获取主机列表的标签"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:tags 表示获取标签")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_dashboards", mcp.WithDescription("获取主机列表的仪表盘"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:dashboards 表示获取仪表盘")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_discoveries", mcp.WithDescription("获取主机列表的低级发现规则"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:discoveries 表示获取低级发现规则")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_discovery_rule", mcp.WithDescription("获取主机列表的发现事件列表"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:discovery_rule 表示获取发现事件列表")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_discovery_data", mcp.WithDescription("获取主机列表的最近一次发现数据"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:discovery_data 表示获取最近一次发现数据")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_discovery_rule_prototype", mcp.WithDescription("获取主机列表的发现规则原型配置"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:discovery_rule_prototype 表示获取发现规则原型配置")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_discovery", mcp.WithDescription("获取主机列表的发现记录"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:discovery 表示获取最近一次发现数据")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_http_tests", mcp.WithDescription("获取主机列表的Web检查"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:http_tests 表示获取主机Web检查")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_interfaces", mcp.WithDescription("获取主机列表的接口"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:interfaces 表示获取主机接口")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_inventory", mcp.WithDescription("获取主机列表的清单数据"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithBoolean("is_detailed", mcp.Required(), mcp.Description("是否查询详细信息: 默认 false")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:inventory 表示获取主机清单数据")),
		),
		handler.GetHostsHandler,
	)
	s.AddTool(
		mcp.NewTool("get_host_inherited_tags", mcp.WithDescription("获取主机列表的继承标签"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("hostids", mcp.Required(), mcp.Description("主机ID列表,数量不要超过4个")),
			mcp.WithString("style", mcp.Required(), mcp.Description("默认:inherited_tags 表示获取主机继承标签")),
		),
		handler.GetHostsHandler,
	)
}

func registerHost(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("create_host",
			mcp.WithDescription("在指定实例中创建主机"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix 实例名称")),
			mcp.WithString("host", mcp.Required(), mcp.Description("主机技术名,host 字段")),
			mcp.WithString("name", mcp.Required(), mcp.Description("主机可见名称")),
			mcp.WithArray("groups", mcp.Required(), mcp.Description("主机组对象数组,指定 groupid 字段")),
			mcp.WithArray("interfaces", mcp.Description(`主机接口数组，元素为对象，必须字段：
		{
		  "type":  1,        // int，1=Agent 2=SNMP 3=IPMI 4=JMX
		  "main":  1,        // int，1 主接口 0 非主（只能有一个主）
		  "useip": 1,        // int，1 用 IP 0 用 DNS
		  "ip":    "1.1.1.1",// string，当 useip=1 时必填
		  "dns":   "",       // string，当 useip=0 时必填
		  "port":  "10050"   // string，端口
		}`)),
			mcp.WithArray("templateids", mcp.Description("模板ID数组")),
			mcp.WithArray("templates", mcp.Description("模板对象数组，包含 templateid 字段")),
			mcp.WithArray("tags", mcp.Description("主机标签数组")),
			mcp.WithArray("macros", mcp.Description("用户宏数组")),
			mcp.WithObject("inventory", mcp.Description("主机清单对象")),
		),
		handler.CreateHostHandler,
	)
	// host.update
	s.AddTool(
		mcp.NewTool("update_host",
			mcp.WithDescription("更新主机属性"),
			mcp.WithString("instance", mcp.Description("Zabbix 实例名称")),
			mcp.WithString("hostid", mcp.Description("主机ID")),
			mcp.WithArray("hostids", mcp.Description("主机ID数组")),
			mcp.WithString("host", mcp.Description("主机技术名")),
			mcp.WithString("name", mcp.Description("主机可见名")),
			// 以下字段在 update 场景中表示替换（replace）语义：传入后会替换对应关联，未列出的将被移除
			mcp.WithArray("groups", mcp.Description("主机组对象数组，仅包含 groupid 字段，用于替换当前主机组")),
			mcp.WithArray("interfaces", mcp.Description("主机接口对象数组，用于替换当前接口")),
			mcp.WithArray("templates", mcp.Description("模板对象数组，仅包含 templateid 字段，用于替换当前关联模板")),
			mcp.WithArray("templates_clear", mcp.Description("模板对象数组，仅包含 templateid 字段，用于解除并清除模板关联")),
			mcp.WithArray("tags", mcp.Description("主机标签数组，用于替换当前标签")),
			mcp.WithArray("macros", mcp.Description("用户宏数组，用于替换当前宏")),
			mcp.WithObject("inventory", mcp.Description("主机资产清单对象，用于替换清单")),
		),
		handler.UpdateHostHandler,
	)
}
