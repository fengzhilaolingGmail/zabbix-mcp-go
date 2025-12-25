/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:20:36
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-25 11:30:37
 * @FilePath: \zabbix-mcp-go\handler\host.go
 * @Description: 文件详情
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler

import (
	"context"
	"fmt"

	"zabbixMcp/logger"
	"zabbixMcp/models"
	"zabbixMcp/server"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetHostsHandler 通过注入的 ClientProvider 调用 host.get 并返回结果
func GetHostsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance := ""
	hostnames := []string{}
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instance = v
		}
		if arr, ok := args["hostnames"].([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok && s != "" {
					hostnames = append(hostnames, s)
				}
			}
		}
	}
	if clientPool == nil {
		return mcp.NewToolResultStructuredOnly(makeResult([]map[string]interface{}{})), nil
	}
	spec := models.HostParams{Output: "extend", SelectInterfaces: "extend"}
	logger.L().Infof("instance: %s, hostname: %v", instance, hostnames)
	if len(hostnames) > 0 {
		spec.Filter = map[string]interface{}{"host": hostnames}
		spec.SelectHostGroups = []string{"groupid", "name"}
		spec.SelectParentTemplates = []string{"templateid", "name"}
	}
	hosts, err := server.GetHosts(ctx, clientPool, spec, instance)
	if err != nil {
		return nil, fmt.Errorf("调用 host.get 失败: %w", err)
	}
	return mcp.NewToolResultStructuredOnly(makeResult(hosts)), nil
}

// 通过主机组查询
// 通过主机名查询 详细信息 ()
