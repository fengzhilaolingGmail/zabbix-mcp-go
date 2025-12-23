/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-20 17:13:46
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-23 19:45:28
 * @FilePath: \zabbix-mcp-go\register\user_group.go
 * @Description: 用户组功能注册
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package register

import (
	"zabbixMcp/handler"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerUserGroup(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("get_groups", mcp.WithDescription("获取所有Zabbix用户组信息"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithString("name", mcp.Description("用户组名称")),
			mcp.WithString("status", mcp.Description("用户组状态: 0启用 1禁用 默认: 0")),
			mcp.WithBoolean("selectUsers", mcp.Description("是否获取用户组下用户列表 默认: false")),
			mcp.WithBoolean("selectRights", mcp.Description("是否获取用户组权限列表 默认: false")),
			mcp.WithBoolean("selectTagFilters", mcp.Description("是否获取用户组标签过滤器列表 默认: false")),
		),
		handler.GetUserGroupsHandler,
	)
	s.AddTool(
		mcp.NewTool("create_group", mcp.WithDescription("创建Zabbix用户组"),
			mcp.WithString("instance", mcp.Required(), mcp.Description("Zabbix实例名称必须填")),
			mcp.WithString("name", mcp.Required(), mcp.Description("用户组名称")),
			mcp.WithObject(
				"hostgroup_rights",
				mcp.Description("用户组权限列表, Key为主机组ID字符串, Value为权限整型"),
				mcp.AdditionalProperties(map[string]any{
					"type":        "integer",
					"description": "权限值, 0=拒绝访问 ,2=只读, 3=读写等",
				}),
			),
			mcp.WithObject(
				"templategroup_rights",
				mcp.Description("用户组模板权限列表, Key为模板组ID字符串, Value为权限整型"),
				mcp.AdditionalProperties(map[string]any{
					"type":        "integer",
					"description": "权限值, 0=拒绝访问 ,2=只读, 3=读写等",
				}),
			),
			mcp.WithObject(
				"tag_filters",
				mcp.Description("用户组标签过滤器列表, Key为标签字符串, Value为标签值字符串"),
				mcp.AdditionalProperties(map[string]any{
					"type":        "string",
					"description": "标签值字符串",
				}),
			),
			mcp.WithArray("userids", mcp.Description("Zabbix用户ID列表")),
		),
		handler.CreateUserGroupHandler,
	)
}
