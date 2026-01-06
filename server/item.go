/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2026-01-02 16:16:35
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-01-02 17:50:58
 * @FilePath: \zabbix-mcp-go\server\item.go
 * @Description: 监控项相关服务
 * Copyright (c) 2026 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package server

import (
	"context"
	"fmt"
	"zabbixMcp/logger"
	"zabbixMcp/models"
	"zabbixMcp/zabbix"
)

func GetItems(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) ([]map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("item.get", spec)
	var items []map[string]interface{}
	callErr = client.Call(ctx, "item.get", adapted, &items)
	if callErr != nil {
		logger.L().Error("get items error: %s", callErr.Error())
		return nil, callErr
	}
	return items, nil
}
