/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-17 20:48:34
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-17 21:13:05
 * @FilePath: \zabbix-mcp-go\register\clientPool.go
 * @Description: 文件解释
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package register

import (
	"zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerClientPool 注册与 Zabbix 客户端池相关的 MCP 工具。
// 目前仅注册工具元信息；具体处理器在 handler 包未实现时暂留为 nil，以免影响构建。
func registerClientPool(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_instances_info",
			mcp.WithDescription("获取所有Zabbix实例的详细信息"),
		),
		handler.GetInstancesInfoHandler,
	)
}
