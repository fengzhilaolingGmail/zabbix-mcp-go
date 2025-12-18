/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-17 20:56:38
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 19:11:23
 * @FilePath: \zabbix-mcp-go\handler\instances.go
 * @Description: 文件解释
 * Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package handler

import (
	"context"

	"zabbixMcp/logger"
	"zabbixMcp/server"
	"zabbixMcp/zabbix"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetInstancesInfoHandler 返回池中所有实例的信息
// 签名为 mcp 工具处理器，返回可序列化的结构或错误
func GetInstancesInfoHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceName := ""
	// req.Params.Arguments may be typed as interface{}; try to assert to a map first
	if args, ok := req.Params.Arguments.(map[string]interface{}); ok {
		if v, ok2 := args["instance"].(string); ok2 {
			instanceName = v
		}
	}
	if clientPool == nil {
		// 返回统一包装的空列表，便于客户端解析
		logger.Info("GetInstancesInfoHandler: clientPool is nil", nil, nil)
		return mcp.NewToolResultStructuredOnly(makeResult([]zabbix.ClientInfo{})), nil
	}
	infos, err := server.GetInstancesInfo(clientPool, instanceName)
	if err != nil {
		logger.L().Errorf("获取实例信息失败: %v", err)
		return nil, err
	}
	return mcp.NewToolResultStructuredOnly(makeResult(infos)), nil
}
