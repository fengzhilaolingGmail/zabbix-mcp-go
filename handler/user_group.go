/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-20 17:15:25
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-20 17:16:08
 * @FilePath: \zabbix-mcp-go\handler\user_group.go
 * @Description: 用户组
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package handler

import (
	"context"
	"fmt"
	"strconv"
	"zabbixMcp/models"
	"zabbixMcp/server"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetUserGroupsHandler 通过注入的 ClientProvider 调用 usergroup.get 并返回结果
func GetUserGroupsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	name := ""
	status := "0"
	selectUsers := false
	selectRights := false
	selectTagFilters := false
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}
		if v, ok2 := args["status"].(string); ok2 {
			status = v
		}
		if v, ok2 := args["selectUsers"].(bool); ok2 {
			selectUsers = v
		}
		if v, ok2 := args["selectRights"].(bool); ok2 {
			selectRights = v
		}
		if v, ok2 := args["selectTagFilters"].(bool); ok2 {
			selectTagFilters = v
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	statusInt, err := strconv.Atoi(status)
	if err != nil {
		return nil, err
	}
	// 使用 server 层处理业务逻辑
	spec := models.UserGroup{Output: "extend", Status: statusInt}
	if name != "" {
		// 兼容低版本
		spec.Filter = map[string]interface{}{"name": name}
	}
	if selectUsers {
		spec.SelectUsers = selectUsers
	}
	if selectRights {
		spec.SelectRights = selectRights
	}
	if selectTagFilters {
		spec.SelectTagFilters = selectTagFilters
	}
	userGroups, err := server.GetUserGroups(ctx, clientPool, spec, instanceName)
	if err != nil {
		return nil, fmt.Errorf("调用 usergroup.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(userGroups)), nil
}
