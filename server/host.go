/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:13:06
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 11:13:07
 * @FilePath: \zabbix-mcp-go\server\host.go
 * @Description: 主机相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package server

import (
	"context"
	"fmt"

	"zabbixMcp/models"
	"zabbixMcp/zabbix"
)

// GetHosts 调用底层 ClientProvider 执行 host.get，并返回解析后的列表
func GetHosts(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec) ([]map[string]interface{}, error) {
	if provider == nil {
		return nil, fmt.Errorf("no zabbix client")
	}
	lease, err := provider.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	var callErr error
	defer func() { lease.Release(callErr) }()
	client := lease.Client()
	adapted := client.AdaptAPIParams("host.get", spec)
	var hosts []map[string]interface{}
	callErr = client.Call(ctx, "host.get", adapted, &hosts)
	if callErr != nil {
		return nil, callErr
	}
	return hosts, nil
}
