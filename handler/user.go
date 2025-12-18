/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 10:49:35
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 10:49:50
 * @FilePath: \zabbix-mcp-go\handler\user_host_handlers.go
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

// GetUsersHandler 通过注入的 ClientProvider 调用 user.get 并返回结果
func GetUsersHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}

	// 使用 server 层处理业务逻辑
	params := map[string]interface{}{"output": "extend"}
	users, err := server.GetUsers(ctx, clientPool, params)
	if err != nil {
		return nil, fmt.Errorf("调用 user.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(users)), nil
}
