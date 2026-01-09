/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 19:30:16
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-06 19:30:33
 * @FilePath: \zabbix-mcp-go\server\graph.go
 * @Description: 文件详情
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */
package server

import (
	"context"
	"fmt"

	"zabbixMcp/logger"
	"zabbixMcp/models"
	"zabbixMcp/zabbix"
)

// CreateGraph 调用底层 ClientProvider 执行 graph.create，并返回创建结果。
// instance 为空时使用任意可用客户端，否则强制选择指定实例。
func CreateGraph(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) (map[string]interface{}, error) {
	if provider == nil {
		return nil, fmt.Errorf("no zabbix client")
	}
	var (
		lease zabbix.ClientLease
		err   error
	)
	if instance != "" {
		lease, err = provider.AcquireByInstance(ctx, instance)
	} else {
		lease, err = provider.Acquire(ctx)
	}
	if err != nil {
		return nil, err
	}
	var callErr error
	defer func() { lease.Release(callErr) }()
	client := lease.Client()
	adapted := client.AdaptAPIParams("graph.create", spec)
	var graphs map[string]interface{}
	callErr = client.Call(ctx, "graph.create", adapted, &graphs)
	if callErr != nil {
		logger.L().Error("create graph error: %s", callErr.Error())
		return nil, callErr
	}
	return graphs, nil
}
