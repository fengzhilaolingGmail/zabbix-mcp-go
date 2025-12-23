/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-20 17:15:25
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-23 20:00:44
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

	// 使用 server 层处理业务逻辑
	spec := models.UserGroup{Output: "extend", Status: status}
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

func CreateUserGroupHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	name := ""
	hostgroupRights := map[string]int{}
	templategroupRights := map[string]int{}
	tagFilters := map[string]string{}
	userids := []string{}
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}
		if raw, ok2 := args["hostgroup_rights"]; ok2 {
			hostgroupRights = toStringIntMap(raw)
		} else if raw, ok2 := args["hostgroupRights"]; ok2 {
			hostgroupRights = toStringIntMap(raw)
		}
		if raw, ok2 := args["templategroup_rights"]; ok2 {
			templategroupRights = toStringIntMap(raw)
		} else if raw, ok2 := args["templategroupRights"]; ok2 {
			templategroupRights = toStringIntMap(raw)
		}
		if raw, ok2 := args["tag_filters"]; ok2 {
			tagFilters = toStringStringMap(raw)
		} else if raw, ok2 := args["tagFilters"]; ok2 {
			tagFilters = toStringStringMap(raw)
		}
		if arr, ok := args["userids"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					userids = append(userids, s)
				}
			}
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	// 使用 server 层处理业务逻辑
	spec := models.UserGroup{
		Name:                name,
		HostgroupRights:     hostgroupRights,
		TemplategroupRights: templategroupRights,
		TagFilters:          tagFilters,
		Users:               userids,
	}
	userGroups, err := server.CreateUserGroup(ctx, clientPool, spec, instanceName)
	if err != nil {
		return nil, fmt.Errorf("调用 usergroup.create 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(userGroups)), nil
}

func toStringIntMap(val interface{}) map[string]int {
	result := map[string]int{}
	switch typed := val.(type) {
	case map[string]int:
		for k, v := range typed {
			result[k] = v
		}
	case map[int]int:
		for k, v := range typed {
			result[strconv.Itoa(k)] = v
		}
	case map[string]interface{}:
		for k, v := range typed {
			if intVal, ok := toInt(v); ok {
				result[k] = intVal
			}
		}
	case map[interface{}]interface{}:
		for k, v := range typed {
			key := fmt.Sprintf("%v", k)
			if intVal, ok := toInt(v); ok {
				result[key] = intVal
			}
		}
	}
	return result
}

func toStringStringMap(val interface{}) map[string]string {
	result := map[string]string{}
	switch typed := val.(type) {
	case map[string]string:
		for k, v := range typed {
			result[k] = v
		}
	case map[string]interface{}:
		for k, v := range typed {
			if strVal, ok := v.(string); ok {
				result[k] = strVal
			}
		}
	case map[interface{}]interface{}:
		for k, v := range typed {
			key := fmt.Sprintf("%v", k)
			if strVal, ok := v.(string); ok {
				result[key] = strVal
			}
		}
	}
	return result
}

func toInt(val interface{}) (int, bool) {
	switch v := val.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case float32:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed, true
		}
	}
	return 0, false
}
