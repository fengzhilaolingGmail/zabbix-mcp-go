/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:20:36
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 19:50:37
 * @FilePath: \zabbix-mcp-go\handler\host.go
 * @Description: 文件详情
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler

import (
	"context"
	"fmt"
	"zabbixMcp/server"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetHostsHandler 通过注入的 ClientProvider 调用 host.get 并返回结果
func GetHostsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}

	params := map[string]interface{}{"output": "extend"}
	hosts, err := server.GetHosts(ctx, clientPool, params)
	if err != nil {
		return nil, fmt.Errorf("调用 host.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}

// 通过主机组查询
// 通过主机名查询 详细信息 ()
