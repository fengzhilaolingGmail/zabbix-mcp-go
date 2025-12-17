/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-17 20:56:38
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-17 21:13:41
 * @FilePath: \zabbix-mcp-go\handler\pool_handler.go
 * @Description: 文件解释
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler

import (
	"context"

	zabbix "zabbixMcp/zabbix"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetInstancesInfoHandler 返回池中所有实例的信息
// 签名为 mcp 工具处理器，返回可序列化的结构或错误
func GetInstancesInfoHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if clientPool == nil {
		// 返回空列表而非错误，便于前端友好处理
		return &mcp.CallToolResult{StructuredContent: []zabbix.ClientInfo{}}, nil
	}
	infos := clientPool.Info()
	return &mcp.CallToolResult{StructuredContent: infos}, nil
}
