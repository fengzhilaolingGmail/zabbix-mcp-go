/*
 * @Author: fengzhilaoling fengzhilaoling@gmail.com
 * @Date: 2025-12-18 11:10:11
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2025-12-20 16:11:53
 * @FilePath: \zabbix-mcp-go\server\user.go
 * @Description: 用户相关功能
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

// GetUsers 调用底层 ClientProvider 执行 user.get，并返回解析后的列表。
// instanceName 为空时使用任意可用客户端，否则强制选择指定实例。
func GetUsers(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) ([]map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("user.get", spec)
	var users []map[string]interface{}
	callErr = client.Call(ctx, "user.get", adapted, &users)
	if callErr != nil {
		logger.L().Error("get user error: %s", callErr.Error())
		return nil, callErr
	}
	return users, nil
}

// 创建用户
func CreateUsers(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance, passwd string) (map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("user.create", spec)
	var users map[string]interface{}
	callErr = client.Call(ctx, "user.create", adapted, &users)
	if callErr != nil {
		logger.L().Error("create user error: %s", callErr.Error())
		return nil, callErr
	}
	users["passwd"] = passwd
	return users, nil
}

// 删除用户
// 用户组: 所有用户组 获取用户组user
// 获取用户组信息
func GetUserGroups(ctx context.Context, provider zabbix.ClientProvider, spec models.ParamSpec, instance string) ([]map[string]interface{}, error) {
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
	adapted := client.AdaptAPIParams("usergroup.get", spec)
	var userGroups []map[string]interface{}
	callErr = client.Call(ctx, "usergroup.get", adapted, &userGroups)
	if callErr != nil {
		logger.L().Error("get user group error: %s", callErr.Error())
		return nil, callErr
	}
	return userGroups, nil
}
