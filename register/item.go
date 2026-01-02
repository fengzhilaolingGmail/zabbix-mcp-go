/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-02 16:33:29
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-02 17:50:18
 * @FilePath: \zabbix-mcp-go\register\item.go
 * @Description: mcp注册监控项相关工具
 * Copyright (c) 2026 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package register

import (
	"zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerItem(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_items",
			mcp.WithDescription("获取实例Zabbix监控项信息,支持所有或模糊匹配. 注意: 必须提供主机过滤(host_ids/hostname)或监控项过滤(item_key/item_name)至少一种"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithArray("host_ids", mcp.Description("主机ID列表 (与hostname二选一,至少提供一个)")),
			mcp.WithString("hostname", mcp.Description("主机名称 (与host_ids二选一,至少提供一个)")),
			mcp.WithString("item_key", mcp.Description("监控项键(key) (与item_name二选一,至少提供一个)")),
			mcp.WithString("item_name", mcp.Description("监控项名称(name),需要使用通配符(与item_key二选一,至少提供一个)")),
		),
		handler.GetItemsHandler,
	)
}
