/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-06 18:56:08
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-09 14:32:09
 * @FilePath: \zabbix-mcp-go\server\dashboard.go
 * @Description: 仪表盘相关功能
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

// CreateDashboard 调用底层 ClientProvider 执行 dashboard.create，并返回创建结果。
// instance 为空时使用任意可用客户端，否则强制选择指定实例。
func CreateDashboard(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) (map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("dashboard.create", spec)
	var dashboard map[string]interface{}
	callErr = client.Call(ctx, "dashboard.create", adapted, &dashboard)
	if callErr != nil {
		logger.L().Error("create dashboard error: %s", callErr.Error())
		return nil, callErr
	}
	return dashboard, nil
}

// GetDashboard 调用底层 ClientProvider 执行 dashboard.get，并返回查询结果。
// instance 为空时使用任意可用客户端，否则强制选择指定实例。
func GetDashboard(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) ([]map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("dashboard.get", spec)
	var dashboard []map[string]interface{}
	callErr = client.Call(ctx, "dashboard.get", adapted, &dashboard)
	if callErr != nil {
		logger.L().Error("get dashboard error: %s", callErr.Error())
		return nil, callErr
	}
	return dashboard, nil
}
