/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:10:11
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-18 11:19:19
 * @FilePath: \zabbix-mcp-go\server\user.go
 * @Description: 用户相关功能
 * @Copyright: Copyright (c) 2025 by fengzhilaoling@gmail.com, All Rights Reserved.
 */

package server

import (
	"context"
	"fmt"

	"zabbixMcp/models"
	"zabbixMcp/zabbix"
)

// GetUsers 调用底层 ClientProvider 执行 user.get，并返回解析后的列表
func GetUsers(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec) ([]map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("user.get", spec)
	var users []map[string]interface{}
	callErr = client.Call(ctx, "user.get", adapted, &users)
	if callErr != nil {
		return nil, callErr
	}
	return users, nil
}
