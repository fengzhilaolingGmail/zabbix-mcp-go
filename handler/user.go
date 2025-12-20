/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 10:49:35
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-19 18:08:14
 * @FilePath: \zabbix-mcp-go\handler\user.go
 * @Description: 文件详情
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler

import (
	"context"
	"fmt"

	"zabbixMcp/models"
	"zabbixMcp/server"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetUsersHandler 通过注入的 ClientProvider 调用 user.get 并返回结果
func GetUsersHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	username := ""
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if v, ok2 := args["username"].(string); ok2 {
			username = v
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	// 使用 server 层处理业务逻辑
	spec := models.UserGetParams{Output: "extend"}
	if username != "" {
		// 兼容低版本
		spec.Alias = username
		spec.Filter = map[string]interface{}{"username": username}
	}
	users, err := server.GetUsers(ctx, clientPool, spec, instanceName)
	if err != nil {
		return nil, fmt.Errorf("调用 user.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(users)), nil
}
