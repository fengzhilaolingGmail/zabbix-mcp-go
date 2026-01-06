/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-04 10:54:09
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-05 11:24:27
 * @FilePath: \zabbix-mcp-go\server\history.go
 * @Description: 历史数据获取功能实现
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

// GetHistory 调用 zabbix API history.get 获取历史数据
func GetHistory(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) ([]map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("history.get", spec)
	var histories []map[string]interface{}
	callErr = client.Call(ctx, "history.get", adapted, &histories)
	if callErr != nil {
		logger.L().Error("get history error: %s", callErr.Error())
		return nil, callErr
	}
	return histories, nil
}
