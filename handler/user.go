/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 10:49:35
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-20 17:10:59
 * @FilePath: \zabbix-mcp-go\handler\user.go
 * @Description: 文件详情
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler

import (
	"context"
	"fmt"
	"strconv"

	"zabbixMcp/models"
	"zabbixMcp/server"
	"zabbixMcp/utils"

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
	spec := models.UserParams{Output: "extend"}
	if username != "" {
		// 兼容低版本
		spec.Alias = username
		spec.Filter = map[string]interface{}{"username": username}
		spec.GetAccess = true
		spec.SelectUsrgrps = []string{"usrgrpid", "name"}
	}
	users, err := server.GetUsers(ctx, clientPool, spec, instanceName)
	if err != nil {
		return nil, fmt.Errorf("调用 user.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(users)), nil
}

// CreateUsersHandler 通过注入的 ClientProvider 调用 user.create 并返回结果
func CreateUsersHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	username := ""
	name := ""
	userGroup := ""
	roleID := ""
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		if v, ok2 := args["username"].(string); ok2 {
			username = v
		}
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}
		if v, ok2 := args["userGroup"].(string); ok2 {
			userGroup = v
		}
		if v, ok2 := args["roleID"].(string); ok2 {
			roleID = v
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	passwd, err := utils.GenerateSecurePassword(12)
	if err != nil {
		return nil, fmt.Errorf("生成密码失败: %w", err)
	}
	// 使用 server 层处理业务逻辑
	spec := models.UserParams{
		UserName:  username,
		Name:      name,
		Passwd:    passwd,
		Roleid:    roleID,
		UserGroup: userGroup,
	}
	users, err := server.CreateUsers(ctx, clientPool, spec, instanceName, passwd)
	if err != nil {
		return nil, fmt.Errorf("调用 user.create 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(users)), nil
}

func UpdateUsersHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	name := ""
	// surname := ""
	userid := ""
	usrgrps := []string{}
	updatePasswd := false
	passwd := ""
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
		// if v, ok2 := args["surname"].(string); ok2 {
		// 	surname = v
		// }
		if v, ok2 := args["name"].(string); ok2 {
			name = v
		}
		if v, ok2 := args["usrgrps"].([]string); ok2 {
			usrgrps = v
		}
		if v, ok2 := args["userid"].(string); ok2 {
			userid = v
		}
		if v, ok2 := args["updatePasswd"].(bool); ok2 {
			updatePasswd = v
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	// 使用 server 层处理业务逻辑
	spec := models.UserParams{Userid: userid}
	// if surname != "" {
	// 	spec.Surname = surname
	// }
	if name != "" {
		spec.Name = name
	}
	if len(usrgrps) > 0 {
		spec.Usrgrps = usrgrps
	}
	if updatePasswd {
		passwd, err := utils.GenerateSecurePassword(12)
		if err != nil {
			return nil, fmt.Errorf("生成密码失败: %w", err)
		}
		spec.Passwd = passwd
		spec.CurrentPasswd = passwd
	}
	users, err := server.UpdateUser(ctx, clientPool, spec, instanceName, passwd)
	if err != nil {
		return nil, fmt.Errorf("调用 user.update 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(users)), nil
}

// GetUsersHandler 通过注入的 ClientProvider 调用 user.get 并返回结果
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
